// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 16, Puzzle 1.
// Proboscidea Volcanium - pressure release.

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func Int(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("%s is not an int: %v", s, err)
	}
	return v
}

// A priority Queue entry.
type PQE struct {
	node string
	dist int
}

// Priority Queue
type PQ struct {
	Q []PQE
}

func (q *PQ) Add(e PQE) {
	if len(q.Q) == 0 {
		q.Q = append(q.Q, e)
		return
	}
	if e.dist < q.Q[0].dist {
		q.Q = append([]PQE{e}, q.Q...)
		return
	}
	for n, t := range q.Q {
		if e.dist < t.dist {
			q.Q = append(q.Q[:n+1], q.Q[n:]...)
			q.Q[n] = e
			return
		}
	}
	q.Q = append(q.Q, e)
}

func (q *PQ) Pop() PQE {
	e := q.Q[0]
	q.Q = q.Q[1:]
	return e
}

type Valve struct {
	name     string
	flowRate int

	paths []*Valve

	// distance to other valves from this valve
	dists map[string]int
}

func (v Valve) String() string {
	return v.name
}

func (v Valve) Describe() string {
	names := []string{}
	for _, o := range v.paths {
		names = append(names, o.name)
	}
	return fmt.Sprintf("Valve %s; flow rate=%d; leads to %s", v.name, v.flowRate, strings.Join(names, ", "))
}

// Calculates (Djikstra) distance to dest and stores into dists
func (v *Valve) CalcDist(dest *Valve) {
	if rev, ok := dest.dists[v.name]; ok {
		v.dists[dest.name] = rev
		return
	}

	dists := DjikstraStart()
	been := map[string]bool{}
	dists[v.name] = 0
	q := PQ{}
	q.Add(PQE{node: v.name, dist: 0})
	for len(q.Q) > 0 {
		c := q.Pop()
		valve := Valves[c.node]
		been[valve.name] = true
		for _, n := range valve.paths {
			if been[n.name] {
				continue
			}
			if c.dist+1 < dists[n.name] {
				dists[n.name] = c.dist + 1
				q.Add(PQE{node: n.name, dist: c.dist + 1})
			}
		}
	}
	v.dists[dest.name] = dists[dest.name]
}

// Returns the potential pressure released if valve opens at tick,
// given limit ticks
func (v Valve) Potential(tick int, limit int) int {
	return (limit - tick - 1) * v.flowRate
}

func GetOrCreateValve(name string) *Valve {
	v := Valves[name]
	if v != nil {
		return v
	}
	v = &Valve{
		name:     name,
		flowRate: -1,
		dists:    map[string]int{},
	}
	Valves[v.name] = v
	return v
}

func GetOrCreateValves(nameList string) []*Valve {
	rv := []*Valve{}
	for _, name := range strings.Split(nameList, ", ") {
		rv = append(rv, GetOrCreateValve(name))
	}
	return rv
}

func ShowDot(save bool) {
	file, err := ioutil.TempFile("", "dot")
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString("digraph {\n")
	for _, v := range Valves {
		file.WriteString(fmt.Sprintf("%s [label=\"%s flow=%d\"]\n", v.name, v.name, v.flowRate))
		for _, ov := range v.paths {
			file.WriteString(fmt.Sprintf("%s -> %s\n", v.name, ov.name))
		}
		file.WriteString("\n")
	}
	file.WriteString("}\n")

	cmd := exec.Command("dot", "-Tpng", "-o", file.Name()+".png", file.Name())
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Couldn't render %s: %v", file.Name(), err)
		return
	}
	//os.Remove(file.Name())
	if save {
		log.Print("Image is at ", file.Name()+".png")
	} else {
		cmd = exec.Command("eog", file.Name()+".png")
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Couldn't display %s: %v", file.Name()+".png", err)
			return
		}
		os.Remove(file.Name() + ".png")
	}
}

func DjikstraStart() map[string]int {
	rv := map[string]int{}
	for n := range Valves {
		rv[n] = math.MaxInt
	}
	return rv
}

func IDKey(vl []*Valve) string {
	ids := []string{}
	for _, v := range vl {
		ids = append(ids, v.name)
	}
	return strings.Join(ids, "")
}

var hit int
var miss int

type CE struct {
	closed   []string
	pressure int
}

func Next(at *Valve, tick int, limit int, closed []*Valve) ([]string, int) {
	key := fmt.Sprintf("%s-%d-%d-%s", at, tick, limit, IDKey(closed))
	if rv, ok := Cache[key]; ok {
		hit++
		return rv.closed, rv.pressure
	}
	miss++
	//fmt.Println(tick, "@", at, " hit=", hit, " miss=", miss)

	release := 0
	tick2 := tick
	if at.flowRate > 0 {
		release = at.Potential(tick, limit)
		tick2++
	}

	best := 0
	var bestClosed []string
	for _, v := range closed {
		when := tick2 + at.dists[v.name]
		if when > limit-1 {
			//fmt.Println(tick, "@", at, ": can't reach ", v, " in time.")
			continue
		}
		remain := []*Valve{}
		for _, vv := range closed {
			if vv != v {
				remain = append(remain, vv)
			}
		}
		d, p := Next(v, when, limit, remain)
		//fmt.Println(tick, "@", at, ": ", when, "@", v, "+=", p)
		if p > best {
			best = p
			bestClosed = d
		}
	}
	newClosed := []string{at.name}
	if best > 0 {
		newClosed = append(newClosed, bestClosed...)
	}
	Cache[key] = CE{closed: newClosed, pressure: release + best}
	return newClosed, release + best
}

var INPUT_RE = regexp.MustCompile(`Valve (.*?) has flow rate=(\d+); tunnel.*? lead.*? to valve.*? (.*)$`)

var Valves map[string]*Valve
var Priority []*Valve
var Cache map[string]CE

func main() {
	s := bufio.NewScanner(os.Stdin)

	Valves = map[string]*Valve{}
	Priority = []*Valve{}
	Cache = map[string]CE{}

	for s.Scan() {
		m := INPUT_RE.FindStringSubmatch(s.Text())
		if len(m) != 4 {
			log.Fatalf("Couldn't parse: %s", s.Text())
		}
		v := GetOrCreateValve(m[1])
		if v.flowRate != -1 {
			log.Fatal("multiple flow definitions of ", v)
		}
		v.flowRate = Int(m[2])
		v.paths = GetOrCreateValves(m[3])
		//fmt.Println(v.Describe())
		if v.flowRate != 0 {
			Priority = append(Priority, v)
		}
	}
	sort.Slice(Priority, func(i, j int) bool {
		return Priority[i].flowRate > Priority[j].flowRate
	})

	// Precompute how far from each other the valves are.
	//fmt.Println(len(Priority), " useful valves ", Priority)
	start := Valves["AA"]
	for _, v := range Priority {
		start.CalcDist(v)
		//fmt.Println(start, " to ", v, " is ", start.dists[v.name])
		for _, vd := range Priority {
			if vd == v {
				continue
			}
			v.CalcDist(vd)
			//fmt.Println(v, " to ", vd, " is ", v.dists[vd.name])
		}
	}

	mypath, mypressure := Next(start, 0, 26, Priority)
	closed := []*Valve{}
NEXT:
	for _, v := range Priority {
		for _, c := range mypath {
			if v.name == c {
				continue NEXT
			}
		}
		closed = append(closed, v)
	}
	fmt.Println(closed)
	elepath, elepressure := Next(start, 0, 26, closed)
	fmt.Println(" Me: ", strings.Join(mypath, " > "), mypressure)
	fmt.Println("Ele: ", strings.Join(elepath, " > "), elepressure)
	fmt.Println(mypressure + elepressure)
}

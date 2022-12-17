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

func Min(a, b int) int {
	if a == -1 {
		return b
	}
	if b == -1 {
		return a
	}
	if a < b {
		return a
	}
	return b
}

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
	desc     string
	pressure int
}

func Prefix(d int) string {
	//if d > 0 {
	//		return strings.Repeat(" ", d)
	//	}
	return ""
}

type ActorState struct {
	who  string
	at   *Valve
	when int
	d    int
}

func (as ActorState) String() string {
	return fmt.Sprintf("%s:%s@%d", as.who, as.at.name, as.when)
}

func (as ActorState) Prefix() string {
	return Prefix(as.d)
}

// for when only actor is moving
func FindBest(tick int, limit int, closed []*Valve, actor ActorState, other ActorState) (string, int) {
	//fmt.Println(actor.Prefix(), tick, "@", actor.at, " options are: ", IDKey(closed))
	best := 0
	bestd := ""
	explored := 0
	for _, v := range closed {
		when := tick + actor.at.dists[v.name]
		if when > limit-1 {
			//fmt.Println(actor, ": can't reach ", v, " in time.")
			continue
		}
		remain := []*Valve{}
		for _, vv := range closed {
			if vv != v {
				remain = append(remain, vv)
			}
		}
		next := ActorState{who: actor.who, when: when, at: v, d: actor.d + 1}
		d, p := Next(Min(when, other.when), limit, remain, next, other)
		//fmt.Println(actor, ": (", v, "@", when, " + ", other, ") +=", p)
		if p > best {
			best = p
			bestd = d
		}
		explored++
	}
	if explored == 0 && tick < other.when {
		// actor couldn't do anything at this step (out of option), but other might have a valve to close still!
		bestd, best = Next(other.when, limit, closed, actor, other)
	}

	return bestd, best
}

func RemoveValve(l []*Valve, v *Valve) []*Valve {
	newL := []*Valve{}
	for _, vv := range l {
		if vv != v {
			newL = append(newL, vv)
		}
	}
	return newL
}

func Next(tick int, limit int, closed []*Valve, a1 ActorState, a2 ActorState) (string, int) {
	key := fmt.Sprintf("%s-%s-%d-%d-%s", a1, a2, tick, limit, IDKey(closed))
	// TODO: is key ordering between a1, a2 important?
	if rv, ok := Cache[key]; ok {
		hit++
		return rv.desc, rv.pressure
	}
	miss++
	//fmt.Println(tick, ":", a1, ", ", a2, "; hit=", hit, " miss=", miss)

	release := 0
	dl := []string{}
	a1_nexttick := tick
	if a1.when == tick && a1.at.flowRate > 0 {
		release += a1.at.Potential(tick, limit)
		a1_nexttick++
		//fmt.Println(a1, " opens ", a1.at.name)
		dl = append(dl, fmt.Sprintf("%s:%s", a1.at.name, a1.who))
		closed = RemoveValve(closed, a1.at)
	}
	a2_nexttick := tick
	if a2.when == tick && a2.at.flowRate > 0 {
		release += a2.at.Potential(tick, limit)
		a2_nexttick++
		//fmt.Println(a2, " opens ", a2.at.name)
		dl = append(dl, fmt.Sprintf("%s:%s", a2.at.name, a2.who))
		closed = RemoveValve(closed, a2.at)
	}

	best := 0
	bestd := ""
	if a1.when == tick && a2.when != tick {
		bestd, best = FindBest(a1_nexttick, limit, closed, a1, a2)
	} else if a2.when == tick && a1.when != tick {
		bestd, best = FindBest(a2_nexttick, limit, closed, a2, a1)
	} else if a1.when == tick && a2.when == tick {
		for _, v := range closed {
			when := a1_nexttick + a1.at.dists[v.name]
			if when > limit-1 {
				//fmt.Println(a1, ": can't reach ", v, " in time.")
				continue
			}
			remain := RemoveValve(closed, v)
			a1next := ActorState{who: a1.who, when: when, at: v, d: a1.d + 1}
			for _, vv := range remain {
				when2 := a2_nexttick + a2.at.dists[vv.name]
				if when2 > limit-1 {
					//fmt.Println(a2, ": can't reach ", v, " in time.")
					continue
				}
				remain2 := RemoveValve(remain, vv)
				a2next := ActorState{who: a2.who, when: when2, at: vv, d: a2.d + 1}
				d, p := Next(Min(when, when2), limit, remain2, a1next, a2next)
				//fmt.Println(a1, ": (", a1next, " THEN ", a2next, ") +=", p)
				if a1.d == 0 {
					fmt.Println("R: ", d, p)
				}
				if p > best {
					best = p
					bestd = d
				}
			}
		}
	}

	d := ""
	if len(dl) > 0 {
		d += fmt.Sprintf("%d:(%s)", tick+1, strings.Join(dl, ", "))
	}
	if len(dl) > 0 && best > 0 {
		d += " > "
	}
	if best > 0 {
		d += bestd
	}
	Cache[key] = CE{desc: d, pressure: release + best}
	return d, release + best
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

	path, pressure := Next(0, 26, Priority, ActorState{who: " me", when: 0, at: start}, ActorState{who: "ele", when: 0, at: start})
	fmt.Println(path)
	fmt.Println(pressure)
}

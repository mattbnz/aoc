// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 16, Puzzle 2.
// Proboscidea Volcanium - double pressure release.

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
	// prio (could be dist for djikstra, or tick for main loop)
	prio int

	// data
	node string // used in djikstra
	dp   DP     // used in mainloop
}

// Priority Queue
type PQ struct {
	Q []PQE
}

func (q *PQ) Add(e PQE) {
	if e.dp.who != "" {
		fmt.Println("Add", e.prio, ": ", e.dp.who, "@", e.dp.at)
	}
	if len(q.Q) == 0 {
		q.Q = append(q.Q, e)
		return
	}
	if e.prio < q.Q[0].prio {
		q.Q = append([]PQE{e}, q.Q...)
		return
	}
	for n, t := range q.Q {
		if e.prio < t.prio {
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
	q.Add(PQE{node: v.name, prio: 0})
	for len(q.Q) > 0 {
		c := q.Pop()
		valve := Valves[c.node]
		been[valve.name] = true
		for _, n := range valve.paths {
			if been[n.name] {
				continue
			}
			if c.prio+1 < dists[n.name] {
				dists[n.name] = c.prio + 1
				q.Add(PQE{node: n.name, prio: c.prio + 1})
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

func Next(at *Valve, tick int, limit int, closed []*Valve) (string, int) {
	key := fmt.Sprintf("%s-%d-%d-%s", at, tick, limit, IDKey(closed))
	if rv, ok := Cache[key]; ok {
		hit++
		return rv.desc, rv.pressure
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
	bestd := ""
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
			bestd = d
		}
	}
	d := at.name
	if best > 0 {
		d += " > " + bestd
	}
	Cache[key] = CE{desc: d, pressure: release + best}
	return d, release + best
}

type OpenLog struct {
	who  string
	when int
	// only populated in final loop
	valve string
}

type PathLog map[string]OpenLog

func (pl PathLog) String() string {
	opens := []OpenLog{}
	for v, l := range pl {
		l.valve = v
		opens = append(opens, l)
	}
	sort.Slice(opens, func(i, j int) bool { return opens[i].when < opens[j].when })
	rv := []string{}
	for _, o := range opens {
		if o.when > 0 {
			rv = append(rv, fmt.Sprintf("%s by %s@%d", o.valve, o.who, o.when))
		}
	}
	return strings.Join(rv, " > ") //+ " == " + fmt.Sprintf("%d", pl.Score(26))
}

func (pl PathLog) Score(limit int) int {
	score := 0
	for v, l := range pl {
		if l.when > 0 {
			score += Valves[v].Potential(l.when, limit)
		}
	}
	return score
}

// makes a copy of a log
func CopyLog(in PathLog) PathLog {
	rv := PathLog{}
	for v, log := range in {
		rv[v] = log
	}
	return rv
}

// Decision point
type DP struct {
	// Who this decision is for
	who string
	// valve decision needs to be made at
	at *Valve
	// pointer to current valve state
	opened *PathLog
	// pointer to final result accumulator
	results *[]PathLog
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

	limit := 26
	results := []PathLog{}
	opened := PathLog{}
	for _, v := range Priority {
		opened[v.name] = OpenLog{} // default val = closed, when=0
	}

	q := PQ{}
	q.Add(PQE{prio: 0, dp: DP{who: " me", at: start, opened: &opened, results: &results}})
	q.Add(PQE{prio: 0, dp: DP{who: "ele", at: start, opened: &opened, results: &results}})
	for len(q.Q) > 0 {
		e := q.Pop()
		tick := e.prio
		fmt.Println(e.prio, ":", e.dp.who, "@", e.dp.at, ":: ", e.dp.opened)
		// open if needed
		tick2 := tick
		if e.dp.at.flowRate > 0 {
			log := *e.dp.opened
			log[e.dp.at.name] = OpenLog{who: e.dp.who, when: tick}
			tick2++
		}

		// pick next
		explored := 0
		for v, vlog := range *e.dp.opened {
			if vlog.when > 0 || vlog.when == -1 {
				continue // valve already open (or enroute to be opened)
			}
			when := tick2 + e.dp.at.dists[v]
			if when > limit-1 {
				//fmt.Println(tick, "@", at, ": can't reach ", v, " in time.")
				continue
			}
			(*e.dp.opened)[e.dp.at.name] = OpenLog{who: e.dp.who, when: -1}
			newLog := CopyLog(*e.dp.opened)
			q.Add(PQE{prio: when, dp: DP{who: e.dp.who, at: Valves[v], opened: &newLog, results: e.dp.results}})
			explored++
		}

		// or if nothing further to explore accumulate result (using a copy as it's final)
		if explored == 0 {
			*e.dp.results = append(*e.dp.results, CopyLog(*e.dp.opened))
		}
	}
	fmt.Println("Found ", len(results), " possible paths!")
	bestScore := 0
	best := PathLog{}
	for _, path := range results {
		score := path.Score(limit)
		if score > bestScore {
			best = path
			bestScore = score
		}
	}
	fmt.Println(best)
	fmt.Println(bestScore)
}

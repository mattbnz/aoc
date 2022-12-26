// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 24, Puzzle 1.
// Blizzard Basin. Another search...

package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Abs(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}

func Max(a, b int) int {
	if a == -1 {
		return b
	}
	if b == -1 {
		return a
	}
	if a > b {
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

type Pos struct {
	row, col int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d,%d", p.row, p.col)
}

type Blizzard int

const (
	LEFT Blizzard = iota
	RIGHT
	UP
	DOWN
)

func (b Blizzard) String() string {
	if b == LEFT {
		return "<"
	} else if b == RIGHT {
		return ">"
	} else if b == UP {
		return "^"
	} else if b == DOWN {
		return "v"
	}
	return "!"
}

func NewBlizzard(s string) (Blizzard, error) {
	if s == "<" {
		return LEFT, nil
	} else if s == ">" {
		return RIGHT, nil
	} else if s == "^" {
		return UP, nil
	} else if s == "v" {
		return DOWN, nil
	}
	return LEFT, fmt.Errorf("%s is not a blizzard!", s)
}

type Cell struct {
	exists bool
	W      bool              // wall
	E      bool              // expedition presence
	B      map[Blizzard]bool // blizard presence
}

func (c Cell) String() string {
	if c.W {
		return "#"
	}
	if c.E {
		return "E"
	}
	bs := ""
	for b, p := range c.B {
		if !p {
			continue
		}
		bs += b.String()
	}
	if len(bs) == 0 {
		return "."
	}
	if len(bs) > 1 {
		return fmt.Sprintf("%d", len(bs))
	}
	return bs
}

func (c Cell) Open() bool {
	if c.W {
		return false
	}
	for _, b := range c.B {
		if b {
			return false
		}
	}
	return true
}

func NewCell(s string) Cell {
	c := Cell{exists: true, B: map[Blizzard]bool{}}
	if s == "#" {
		c.W = true
		return c
	}
	b, err := NewBlizzard(s)
	if err == nil {
		c.B[b] = true
	}
	return c
}

const NOW = -1

type Grid struct {
	c              map[Pos]map[int]Cell
	maxrow, maxcol int
	epath          map[int]Pos

	minute int
}

func (g Grid) String() string {
	return fmt.Sprintf("Grid %s @ minute %d", g.Key(), g.minute)
}

func (g Grid) Key() string {
	s := ""
	for row := 2; row < g.maxrow; row++ {
		for col := 2; col < g.maxcol; col++ {
			s += g.C(Pos{row, col}, NOW).String()
		}
	}
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))[:12]
}

// Returns a copy of this grid at the current minute (doesn't preserve past minutes)
func (g Grid) Copy() Grid {
	ng := Grid{
		maxrow: g.maxrow,
		maxcol: g.maxcol,
		minute: g.minute,
		c:      map[Pos]map[int]Cell{},
		epath:  map[int]Pos{},
	}
	for p, mins := range g.c {
		ng.c[p] = map[int]Cell{g.minute: mins[g.minute]}
	}
	for m, p := range g.epath {
		ng.epath[m] = p
	}
	return ng
}

func (g Grid) Prefix() string {
	return "Grid " + g.Key() + ", "
}

func (g Grid) Print(prefix string, suffix string) {
	fmt.Printf("%sminute %d: %s\n", prefix, g.minute, suffix)
	for row := 1; row <= g.maxrow; row++ {
		for col := 1; col <= g.maxcol; col++ {
			fmt.Print(g.C(Pos{row, col}, NOW))
		}
		fmt.Println()
	}
	fmt.Println()
}

func (g Grid) Path(lastN int) string {
	start := 0
	rv := []string{}
	if lastN > 0 && len(g.epath) > lastN {
		start = Max(start, g.minute-lastN)
		rv = append(rv, "...")
	}
	for i := start; i <= g.minute; i++ {
		p, ok := g.epath[i]
		if ok {
			rv = append(rv, p.String())
		} else {
			rv = append(rv, "??")
		}
	}
	return strings.Join(rv, " > ")
}

func (g Grid) C(p Pos, min int) Cell {
	if min == NOW {
		min = g.minute
	}
	return g.c[p][min]
}

func (g *Grid) SetC(p Pos, min int, v Cell) {
	if min == NOW {
		min = g.minute
	}
	if _, found := g.c[p]; found {
		g.c[p][min] = v
	} else {
		g.c[p] = map[int]Cell{min: v}
	}
	g.maxrow = Max(g.maxrow, p.row)
	g.maxcol = Max(g.maxcol, p.col)
	if v.E && min >= len(g.epath) {
		g.epath[min] = p
	}
}

// Make a copy of every cell, with the expedition and blizzards removed
func (g *Grid) PrepMinute(min int) {
	for p, mins := range g.c {
		if p.row == 1 || p.row == g.maxrow || p.col == 1 || p.col == g.maxcol {
			c := mins[g.minute]
			c.E = false
			g.SetC(p, min, c) // Wall, copy but clear expedition
		} else {
			g.SetC(p, min, NewCell(""))
		}
	}
}

// Increments a minute, aka moves all the blizzards
// Does not move the expedition.
func (g *Grid) Inc() {
	next := g.minute + 1
	g.PrepMinute(next)

	for row := 1; row <= g.maxrow; row++ {
		for col := 1; col <= g.maxcol; col++ {
			p := Pos{row, col}
			c := g.C(p, NOW)
			if !c.exists || c.W {
				continue
			}
			for b, here := range c.B {
				if !here {
					continue
				}
				np := g.IncBlizzard(p, b)
				nc := g.C(np, next)
				nc.B[b] = true
				g.SetC(np, next, nc)
			}
		}
	}
	g.minute = next
}

// Returns next pos of blizzard moving in Blizzard dir from p
func (g *Grid) IncBlizzard(p Pos, b Blizzard) Pos {
	np := Pos{p.row, p.col}
	if b == LEFT {
		np.col--
	} else if b == RIGHT {
		np.col++
	} else if b == UP {
		np.row--
	} else if b == DOWN {
		np.row++
	}
	if np.row == g.maxrow {
		np.row = 2 // 1 == wall
	}
	if np.row == 1 {
		np.row = g.maxrow - 1
	}
	if np.col == g.maxcol {
		np.col = 2
	}
	if np.col == 1 {
		np.col = g.maxcol - 1
	}
	return np
}

func (g *Grid) FindE(min int) Pos {
	if min == NOW {
		min = g.minute
	}
	for p, mins := range g.c {
		if mins[min].E {
			return p
		}
	}
	log.Fatalf("Expedition has gone missing in minute %d", g.minute)
	return Pos{}
}

func (g *Grid) DistFrom(min int, p2 Pos) int {
	p := g.FindE(min)
	return Abs(p2.row-p.row) + Abs(p2.col-p.col)
}

// Returns possible positions for the expedition in this minute
func (g *Grid) Moves() []Pos {
	c := g.FindE(g.minute - 1)
	rv := []Pos{}
	for _, p := range []Pos{c, {c.row - 1, c.col}, {c.row + 1, c.col}, {c.row, c.col - 1}, {c.row, c.col + 1}} {
		t := g.C(p, -1)
		if t.exists && t.Open() {
			rv = append(rv, p)
		}
	}
	return rv
}

func NewGrid() Grid {
	g := Grid{c: map[Pos]map[int]Cell{}, epath: map[int]Pos{}}
	g.maxrow = -1
	g.maxcol = -1
	return g
}

// A priority Queue entry.
type PQE struct {
	grid  Grid
	speed float64
}

// Priority Queue
type PQ struct {
	Q []PQE

	lastSpeeds []float64
	start      Pos
}

func (q *PQ) Add(g Grid) PQE {
	qe := PQE{grid: g}
	qe.speed = float64(g.DistFrom(NOW, q.start)) / float64(g.minute)
	if q.lastSpeeds == nil {
		q.lastSpeeds = make([]float64, 10)
	}
	q.lastSpeeds = append(q.lastSpeeds[1:], qe.speed)
	if len(q.Q) == 0 {
		q.Q = append(q.Q, qe)
		return qe
	}
	if qe.speed > q.Q[0].speed {
		q.Q = append([]PQE{qe}, q.Q...)
		return qe
	}
	for n, t := range q.Q {
		if qe.speed > t.speed {
			q.Q = append(q.Q[:n+1], q.Q[n:]...)
			q.Q[n] = qe
			return qe
		}
	}
	q.Q = append(q.Q, qe)
	return qe
}

func (q *PQ) Pop() PQE {
	e := q.Q[0]
	q.Q = q.Q[1:]
	return e
}

func (q *PQ) AvgSpeed() float64 {
	sum := 0.0
	for _, e := range q.lastSpeeds {
		sum += e
	}
	return sum / float64(len(q.lastSpeeds))
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	startGrid := NewGrid()

	row := 1
	for s.Scan() {
		for col, c := range s.Text() {
			p := Pos{row, col + 1}
			c := NewCell(string(c))
			if row == 1 && c.Open() {
				c.E = true
			}
			startGrid.SetC(p, NOW, c)
		}
		row++
	}
	exit := Pos{}
	for c := 1; c <= startGrid.maxcol; c++ {
		p := Pos{startGrid.maxrow, c}
		if startGrid.C(p, NOW).Open() {
			exit = p
		}
	}
	if exit.row == 0 {
		log.Fatal("Didn't find exit!")
	}
	startGrid.Print(startGrid.Prefix(), "")
	start := startGrid.FindE(NOW)
	fmt.Printf("Expedition is at %s\n", start)
	fmt.Printf("Exit is at %s, %d away\n\n", exit, startGrid.DistFrom(NOW, exit))

	q := PQ{start: start}
	q.Add(startGrid)
	seen := map[string]bool{}
	var best *Grid
	var maxDist int
	for len(q.Q) > 0 {
		qe := q.Pop()
		g := qe.grid
		k := g.Key()
		if _, done := seen[k]; done {
			//fmt.Printf("Skipping %s - already done\n", g)
			continue
		}
		prefix := g.Prefix()
		seen[k] = true
		if best != nil && g.minute > best.minute {
			fmt.Printf("%sGiving up on > best minute %d\n", prefix, best.minute)
			continue
		}
		dist := g.DistFrom(NOW, start)
		if dist > maxDist {
			fmt.Printf("%sNew max distance of %d\n", prefix, dist)
			maxDist = dist
		} else if dist < maxDist/2 {
			//g.Print(prefix, "Giving up - not making progress")
			continue
		}
		if qe.speed < q.AvgSpeed()/2 {
			g.Print(prefix, "Giving up, slow")
			continue
		}
		fmt.Println(prefix, fmt.Sprintf("spd=%.2f", qe.speed), g.Path(15))
		at := g.FindE(NOW)
		g.Inc()
		moves := g.Moves()
		for _, m := range moves {
			if m == exit {
				fmt.Printf("%s Reached exit at minute %d!\n", prefix, g.minute)
				if best == nil || best.minute > g.minute {
					fmt.Println(" - and it's our new best!", g.minute)
					best = &g
				}
				g.Print(prefix, "Exit to "+m.String())
				continue
			}
			nextGrid := g.Copy()
			nc := g.C(m, NOW)
			nc.E = true
			nextGrid.SetC(m, NOW, nc)
			qe := q.Add(nextGrid)
			//nextGrid.Print(prefix, fmt.Sprintf("%s move to %s  (%s), ql=%d", at, m, nextGrid.Key(), len(q)))
			fmt.Println(prefix,
				fmt.Sprintf("%s move to %s (%s, spd=%.2f), ql=%d",
					at, m, nextGrid.Key(), qe.speed, len(q.Q)))
		}
	}

	if best != nil {
		fmt.Printf("Reached exit in %d minutes\n", best.minute)
	} else {
		log.Fatal("Did not reach the exit!")
	}
}

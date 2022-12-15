// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 15, Puzzle 1.
// Beacon Exclusion Zone.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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

func Abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

type Pos struct {
	row, col int
}

func (p Pos) String() string {
	return fmt.Sprintf("x=%d,y=%d", p.col, p.row)
}

func (p Pos) Dist(o Pos) int {
	return Abs(p.col-o.col) + Abs(p.row-o.row)
}

func NewPos(s string) Pos {
	col, row, ok := strings.Cut(s, ",")
	if !ok {
		log.Fatal("bad pos: ", s)
	}
	return Pos{Int(row), Int(col)}
}

const AIR = 0
const SENSOR = 1
const BEACON = 2
const COVERED = 3

type Grid struct {
	c              map[Pos]int
	minrow, maxrow int
	mincol, maxcol int

	maxdrop int
}

func (g Grid) Print() {
	g.print(g.maxrow)
}

func (g Grid) PrintPartial() {
	g.print(g.maxdrop + 1)
}

func (g Grid) print(maxrow int) {
	/*	colminS := fmt.Sprintf("% 3d", g.mincol)
		colmaxS := fmt.Sprintf("% 3d", g.maxcol)
		fmt.Printf("      %c%s%c\n", colminS[0], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[0])
		fmt.Printf("      %c%s%c\n", colminS[1], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[1])
		fmt.Printf("      %c%s%c\n", colminS[2], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[2])
		fmt.Printf("      %c%s%c\n", colminS[3], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[3])
	*/
	for row := g.minrow; row <= g.maxrow; row++ {
		if row > maxrow {
			return
		}
		fmt.Printf("% 10d: ", row)
		for col := g.mincol; col <= g.maxcol; col++ {
			c := g.C(Pos{row, col})
			if c == SENSOR {
				fmt.Printf("S")
			} else if c == BEACON {
				fmt.Printf("B")
			} else if c == COVERED {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
}

func (g Grid) C(p Pos) int {
	return g.c[p] // default val if missing == AIR
}

func (g *Grid) SetC(p Pos, v int) {
	g.c[p] = v
	g.minrow = Min(g.minrow, p.row)
	g.mincol = Min(g.mincol, p.col)
	g.maxrow = Max(g.maxrow, p.row)
	g.maxcol = Max(g.maxcol, p.col)
}

func (g *Grid) Add(s Sensor) {
	g.SetC(s.at, SENSOR)
	g.SetC(s.closest, BEACON)

	scope := s.at.Dist(s.closest)
	DPrintf("%s and %s has scope %d\n", s.at, s.closest, scope)
	for row := s.at.row - scope; row != s.at.row+scope+1; row++ {
		for col := s.at.col - scope; col != s.at.col+scope+1; col++ {
			p := Pos{row: row, col: col}
			if p == s.at || p == s.closest {
				continue
			}
			if g.C(p) != AIR {
				continue
			}
			if s.at.Dist(p) <= scope {
				g.SetC(p, COVERED)
			}
		}
	}
}

var DBG = true

func DPrintln(a ...any) {
	if DBG {
		fmt.Println(a...)
	}
}

func DPrintf(f string, args ...any) {
	if DBG {
		fmt.Printf(f, args...)
	}
}

func NewGrid() Grid {
	g := Grid{c: map[Pos]int{}}
	g.minrow = -1
	g.maxrow = -1
	g.mincol = -1
	g.maxcol = -1
	return g
}

type Sensor struct {
	at      Pos
	closest Pos
}

func (s Sensor) String() string {
	return fmt.Sprintf("Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", s.at.col, s.at.row, s.closest.col, s.closest.row)
}

// var INPUT_RE = regexp.MustCompile(`Sensor at x=([-\d+]), y=([-\d+]): closest beacon is at x=([-\d+]), y=([-\d+])`)
var INPUT_RE = regexp.MustCompile(`Sensor at x=([-\d]+), y=([-\d]+): closest beacon is at x=([-\d]+), y=([-\d]+)`)

func main() {
	s := bufio.NewScanner(os.Stdin)

	g := NewGrid()

	for s.Scan() {
		m := INPUT_RE.FindStringSubmatch(s.Text())
		if len(m) != 5 {
			fmt.Println(len(m), m)
			log.Fatalf("Couldn't parse: %s", s.Text())
		}
		s := Sensor{}
		s.at = Pos{col: Int(m[1]), row: Int(m[2])}
		s.closest = Pos{col: Int(m[3]), row: Int(m[4])}
		g.Add(s)
	}
	g.Print()

	sum := 0
	for col := g.mincol; col != g.maxcol+1; col++ {
		c := g.C(Pos{row: 10, col: col})
		if c == COVERED || c == SENSOR {
			sum++
		}
	}
	fmt.Println(sum)
}

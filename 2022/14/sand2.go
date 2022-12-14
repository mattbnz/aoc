// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 14, Puzzle 2.
// Regolith Reservoir - Falling sand.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

type Pos struct {
	row, col int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d,%d", p.row, p.col)
}

func NewPos(s string) Pos {
	col, row, ok := strings.Cut(s, ",")
	if !ok {
		log.Fatal("bad pos: ", s)
	}
	return Pos{Int(row), Int(col)}
}

const AIR = 0
const ROCK = 1
const SAND = 2

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
	colminS := fmt.Sprintf("% 3d", g.mincol)
	colmaxS := fmt.Sprintf("% 3d", g.maxcol)
	fmt.Printf("      %c%s%c\n", colminS[0], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[0])
	fmt.Printf("      %c%s%c\n", colminS[1], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[1])
	fmt.Printf("      %c%s%c\n", colminS[2], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[2])
	fmt.Printf("      %c%s%c\n", colminS[3], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[3])
	for row := g.minrow; row <= g.maxrow; row++ {
		if row > maxrow {
			return
		}
		fmt.Printf("% 4d: ", row)
		for col := g.mincol; col <= g.maxcol; col++ {
			c := g.C(Pos{row, col})
			if row == 0 && col == (500-g.mincol) {
				fmt.Printf("+")
			} else if c == SAND {
				fmt.Printf("o")
			} else if c == ROCK {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
}

func (g Grid) C(p Pos) int {
	if p.row == g.maxrow {
		return ROCK // infinite floor
	}
	return g.c[p] // default val if missing == AIR
}

func (g *Grid) SetC(p Pos, v int) {
	g.c[p] = v
	g.minrow = Min(g.minrow, p.row)
	g.mincol = Min(g.mincol, p.col)
	g.maxrow = Max(g.maxrow, p.row)
	g.maxcol = Max(g.maxcol, p.col)
}

func (g *Grid) AddLines(l [][]Pos) {
	for _, line := range l {
		col, row := line[0].col, line[0].row
		e := line[1]
		for col != e.col || row != e.row {
			g.SetC(Pos{row, col}, ROCK)
			if col < e.col {
				col++
			} else if col > e.col {
				col--
			} else if row < e.row {
				row++
			} else if row > e.row {
				row--
			} else {
				log.Fatal("Bad line", line[0], line[1])
			}
		}
		g.SetC(Pos{row, col}, ROCK)
	}
	g.maxrow += 2
}

var DBG = false

func DPrintln(a ...any) {
	if DBG {
		fmt.Println(a...)
	}
}

// returns true if the sand came to rest; false if it fell to eternity
func (g *Grid) DropSand() bool {
	col, row := 500, 0
	for row < g.maxrow {
		// basic case, can drop
		if g.C(Pos{row + 1, col}) == AIR {
			row++
			DPrintln("dropping to ", col, row)
			continue
		}
		// try left
		if g.C(Pos{row + 1, col - 1}) == AIR {
			row++
			col--
			DPrintln("dropping left to ", col, row)
			continue
		}
		// try right
		if g.C(Pos{row + 1, col + 1}) == AIR {
			row++
			col++
			DPrintln("dropping right to ", col, row)
			continue
		}
		// can't move
		g.SetC(Pos{row, col}, SAND)
		g.maxdrop = Max(g.maxdrop, row)
		DPrintln("settled at ", col, row, " maxdrop ", g.maxdrop)
		if col == 500 && row == 0 {
			return false
		}
		return true
	}

	// fell out
	DPrintln("out of bounds at ", col, row)
	return false
}

func NewGrid() Grid {
	g := Grid{c: map[Pos]int{}}
	g.minrow = 0
	g.maxrow = -1
	g.mincol = -1
	g.maxcol = -1
	return g
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	lines := [][]Pos{}

	for s.Scan() {
		l := strings.Split(s.Text(), " -> ")
		for i := 1; i < len(l); i++ {
			lines = append(lines, []Pos{NewPos(l[i-1]), NewPos(l[i])})
		}
	}
	grid := NewGrid()
	grid.AddLines(lines)
	grid.Print()

	sand := 0
	for grid.DropSand() {
		sand++
		//grid.PrintPartial()
	}
	fmt.Println()
	grid.Print()
	fmt.Println(sand + 1)
}

// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 13, Puzzle 1.
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
	c              [][]int
	minrow, maxrow int
	mincol, maxcol int
}

func (g Grid) Print() {
	colminS := fmt.Sprintf("% 3d", g.mincol)
	colmaxS := fmt.Sprintf("% 3d", g.maxcol)
	fmt.Printf("    %c%s%c\n", colminS[0], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[0])
	fmt.Printf("    %c%s%c\n", colminS[1], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[1])
	fmt.Printf("    %c%s%c\n", colminS[2], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[2])
	fmt.Printf("    %c%s%c\n", colminS[3], strings.Repeat(" ", g.maxcol-g.mincol-1), colmaxS[3])
	for row, cols := range g.c {
		fmt.Printf("% 2d: ", row)
		for _, c := range cols {
			if c == SAND {
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

func (g Grid) AddLines(l [][]Pos) {
	for _, line := range l {
		col, row := line[0].col, line[0].row
		e := line[1]
		for col != e.col || row != e.row {
			g.c[row-g.minrow][col-g.mincol] = ROCK

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
		g.c[row-g.minrow][col-g.mincol] = ROCK
	}
}

func NewGrid(minrow, mincol, maxrow, maxcol int) Grid {
	g := Grid{}
	for r := minrow; r <= maxrow; r++ {
		g.c = append(g.c, make([]int, maxcol-mincol+1))
	}
	g.minrow = minrow
	g.maxrow = maxrow
	g.mincol = mincol
	g.maxcol = maxcol
	return g
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	lines := [][]Pos{}
	minrow, maxrow := 0, -1
	mincol, maxcol := -1, -1

	for s.Scan() {
		l := strings.Split(s.Text(), " -> ")
		for i := 1; i < len(l); i++ {
			line := []Pos{NewPos(l[i-1]), NewPos(l[i])}
			minrow = Min(Min(minrow, line[0].row), line[1].row)
			mincol = Min(Min(mincol, line[0].col), line[1].col)
			maxrow = Max(Max(maxrow, line[0].row), line[1].row)
			maxcol = Max(Max(maxcol, line[0].col), line[1].col)
			lines = append(lines, line)
		}
	}
	grid := NewGrid(minrow, mincol, maxrow, maxcol)
	grid.AddLines(lines)
	grid.Print()
}

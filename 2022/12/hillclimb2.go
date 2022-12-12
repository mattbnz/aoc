// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 12, Puzzle 2.
// Hill Climbing Algorithm - Best Path.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Abs(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}

type Pos struct {
	row, col int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d,%d", p.row, p.col)
}

func (p Pos) DistFrom(op Pos) int {
	return Abs(op.row-p.row) + Abs(op.col-p.col)
}

// Generated positions may not be valid...
func (p Pos) Left() Pos {
	return Pos{row: p.row, col: p.col - 1}
}
func (p Pos) Right() Pos {
	return Pos{row: p.row, col: p.col + 1}
}
func (p Pos) Up() Pos {
	return Pos{row: p.row - 1, col: p.col}
}
func (p Pos) Down() Pos {
	return Pos{row: p.row + 1, col: p.col}
}

type Cell struct {
	id     Pos
	height int

	dTo   int // distance to this Cell (used to avoid backtracking)
	dFrom int // distance from this Cell (the answer we want)
}

var Reset = "\033[0m"

var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"
var Colors = []string{Red, Green, Yellow, Blue, Purple, Cyan, Gray}

func (c Cell) String() string {
	return fmt.Sprintf("%s(%c)", c.id, c.height+'a')
}

func (c Cell) PrintColorHeight() {
	fmt.Printf(Colors[c.height%len(Colors)])
	fmt.Printf("%c", c.height+'a')
}

type Grid struct {
	cells [][]*Cell

	start Pos
	end   Pos
}

func (g *Grid) ParseRow(r string) {
	row := []*Cell{}
	for i, cS := range r {
		c := &Cell{id: Pos{len(g.cells), i}}
		if cS == 'S' {
			g.start = c.id
			c.height = 0
		} else if cS == 'E' {
			g.end = c.id
			c.height = 25
		} else if cS >= 'a' && cS <= 'z' {
			c.height = int(cS - 'a')
		} else {
			log.Fatalf("bad row: %s", r)
		}
		row = append(row, c)
	}
	g.cells = append(g.cells, row)
}

func (g *Grid) Print() {
	for r, cols := range g.cells {
		for c, h := range cols {
			if r == g.start.row && c == g.start.col {
				fmt.Print(Reset)
				fmt.Printf("S")
			} else if r == g.end.row && c == g.end.col {
				fmt.Print(Reset)
				fmt.Printf("E")
			} else {
				h.PrintColorHeight()
			}
		}
		fmt.Println("")
	}
	fmt.Print(Reset)
}

// returns a cell from the grid (or nil, if pos is invalid)
func (g *Grid) Cell(p Pos) *Cell {
	if p.row < 0 || p.row > len(g.cells)-1 || p.col < 0 || p.col > len(g.cells[0])-1 {
		return nil
	}
	return g.cells[p.row][p.col]
}

// Returns destination cell if the move is legal, or nil if not.
// - If it is, returns the current distance of to from the end (may be zero if not yet known)
// - If it's not, returns -1
func (g *Grid) CanMove(pFrom, pTo Pos) *Cell {
	to := g.Cell(pTo)
	from := g.Cell(pFrom)
	if from == nil || to == nil {
		g.Debug(-1, "!CanMove (bad cell) %s->%s", pFrom, pTo)
		return nil
	}
	if to.height > from.height && to.height-from.height > 1 {
		g.Debug(-1, "!CanMove (too high) %s->%s", pFrom, pTo)
		return nil
	}
	g.Debug(-1, "CanMove %s->%s", pFrom, pTo)
	return to
}

// Generates a list of squares we can move to from given Pos, we can move to a square
// iff (AND)
// - it's on the grid
// - it's at or above our current height (avoid going back down - or at least let's start trying that)
func (g *Grid) MoveOptions(from Pos) []*Cell {
	rv := []*Cell{}
	if c := g.CanMove(from, from.Left()); c != nil {
		rv = append(rv, c)
	}
	if c := g.CanMove(from, from.Right()); c != nil {
		rv = append(rv, c)
	}
	if c := g.CanMove(from, from.Up()); c != nil {
		rv = append(rv, c)
	}
	if c := g.CanMove(from, from.Down()); c != nil {
		rv = append(rv, c)
	}
	return rv
}

func (g *Grid) Debug(dist int, msg string, args ...any) {
	if dist < 1 && os.Getenv("DEBUG") == "" {
		return
	}
	if dist > 0 {
		fmt.Printf(strings.Repeat(" ", (dist-1)%120))
	}
	fmt.Printf(msg, args...)
	fmt.Println("")
}

func (g *Grid) Distance(fromP Pos, to Pos, dist int) int {
	from := g.Cell(fromP)
	if from.dFrom != 0 {
		return from.dFrom
	}
	if from.dTo != 0 && from.dTo <= dist {
		// Already been here... (looping)
		return -1
	}
	g.Debug(dist, "%d: %s", dist, from)
	from.dTo = dist

	best := -1
	for _, c := range g.MoveOptions(fromP) {
		if c.id.DistFrom(to) == 0 {
			g.Debug(dist, "%d: ** Found the destination (%s) from %s", dist, c, from)
			// actually this is the end, so we don't need to explore further!
			best = 0
			break
		}
		d := g.Distance(c.id, to, dist+1)
		if d != -1 && (best == -1 || d < best) {
			best = d
		}
	}
	if best == -1 {
		return -1 // no path, dead end
	}
	from.dFrom = best + 1
	return from.dFrom
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	grid := Grid{}
	for s.Scan() {
		grid.ParseRow(s.Text())
	}
	grid.Print()
	fmt.Println()

	end := grid.Cell(grid.end)

	rv := []int{}
	for i := 0; i < len(grid.cells); i++ {
		rv = append(rv, grid.Distance(Pos{row: i, col: 0}, end.id, 1))
	}

	min := rv[0]
	for _, n := range rv {
		fmt.Println(n)
		if n < min {
			min = n
		}
	}
	fmt.Println("")
	fmt.Println(min)
}

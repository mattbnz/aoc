package day10

import (
	"bufio"
	"fmt"
	"io"
)

type Pos struct {
	row, col int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d,%d", p.row, p.col)
}

type CardinalDirection int

const (
	NORTH CardinalDirection = iota
	EAST
	SOUTH
	WEST
)

func (b CardinalDirection) String() string {
	if b == NORTH {
		return "^"
	} else if b == EAST {
		return ">"
	} else if b == SOUTH {
		return "v"
	} else if b == WEST {
		return "<"
	}
	return "!"
}

func NewCardinalDirection(s string) (CardinalDirection, error) {
	if s == "<" {
		return NORTH, nil
	} else if s == ">" {
		return EAST, nil
	} else if s == "^" {
		return SOUTH, nil
	} else if s == "v" {
		return WEST, nil
	}
	return NORTH, fmt.Errorf("%s is not a cardinal direction", s)
}

type Cell interface {
	// Prints a representation of the cell.
	String() string

	// Returns a new instance of Cell based on the given string.
	New(string) Cell
}

type BaseCell struct {
	Symbol string
}

func (c BaseCell) String() string {
	return c.Symbol
}

type Grid struct {
	c              map[Pos]Cell
	maxrow, maxcol int
}

func (g Grid) String() string {
	return fmt.Sprintf("Grid of %dx%d", g.maxrow, g.maxcol)
}

func (g Grid) C(p Pos) Cell {
	return g.c[p]
}

// Returns a copy of this grid at the current minute (doesn't preserve past minutes)
func (g Grid) Copy() Grid {
	ng := Grid{
		maxrow: g.maxrow,
		maxcol: g.maxcol,
	}
	for p, c := range g.c {
		ng.c[p] = c
	}
	return ng
}

func (g Grid) Print() {
	for row := 1; row <= g.maxrow; row++ {
		for col := 1; col <= g.maxcol; col++ {
			fmt.Print(g.C(Pos{row, col}))
		}
		fmt.Println()
	}
	fmt.Println()
}

func NewGrid[C Cell](r io.Reader) Grid {
	g := Grid{c: map[Pos]Cell{}}
	g.maxrow = -1
	g.maxcol = -1

	s := bufio.NewScanner(r)

	var cFactory C

	row := 1
	for s.Scan() {
		for col, cStr := range s.Text() {
			p := Pos{row, col + 1}
			c := cFactory.New(string(cStr))
			g.c[p] = c
			g.maxcol = Max(g.maxcol, col+1)
		}
		g.maxrow = Max(g.maxrow, row)
		row++
	}
	return g
}

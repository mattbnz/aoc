package day10

import (
	"bufio"
	"fmt"
	"io"
)

// 1 Based row, col indices
type Pos struct {
	row, col int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d,%d", p.row, p.col)
}

func (p Pos) IsZero() bool {
	return p.row == 0 && p.col == 0
}

type CardinalDirection int

const (
	NO_DIRECTION CardinalDirection = iota
	NORTH
	EAST
	SOUTH
	WEST
)

var CardinalDirections = []CardinalDirection{NORTH, EAST, SOUTH, WEST}

func (b CardinalDirection) String() string {
	if b == NORTH {
		return "N"
	} else if b == EAST {
		return "E"
	} else if b == SOUTH {
		return "S"
	} else if b == WEST {
		return "W"
	}
	return "!"
}

func (b CardinalDirection) Opposite() CardinalDirection {
	switch b {
	case NORTH:
		return SOUTH
	case SOUTH:
		return NORTH
	case EAST:
		return WEST
	case WEST:
		return EAST
	}
	return NO_DIRECTION
}

func NewCardinalDirection(s string) (CardinalDirection, error) {
	if s == "^" || s == "N" {
		return NORTH, nil
	} else if s == "v" || s == "S" {
		return SOUTH, nil
	} else if s == ">" || s == "E" {
		return EAST, nil
	} else if s == "<" || s == "W" {
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

func (g Grid) Next(p Pos, dir CardinalDirection) (np Pos, c Cell, found bool) {
	switch dir {
	case NORTH:
		np.row, np.col = p.row-1, p.col
	case SOUTH:
		np.row, np.col = p.row+1, p.col
	case EAST:
		np.row, np.col = p.row, p.col+1
	case WEST:
		np.row, np.col = p.row, p.col-1
	default:
		found = false
		return
	}
	if np.row < 1 || np.col < 1 || np.row > g.maxrow || np.col > g.maxcol {
		found = false
		return
	}
	found = true
	c = g.C(np)
	return
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

func (g Grid) Each(cb func(Pos, Cell) bool) {
	for row := 1; row <= g.maxrow; row++ {
		for col := 1; col <= g.maxcol; col++ {
			p := Pos{row, col}
			if !cb(p, g.C(p)) {
				return
			}
		}
	}
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

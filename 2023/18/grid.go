package day18

import (
	"bufio"
	"fmt"
	"io"
	"math"
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
	New(string, Pos) Cell
}

type BaseCell struct {
	id     Pos
	Symbol string
}

func (c BaseCell) String() string {
	return c.Symbol
}

func (c BaseCell) New(s string, p Pos) Cell {
	return BaseCell{Symbol: s, id: p}
}

// FlexGrid can have indices below 1 vs Grid.
type FlexGrid struct {
	c              map[Pos]Cell
	minrow, mincol int
	maxrow, maxcol int
}

func (g FlexGrid) String() string {
	return fmt.Sprintf("Grid of %dx%d", g.maxrow-g.minrow, g.maxcol-g.mincol)
}

func (g FlexGrid) C(p Pos) Cell {
	c := g.c[p]
	return c
}

func (g *FlexGrid) SetC(p Pos, v Cell) {
	g.c[p] = v
	if p.row < g.minrow {
		g.minrow = p.row
	}
	if p.row > g.maxrow {
		g.maxrow = p.row
	}
	if p.col < g.mincol {
		g.mincol = p.col
	}
	if p.col > g.maxcol {
		g.maxcol = p.col
	}
}

func (g FlexGrid) Next(p Pos, dir CardinalDirection) (np Pos, c Cell, found bool) {
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
	if np.row < g.minrow || np.col < g.mincol || np.row > g.maxrow || np.col > g.maxcol {
		found = false
		return
	}
	found = true
	c = g.C(np)
	return
}

// Returns a copy of this grid at the current minute (doesn't preserve past minutes)
func (g FlexGrid) Copy() FlexGrid {
	ng := FlexGrid{
		minrow: g.minrow,
		maxrow: g.maxrow,
		mincol: g.mincol,
		maxcol: g.maxcol,
		c:      map[Pos]Cell{},
	}
	for p, c := range g.c {
		ng.c[p] = c
	}
	return ng
}

func (g FlexGrid) Print() {
	for row := g.minrow; row <= g.maxrow; row++ {
		for col := g.mincol; col <= g.maxcol; col++ {
			c := g.C(Pos{row, col})
			if c == nil {
				fmt.Print(".")
			} else {
				fmt.Print(c)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (g FlexGrid) PrintNumbered() {
	fmt.Print("X")
	for col := g.mincol; col <= g.maxcol; col++ {
		fmt.Printf("%d", Abs(col)%10)
	}
	fmt.Println()
	for row := g.minrow; row <= g.maxrow; row++ {
		fmt.Printf("%d", Abs(row)%10)
		for col := g.mincol; col <= g.maxcol; col++ {
			c := g.C(Pos{row, col})
			if c == nil {
				fmt.Print(".")
			} else {
				fmt.Print(c)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (g FlexGrid) Each(cb func(Pos, Cell) bool) {
	for row := g.minrow; row <= g.maxrow; row++ {
		for col := g.mincol; col <= g.maxcol; col++ {
			p := Pos{row, col}
			if !cb(p, g.C(p)) {
				return
			}
		}
	}
}

func (g FlexGrid) Equal(o FlexGrid) bool {
	if g.maxrow != o.maxrow || g.maxcol != o.maxcol {
		return false
	}

	for row := g.minrow; row <= g.maxrow; row++ {
		for col := g.mincol; col <= g.maxcol; col++ {
			p := Pos{row, col}
			if g.C(p) != o.C(p) {
				return false
			}
		}
	}
	return true
}

func EmptyGrid() (g FlexGrid) {
	g.c = map[Pos]Cell{}
	g.minrow = int(math.Pow(2, 64)) - 1
	g.mincol = int(math.Pow(2, 64)) - 1
	g.maxrow = -1
	g.maxcol = -1
	return
}

func NewGrid[C Cell](r io.Reader) FlexGrid {
	return NewGridFromScanner[C](bufio.NewScanner(r))
}

func NewGridFromScanner[C Cell](s *bufio.Scanner) FlexGrid {
	g := EmptyGrid()

	var cFactory C

	row := 1
	for s.Scan() {
		if s.Text() == "" {
			return g
		}
		for col, cStr := range s.Text() {
			p := Pos{row, col + 1}
			c := cFactory.New(string(cStr), p)
			g.SetC(p, c)
			g.maxcol = Max(g.maxcol, col+1)
		}
		g.maxrow = Max(g.maxrow, row)
		row++
	}
	return g
}

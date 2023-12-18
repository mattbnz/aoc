package day18

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

const TRENCH = "#"
const GROUND = "*"
const DEFAULT = "."

var digDir = map[string]CardinalDirection{
	"U": NORTH,
	"R": EAST,
	"D": SOUTH,
	"L": WEST,
}

type LagoonCell struct {
	BaseCell

	visited bool
}

type Lagoon struct {
	FlexGrid

	Digger Pos
}

// Moves l.Digger in the specified direction.
func (l *Lagoon) Move(dir CardinalDirection) {
	switch dir {
	case NORTH:
		l.Digger = Pos{l.Digger.row - 1, l.Digger.col}
	case EAST:
		l.Digger = Pos{l.Digger.row, l.Digger.col + 1}
	case SOUTH:
		l.Digger = Pos{l.Digger.row + 1, l.Digger.col}
	case WEST:
		l.Digger = Pos{l.Digger.row, l.Digger.col - 1}
	}
}

// Digs a trench of length in dir from l.Digger
func (l *Lagoon) Trench(dir CardinalDirection, length int) {
	for n := 1; n <= length; n++ {
		l.Move(dir)
		l.Dig()
	}
}

// Digs a hole at l.Digger
func (l *Lagoon) Dig() {
	l.SetC(l.Digger, &LagoonCell{BaseCell: BaseCell{id: l.Digger, Symbol: TRENCH}})
}

func NewLagoon(filename string) (l Lagoon, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	origin := Pos{1, 1}
	l.FlexGrid = EmptyGrid()
	l.Digger = origin
	l.Dig()
	for s.Scan() {
		if s.Text() == "" {
			return
		}
		parts := strings.SplitN(s.Text(), " ", -1)
		if len(parts) != 3 {
			err = fmt.Errorf("bad input line: %s", s.Text())
			return
		}
		dir, ok := digDir[parts[0]]
		if !ok {
			err = fmt.Errorf("bad dig direction: %s", parts[0])
			return
		}
		length, lErr := strconv.Atoi(parts[1])
		if lErr != nil {
			err = lErr
			return
		}
		// Colour ignored for part1.
		l.Trench(dir, length)
	}
	if l.Digger != origin {
		glog.Warningf("Digger stopped at %s not %s after all input lines!!", l.Digger, origin)
	}
	return
}

func (c *Lagoon) TrenchLength() (rv int) {
	c.Each(func(_ Pos, c Cell) bool {
		if c == nil {
			return true
		}
		if c.(*LagoonCell).Symbol == TRENCH {
			rv++
		}
		return true
	})
	return
}

func (c *Lagoon) C(p Pos) *LagoonCell {
	rv := c.FlexGrid.C(p)
	if rv == nil {
		lc := &LagoonCell{BaseCell: BaseCell{id: p, Symbol: "."}}
		c.FlexGrid.c[p] = lc
		return lc
	}
	return rv.(*LagoonCell)
}

func (c *Lagoon) Visit(p Pos) {
	cell := c.C(p)
	if cell.Symbol == TRENCH {
		return // nothing to do for trench cells
	}
	if cell.visited {
		return
	}
	cell.visited = true
	cell.Symbol = GROUND
	for _, d := range CardinalDirections {
		np, _, found := c.FlexGrid.Next(p, d)
		if !found {
			continue
		}
		nc := c.C(np)
		if nc.Symbol == TRENCH {
			continue
		}
		c.Visit(np)
	}
}

func (c *Lagoon) VisitAll() {
	// Visit (recursively from every outside cell inwards)
	for col := c.mincol; col <= c.maxcol; col++ {
		c.Visit(Pos{c.minrow, col})
		c.Visit(Pos{c.maxrow, col})
	}
	for row := c.minrow; row <= c.maxrow; row++ {
		c.Visit(Pos{row, c.mincol})
		c.Visit(Pos{row, c.maxcol})
	}
}

func (c *Lagoon) RowVolume(row int) (rv int) {
	c.VisitAll()
	s := ""
	for col := c.mincol; col <= c.maxcol; col++ {
		cell := c.C(Pos{row, col})
		if cell.Symbol == TRENCH || cell.Symbol == DEFAULT {
			rv++
		}
		s += cell.Symbol
	}
	glog.V(1).Infof("Row % 4d: %s %d", row, s, rv)
	return
}

func (c *Lagoon) Volume() (rv int) {
	c.VisitAll()
	c.Each(func(p Pos, _ Cell) bool {
		cell := c.C(p)
		if cell.Symbol == TRENCH || cell.Symbol == DEFAULT {
			rv++
		}
		return true
	})
	return
}

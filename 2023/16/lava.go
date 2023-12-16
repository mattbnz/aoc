package day16

import (
	"os"

	"github.com/golang/glog"
)

type MirrorCell struct {
	BaseCell

	Cache  map[CardinalDirection]int
	onPath map[CardinalDirection]bool
}

func (c *MirrorCell) New(s string, p Pos) Cell {
	return &MirrorCell{
		BaseCell: BaseCell{id: p, Symbol: s},
		Cache:    map[CardinalDirection]int{},
		onPath:   map[CardinalDirection]bool{},
	}
}

func (c *MirrorCell) Beam(from CardinalDirection) (next []CardinalDirection, cacheKey CardinalDirection) {
	cacheKey = from
	var to CardinalDirection
	if c.Symbol == "." {
		to = from
	} else if c.Symbol == "\\" {
		switch from {
		case NORTH:
			to = WEST
		case SOUTH:
			to = EAST
		case EAST:
			to = SOUTH
		case WEST:
			to = NORTH
		}
	} else if c.Symbol == "/" {
		switch from {
		case NORTH:
			to = EAST
		case SOUTH:
			to = WEST
		case EAST:
			to = NORTH
		case WEST:
			to = SOUTH
		}
	} else if c.Symbol == "|" {
		switch from {
		case EAST:
			to = NORTH
			next = append(next, SOUTH)
			cacheKey = EAST
		case WEST:
			to = NORTH
			next = append(next, SOUTH)
			cacheKey = EAST
		default:
			to = from
		}
	} else if c.Symbol == "-" {
		switch from {
		case NORTH:
			to = WEST
			next = append(next, EAST)
			cacheKey = NORTH
		case SOUTH:
			to = WEST
			next = append(next, EAST)
			cacheKey = NORTH
		default:
			to = from
		}
	}
	next = append(next, to)
	return
}

type MirrorGrid struct {
	Grid

	beamCount int
}

func NewMirrorGrid(filename string) (m MirrorGrid, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	m.Grid = NewGrid[*MirrorCell](f)
	return
}

func (g *MirrorGrid) Beam(id int, depth int, in Pos, heading CardinalDirection) int {
	c := g.Grid.C(in).(*MirrorCell)
	next, cacheKey := c.Beam(heading)

	if length, found := c.Cache[cacheKey]; found {
		if depth == 0 {
			return length
		}
		return 0
	}

	inc := 1
	sym := "+"
	if len(c.Cache) > 0 || len(c.onPath) > 0 {
		inc = 0 // node already visited by another direction, don't count again.
		sym = " "
	}
	glog.V(2).Infof("Beam % 2d (% 2d) %s => %s will go %s", id, depth, heading, in, next)

	length := 0
	wentAnywhere := false
	for _, nHeading := range next {
		if c.onPath[nHeading] {
			glog.V(2).Infof("Beam % 2d can't go %s from %s (already on that path!)", id, nHeading, in)
			continue
		}
		wentAnywhere = true
		newNext, _, found := g.Next(c.id, nHeading)
		if !found {
			glog.V(2).Infof("Beam % 2d can't go %s from %s", id, nHeading, in)
			continue
		}
		c.onPath[nHeading] = true
		length += g.Beam(id, depth+1, newNext, nHeading)
		delete(c.onPath, nHeading)
	}
	if wentAnywhere {
		length += inc
	}
	glog.V(2).Infof("Beam % 2d (% 2d) %s => %s%s => %s returns %d", id, depth, heading, in, sym, next, length)
	c.Cache[cacheKey] = length
	return length
}

// Returns a copy of this grid at the current minute (doesn't preserve past minutes)
func (g *MirrorGrid) Copy() Grid {
	ng := Grid{
		maxrow: g.maxrow,
		maxcol: g.maxcol,
		c:      map[Pos]Cell{},
	}
	for p, c := range g.c {
		ng.c[p] = (&MirrorCell{}).New(c.(*MirrorCell).Symbol, p)
	}
	return ng
}

func (g *MirrorGrid) Best() (best int) {
	// North and South
	for c := 1; c <= g.maxcol; c++ {
		sG := MirrorGrid{Grid: g.Copy()}
		n := sG.Beam(1, 1, Pos{1, c}, SOUTH)
		if n > best {
			glog.V(1).Infof("South from %d gives %d as new best", c, n)
			best = n
		}
		nG := MirrorGrid{Grid: g.Copy()}
		n = nG.Beam(1, 1, Pos{g.maxrow, c}, NORTH)
		if n > best {
			glog.V(1).Infof("North from %d gives %d as new best", c, n)
			best = n
		}
	}
	// East and West
	for r := 1; r <= g.maxrow; r++ {
		eG := MirrorGrid{Grid: g.Copy()}
		n := eG.Beam(1, 1, Pos{r, 1}, EAST)
		if n > best {
			glog.V(1).Infof("East from %d gives %d as new best", r, n)
			best = n
		}
		wG := MirrorGrid{Grid: g.Copy()}
		n = wG.Beam(1, 1, Pos{r, g.maxcol}, WEST)
		if n > best {
			glog.V(1).Infof("West from %d gives %d as new best", r, n)
			best = n
		}
	}
	return
}

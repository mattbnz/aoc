package day10

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
)

type Corner int

const (
	NO_CORNER Corner = iota
	TL
	TR
	BR
	BL
)

var Corners = []Corner{TL, TR, BR, BL}

func (b Corner) String() string {
	if b == TL {
		return "TL"
	} else if b == TR {
		return "TR"
	} else if b == BR {
		return "BR"
	} else if b == BL {
		return "BL"
	}
	return "!"
}

// Returns the opposite corner on the tile
func (b Corner) Opposite(heading CardinalDirection) Corner {
	return map[Corner]map[CardinalDirection]Corner{
		TL: {EAST: TR, SOUTH: BL},
		TR: {WEST: TL, SOUTH: BR},
		BR: {WEST: BL, NORTH: TR},
		BL: {EAST: BR, NORTH: TL},
	}[b][heading]
}

// Returns the neighbouring corner on an adjacent tile
func (b Corner) Neighbour(heading CardinalDirection) Corner {
	return map[Corner]map[CardinalDirection]Corner{
		TL: {WEST: TR, NORTH: BL},
		TR: {EAST: TL, NORTH: BR},
		BR: {EAST: BL, SOUTH: TR},
		BL: {WEST: BR, SOUTH: TL},
	}[b][heading]
}

type CornerState struct {
	visited bool
	wet     bool
}

type PipeCell struct {
	BaseCell
	id Pos

	IsStart bool // marker for start section.
	// Distance from start in pipe (if this is a connected section)
	Distance int

	// Whether there is a pipe connection in this direction.
	Connects map[CardinalDirection]bool
	// Whether this corner of the tile is wet (aka 'reachable' from the outside).
	WetCorner map[Corner]CornerState
}

var _ Cell = &PipeCell{}

func (c *PipeCell) New(s string, p Pos) Cell {
	nc := &PipeCell{
		id:        p,
		Connects:  map[CardinalDirection]bool{},
		WetCorner: map[Corner]CornerState{},
	}
	nc.Symbol = s
	switch nc.Symbol {
	case "|":
		nc.Connects[NORTH] = true
		nc.Connects[SOUTH] = true
	case "-":
		nc.Connects[WEST] = true
		nc.Connects[EAST] = true
	case "L":
		nc.Connects[NORTH] = true
		nc.Connects[EAST] = true
	case "J":
		nc.Connects[NORTH] = true
		nc.Connects[WEST] = true
	case "7":
		nc.Connects[SOUTH] = true
		nc.Connects[WEST] = true
	case "F":
		nc.Connects[SOUTH] = true
		nc.Connects[EAST] = true
	case "S":
		nc.IsStart = true
	}
	return nc
}

func (c *PipeCell) Travel(from CardinalDirection) (to CardinalDirection, err error) {
	if c.Symbol == "|" && (from == NORTH || from == SOUTH) {
		to = from
	} else if c.Symbol == "-" && (from == EAST || from == WEST) {
		to = from
	} else {
		if !c.Connects[from.Opposite()] {
			err = fmt.Errorf("%s does not connect to the %s", c, from)
			return
		}
		for d, connects := range c.Connects {
			if d != from.Opposite() && connects {
				to = d
				return
			}
		}
		err = fmt.Errorf("%s does not connect to the %s", c, from)
	}
	return
}

type Maze struct {
	g     Grid
	depth int
}

func NewMaze(filename string) (Maze, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return Maze{}, err
	}
	defer f.Close()

	maze := Maze{g: NewGrid[*PipeCell](f)}
	return maze, nil
}

func (m Maze) FindStart() (p Pos, c *PipeCell, err error) {
	err = fmt.Errorf("start marker not found")
	m.g.Each(func(iPos Pos, iCell Cell) bool {
		pCell := iCell.(*PipeCell)
		if pCell.IsStart {
			p = iPos
			c = pCell
			err = nil
			return false
		}
		return true
	})
	return
}

func (m Maze) C(p Pos) *PipeCell {
	c := m.g.C(p)
	return c.(*PipeCell)
}

func (m Maze) LongestPath() (steps int, err error) {
	startPos, startCell, err := m.FindStart()
	if err != nil {
		return
	}
	glog.Infof("Found starting cell %s at %s", startCell, startPos)
	var moveA, moveB CardinalDirection
	var posA, posB Pos
	for _, dir := range CardinalDirections {
		otherPos, otherCell, found := m.g.Next(startPos, dir)
		if !found {
			continue
		}
		oCell := otherCell.(*PipeCell)
		if oCell.Connects[dir.Opposite()] {
			startCell.Connects[dir] = true
			glog.Infof("Found connection %s to %s from starting cell %s", dir, otherPos, startCell)
			if posA.IsZero() {
				moveA = dir
				posA = otherPos
			} else if posB.IsZero() {
				moveB = dir
				posB = otherPos
			} else {
				err = fmt.Errorf("Start cell has more than two potential connections")
				return
			}
		}
	}
	glog.Infof("Starting Moves: %s to %s, %s to %s", moveA, posA, moveB, posB)
	if moveA == NO_DIRECTION || moveB == NO_DIRECTION || posA.IsZero() || posB.IsZero() {
		err = fmt.Errorf("failed to find starting moves")
		return
	}
	steps = 1

	for {
		if posA == posB {
			// Paths have converged
			return
		}
		cellA := m.C(posA)
		cellB := m.C(posB)
		cellA.Distance = steps
		cellB.Distance = steps

		moveA, err = cellA.Travel(moveA)
		if err != nil {
			return
		}
		moveB, err = cellB.Travel(moveB)
		if err != nil {
			return
		}
		var found bool
		posA, _, found = m.g.Next(posA, moveA)
		if !found {
			err = fmt.Errorf("invalid cell moving %s from %s", moveA, posA)
		}
		posB, _, found = m.g.Next(posB, moveB)
		if !found {
			err = fmt.Errorf("invalid cell moving %s from %s", moveB, posB)
		}
		steps++
	}
}

func (m Maze) CountEnclosed() (enclosed int, err error) {
	// Trace the path to work out the bounds
	_, err = m.LongestPath()
	if err != nil {
		return
	}

	/*m.g.Each(func(iPos Pos, iCell Cell) bool {
		c := iCell.(*PipeCell)
		c.Flood(&m, TL)
		return true
	})*/
	c := m.C(Pos{1, 1})
	c.Flood(&m, TL)

	m.g.Each(func(iPos Pos, iCell Cell) bool {
		c := iCell.(*PipeCell)
		if c.Distance == 0 && !c.IsStart && !c.HasWetCorner() {
			enclosed++
		}
		return true
	})
	return
}

func (c *PipeCell) HasWetCorner() bool {
	for corner, state := range c.WetCorner {
		if state.wet {
			glog.V(2).Infof("%s has wet corner %s", c.id, corner)
			return true
		}
	}
	return false
}

// Marks corners in specified directions as wet
func (c *PipeCell) FloodNeighbours(m *Maze, from Corner, dirA, dirB CardinalDirection) {
	_, aCell, aOK := m.g.Next(c.id, dirA)
	if aOK {
		aCell.(*PipeCell).Flood(m, from.Neighbour(dirA))
	}
	bPos, bCell, bOK := m.g.Next(c.id, dirB)
	if bOK {
		bCorner := from.Neighbour(dirB)
		bCell.(*PipeCell).Flood(m, bCorner)
		_, cCell, cOK := m.g.Next(bPos, dirA)
		if cOK {
			cCell.(*PipeCell).Flood(m, bCorner.Neighbour(dirA))
		}
	}
}

func PickCorner(a, b Corner) Corner {
	if a != NO_CORNER {
		return a
	}
	return b
}

func (m *Maze) Depth(inc int) {
	m.depth += inc
}

func (m *Maze) Prefix() string {
	return fmt.Sprintf("% 3d %s", m.depth, strings.Repeat(" ", m.depth/10))
}

// Marks one or more corner assuming water touching the given corner
func (c *PipeCell) Flood(m *Maze, from Corner) {
	if c.WetCorner[from].visited {
		return
	}
	c.WetCorner[from] = CornerState{visited: true, wet: true}
	glog.V(2).Infof("%s%s.%s is now wet", m.Prefix(), c.id, from)
	m.Depth(1)
	defer func() { m.Depth(-1) }()
	// Neighbours
	switch from {
	case TL:
		c.FloodNeighbours(m, from, NORTH, WEST)
	case TR:
		c.FloodNeighbours(m, from, NORTH, EAST)
	case BL:
		c.FloodNeighbours(m, from, SOUTH, WEST)
	case BR:
		c.FloodNeighbours(m, from, SOUTH, EAST)
	}
	// Opposites
	switch c.Symbol {
	case ".":
		for _, corner := range Corners {
			if corner == from {
				continue
			}
			c.Flood(m, corner)
		}
	case "|":
		next := PickCorner(from.Opposite(NORTH), from.Opposite(SOUTH))
		c.Flood(m, next)
	case "-":
		next := PickCorner(from.Opposite(EAST), from.Opposite(WEST))
		c.Flood(m, next)
	case "F":
		if from == BR {
			return
		} else if from == TR {
			c.Flood(m, TL)
			c.Flood(m, BL)
		} else {
			c.Flood(m, TL)
			c.Flood(m, BL)
			c.Flood(m, TR)
		}
	case "7":
		if from == BL {
			return
		} else if from == TL {
			c.Flood(m, TR)
			c.Flood(m, BR)
		} else {
			c.Flood(m, TR)
			c.Flood(m, BR)
			c.Flood(m, TL)
		}
	case "L":
		if from == TR {
			return
		} else if from == BR {
			c.Flood(m, TL)
			c.Flood(m, BL)
		} else {
			c.Flood(m, TL)
			c.Flood(m, BL)
			c.Flood(m, BR)
		}
	case "J":
		if from == TL {
			return
		} else if from == BL {
			c.Flood(m, TR)
			c.Flood(m, BR)
		} else {
			c.Flood(m, TR)
			c.Flood(m, BR)
			c.Flood(m, BL)
		}
	}
}

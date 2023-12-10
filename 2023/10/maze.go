package day10

import (
	"fmt"
	"os"

	"github.com/golang/glog"
)

type PipeCell struct {
	BaseCell

	Connects map[CardinalDirection]bool
	IsStart  bool

	Distance int
}

var _ Cell = &PipeCell{}

func (c *PipeCell) New(s string) Cell {
	nc := &PipeCell{Connects: map[CardinalDirection]bool{}}
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
	g Grid
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

package day10

import "os"

type PipeCell struct {
	BaseCell

	Connects map[CardinalDirection]bool
	IsStart  bool
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

func (m Maze) LongestPath() int {
	return 0
}

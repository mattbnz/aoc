package day18

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

const HOLE = "#"
const GROUND = "."

var digDir = map[string]CardinalDirection{
	"U": NORTH,
	"R": EAST,
	"D": SOUTH,
	"L": WEST,
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
	l.SetC(l.Digger, BaseCell{id: l.Digger, Symbol: HOLE})
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
		if c.(BaseCell).Symbol == HOLE {
			rv++
		}
		return true
	})
	return
}

func (c *Lagoon) RowVolume(row int) (rv int) {
	rowStr := ""
	for col := c.mincol; col <= c.maxcol; col++ {
		p := Pos{row, col}
		cell := c.C(p).(BaseCell)
		rowStr += cell.Symbol
	}
	idx := 0
	trenchStart := -1
	maybeGroundStart := -1
	for idx <= len(rowStr) {
		// trenchStart not set == not in trench OR hole, look for a trench.
		if trenchStart == -1 {
			trenchStart = strings.Index(rowStr[idx:], HOLE)
			if trenchStart == -1 {
				break // no more trenches
			}
			trenchStart += idx
			idx = trenchStart + 1
			continue
		}
		// in a trench coming out of a hole
		if maybeGroundStart != -1 {
			ground := strings.Index(rowStr[idx:], GROUND)
			if ground == -1 {
				ground = len(rowStr)
			} else {
				ground += idx
			}
			glog.V(2).Infof("row % 4d: Found trench (with hole) from %d - %d: %d", row, trenchStart, ground, ground-trenchStart)
			rv += ground - trenchStart
			trenchStart = -1
			maybeGroundStart = -1
			idx = ground
			continue
		}
		// in a trench, maybe going into a hole.
		maybeHoleStart := strings.Index(rowStr[idx:], GROUND)
		if maybeHoleStart == -1 {
			// trench extends to end of row
			rv += len(rowStr) - trenchStart
			glog.V(2).Infof("row % 4d: Found trench (no hole) from %d - END: %d", row, trenchStart, len(rowStr)-trenchStart)
			trenchStart = -1
			break
		}
		maybeHoleStart += idx + 1
		// if we can find another trench ahead, then we're in a hole, otherwise the trench ended on the previous tile.
		holeEnd := strings.Index(rowStr[maybeHoleStart:], HOLE)
		if holeEnd == -1 {
			// trench ended on the previous tile
			rv += maybeHoleStart - 1 - trenchStart
			glog.V(2).Infof("row % 4d: Found trench (R) from %d - %d: %d", row, trenchStart, maybeHoleStart-1, maybeHoleStart-1-trenchStart)
			trenchStart = -1
			break
		}
		holeEnd += maybeHoleStart
		maybeGroundStart = holeEnd + 1
		idx = maybeGroundStart
	}
	if trenchStart != -1 {
		// Ended in a trench
		rv += len(rowStr) - trenchStart
		glog.V(2).Infof("row % 4d: Ended in trench from %d: %d", row, trenchStart, len(rowStr)-trenchStart)
	}
	glog.V(1).Infof("row % 4d: %s %d", row, rowStr, rv)
	return
}

func (c *Lagoon) Volume() (rv int) {
	for row := c.minrow; row <= c.maxrow; row++ {
		rv += c.RowVolume(row)
	}
	return
}

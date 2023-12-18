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
	inTrench := false
	seenDot := true
	rowDbg := ""
	for col := c.mincol; col <= c.maxcol; col++ {
		p := Pos{row, col}
		cell := c.C(p).(BaseCell)
		if cell.Symbol != HOLE {
			if inTrench {
				rv++
			}
			seenDot = true
		} else {
			if inTrench {
				if seenDot {
					inTrench = false
				}
			} else {
				if seenDot {
					inTrench = true
				}
			}
			seenDot = false
			rv++
		}
		rowDbg += cell.Symbol
	}
	glog.V(1).Infof("row % 4d: %s %d", row, rowDbg, rv)
	return
}

func (c *Lagoon) Volume() (rv int) {
	for row := c.minrow; row <= c.maxrow; row++ {
		rv += c.RowVolume(row)
	}
	return
}

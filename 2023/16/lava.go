package day16

import (
	"os"
	"sync"

	"github.com/golang/glog"
)

type MirrorCell struct {
	BaseCell

	EnergizedFrom map[CardinalDirection]bool
	m             sync.Mutex // protects EnergizedFrom
}

func (c *MirrorCell) New(s string, p Pos) Cell {
	return &MirrorCell{
		BaseCell:      BaseCell{id: p, Symbol: s},
		EnergizedFrom: map[CardinalDirection]bool{},
	}
}

type MirrorGrid struct {
	Grid

	beamCount int
	wg        sync.WaitGroup
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

func (g *MirrorGrid) StartBeam(in Pos, heading CardinalDirection) {
	var id = g.beamCount
	g.beamCount++
	g.wg.Add(1)
	glog.V(1).Infof("Starting beam % 2d from %s going %s", id, in, heading)
	go func() {
		for {
			glog.V(1).Infof("Beam % 2d at %s heading %s", id, in, heading)
			next, going, cont := g.doCell(in, heading)
			if !cont {
				glog.V(1).Infof("Beam % 2d: ended at %s going %s ", id, next, going)
				break
			}
			in = next
			heading = going
		}
		g.wg.Done()
	}()
}

func (g *MirrorGrid) Wait() {
	g.wg.Wait()
}

func (g *MirrorGrid) doCell(in Pos, heading CardinalDirection) (next Pos, going CardinalDirection, cont bool) {
	c := g.Grid.C(in).(*MirrorCell)
	c.m.Lock()
	defer c.m.Unlock()

	if c.EnergizedFrom[heading] {
		cont = false
		return
	}
	c.EnergizedFrom[heading] = true

	var newBeam CardinalDirection

	if c.Symbol == "." {
		going = heading
	} else if c.Symbol == "\\" {
		switch heading {
		case NORTH:
			going = WEST
		case SOUTH:
			going = EAST
		case EAST:
			going = SOUTH
		case WEST:
			going = NORTH
		}
	} else if c.Symbol == "/" {
		switch heading {
		case NORTH:
			going = EAST
		case SOUTH:
			going = WEST
		case EAST:
			going = NORTH
		case WEST:
			going = SOUTH
		}
	} else if c.Symbol == "|" {
		switch heading {
		case EAST:
			going = NORTH
			newBeam = SOUTH
		case WEST:
			going = NORTH
			newBeam = SOUTH
		default:
			going = heading
		}
	} else if c.Symbol == "-" {
		switch heading {
		case NORTH:
			going = WEST
			newBeam = EAST
		case SOUTH:
			going = WEST
			newBeam = EAST
		default:
			going = heading
		}
	}
	if newBeam != NO_DIRECTION {
		newNext, _, found := g.Next(c.id, newBeam)
		if found {
			g.StartBeam(newNext, newBeam)
		}
	}
	next, _, cont = g.Next(c.id, going)
	return
}

func (g *MirrorGrid) Energized() (sum int) {
	g.Each(func(_ Pos, c Cell) bool {
		if len(c.(*MirrorCell).EnergizedFrom) > 0 {
			sum++
		}
		return true
	})
	return
}

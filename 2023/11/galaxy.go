package day11

import (
	"os"

	"github.com/golang/glog"
)

type Space struct {
	g Grid

	Galaxies map[int]Pos

	RowSizes map[int]int
	ColSizes map[int]int
}

func NewSpace(filename string) (Space, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return Space{}, err
	}
	defer f.Close()

	space := Space{g: NewGrid[BaseCell](f), Galaxies: map[int]Pos{}}
	return space, nil
}

func (s *Space) Expand() {
	row := 1
	for {
		if row > s.g.maxrow {
			break
		}
		hasGalaxy := false
		for col := 1; col <= s.g.maxcol; col++ {
			if s.g.C(Pos{row, col}).(BaseCell).Symbol != "." {
				hasGalaxy = true
				break
			}
		}
		if hasGalaxy {
			row++
			continue
		}
		s.g.DupRow(row)
		row += 2
	}

	col := 1
	for {
		if col > s.g.maxcol {
			break
		}
		hasGalaxy := false
		for row := 1; row <= s.g.maxrow; row++ {
			if s.g.C(Pos{row, col}).(BaseCell).Symbol != "." {
				hasGalaxy = true
				break
			}
		}
		if hasGalaxy {
			col++
			continue
		}
		s.g.DupCol(col)
		col += 2
	}
}

func (s *Space) FindGalaxies() {
	gN := 1
	s.g.Each(func(p Pos, c Cell) bool {
		if c.(BaseCell).Symbol == "#" {
			s.Galaxies[gN] = p
			glog.V(1).Infof("Galaxy %d at %s", gN, p)
			gN++
		}
		return true
	})
	glog.Infof("Found %d galaxies", len(s.Galaxies))
}

func (s *Space) PathLength(a, b Pos) int {
	return Abs(a.row-b.row) + Abs(a.col-b.col)
}

type gKey struct {
	a, b int
}

func (s *Space) PathSum() (rv int) {
	s.Expand()
	s.FindGalaxies()
	cache := map[gKey]int{}
	for gA, aPos := range s.Galaxies {
		for gB, bPos := range s.Galaxies {
			key := gKey{Min(gA, gB), Max(gA, gB)}
			if _, done := cache[key]; done {
				continue
			}
			length := s.PathLength(aPos, bPos)
			rv += length
			cache[key] = length
			glog.V(1).Infof("Shortest path from %d@%s to %d@%s is %d", gA, aPos, gB, bPos, length)
		}
	}
	return
}

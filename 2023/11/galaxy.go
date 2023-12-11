package day11

import (
	"os"

	"github.com/golang/glog"
)

type Space struct {
	g Grid

	Galaxies map[int]Pos

	ExpandedRows []int
	ExpandedCols []int
	ExpandedBy   int
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

func (s *Space) ExpandBy(by int) {
	for row := 1; row <= s.g.maxrow; row++ {
		hasGalaxy := false
		for col := 1; col <= s.g.maxcol; col++ {
			if s.g.C(Pos{row, col}).(BaseCell).Symbol != "." {
				hasGalaxy = true
				break
			}
		}
		if hasGalaxy {
			continue
		}
		s.ExpandedRows = append(s.ExpandedRows, row)
	}

	for col := 1; col <= s.g.maxcol; col++ {
		hasGalaxy := false
		for row := 1; row <= s.g.maxrow; row++ {
			if s.g.C(Pos{row, col}).(BaseCell).Symbol != "." {
				hasGalaxy = true
				break
			}
		}
		if hasGalaxy {
			continue
		}
		s.ExpandedCols = append(s.ExpandedCols, col)
	}
	s.ExpandedBy = by
}

func (s *Space) FindGalaxies() {
	gN := 1
	s.g.Each(func(p Pos, c Cell) bool {
		if c.(BaseCell).Symbol == "#" {
			ep := s.ExpandPos(p)
			s.Galaxies[gN] = ep
			glog.V(1).Infof("Galaxy %d at %s (from %s)", gN, ep, p)
			gN++
		}
		return true
	})
	glog.Infof("Found %d galaxies", len(s.Galaxies))
}

// returns the virtual (expanded pos p)
func (s *Space) ExpandPos(p Pos) (ep Pos) {
	if s.ExpandedBy == 0 {
		return p
	}
	rEidx := 0
	ep.row = 1
	for row := 1; row < p.row; row++ {
		if rEidx < len(s.ExpandedRows) && s.ExpandedRows[rEidx] == row {
			ep.row += s.ExpandedBy - 1
			rEidx++
		}
		ep.row++
	}

	cEidx := 0
	ep.col = 1
	for col := 1; col < p.col; col++ {
		if cEidx < len(s.ExpandedCols) && s.ExpandedCols[cEidx] == col {
			ep.col += s.ExpandedBy - 1
			cEidx++
		}
		ep.col++
	}
	return
}

func (s *Space) PathLength(a, b Pos) int {
	return Abs(a.row-b.row) + Abs(a.col-b.col)
}

type gKey struct {
	a, b int
}

func (s *Space) PathSum(expansionFactor int) (rv int) {
	if expansionFactor == 2 {
		s.Expand()
	} else {
		s.ExpandBy(expansionFactor)
	}
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

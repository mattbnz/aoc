package day13

import (
	"bufio"
	"os"

	"github.com/golang/glog"
)

type MirrorValley struct {
	Mirrors []Grid
}

func (m MirrorValley) ReflectsVertical(g Grid) (int, bool) {
	for col := 2; col <= g.maxcol; col++ {
		if m.colsMatch(g, col, col-1) && m.checkColReflects(g, col) {
			glog.Infof("Found good reflection around col %d", col)
			return col, true
		}
	}
	return -1, false
}

func (m MirrorValley) ReflectsHorizontal(g Grid) (int, bool) {
	for row := 2; row <= g.maxrow; row++ {
		if m.rowsMatch(g, row, row-1) && m.checkRowReflects(g, row) {
			glog.Infof("Found good reflection around row %d", row)
			return row, true
		}
	}
	return -1, false
}

func (m MirrorValley) colsMatch(g Grid, c1, c2 int) bool {
	for r := 1; r <= g.maxrow; r++ {
		p1, p2 := Pos{r, c1}, Pos{r, c2}
		v1, v2 := g.C(p1), g.C(p2)
		if v1.(BaseCell).Symbol != v2.(BaseCell).Symbol {
			glog.V(1).Infof("mismatch between %s (%s) and %s (%s)", p1, v1, p2, v2)
			return false
		}
	}
	glog.V(1).Infof("found matching cols %d and %d", c1, c2)
	return true
}

func (m MirrorValley) rowsMatch(g Grid, r1, r2 int) bool {
	for c := 1; c <= g.maxcol; c++ {
		p1, p2 := Pos{r1, c}, Pos{r2, c}
		v1, v2 := g.C(p1), g.C(p2)
		if v1.(BaseCell).Symbol != v2.(BaseCell).Symbol {
			glog.V(1).Infof("mismatch between %s (%s) and %s (%s)", p1, v1, p2, v2)
			return false
		}
	}
	glog.V(1).Infof("found matching rows %d and %d", r1, r2)
	return true
}

// checks the reflection point before row matches the rest of the field.
func (m MirrorValley) checkRowReflects(g Grid, row int) bool {
	step := 1
	for {
		l := row - 1 - step
		r := row + step
		if l < 1 || r > g.maxrow {
			break
		}
		if !m.rowsMatch(g, l, r) {
			return false
		}
		step++
	}
	return true
}

// checks the reflection point before col matches the rest of the field.
func (m MirrorValley) checkColReflects(g Grid, col int) bool {
	step := 1
	for {
		l := col - 1 - step
		r := col + step
		if l < 1 || r > g.maxcol {
			break
		}
		if !m.colsMatch(g, l, r) {
			return false
		}
		step++
	}
	return true
}

func (m MirrorValley) Score(g Grid) int {
	r, ok := m.ReflectsHorizontal(g)
	if ok {
		return (r - 1) * 100
	}
	c, ok := m.ReflectsVertical(g)
	if ok {
		return c - 1
	}
	return 0
}

func (m MirrorValley) ScoreSum() (rv int) {
	for _, g := range m.Mirrors {
		rv += m.Score(g)
	}
	return
}

func NewMirrorValley(filename string) (valley MirrorValley, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	for {
		g := NewGridFromScanner[BaseCell](s)
		if g.maxrow == -1 {
			break
		}
		valley.Mirrors = append(valley.Mirrors, g)
	}
	return
}

package day13

import (
	"bufio"
	"fmt"
	"os"

	"github.com/golang/glog"
)

type MirrorValley struct {
	Mirrors []Grid
	Smudged bool
}

func (m MirrorValley) ReflectsVertical(g Grid, ignoreCol int) (reflectsOn int, nearMatches []Mismatch, ok bool) {
	reflectsOn = -1
	for col := 2; col <= g.maxcol; col++ {
		if col == ignoreCol {
			continue
		}
		mismatched := m.colMismatch(g, col, col-1)
		if len(mismatched) != 0 {
			continue
		}
		nm, colOk := m.checkColReflects(g, col)
		if colOk {
			if reflectsOn == -1 {
				glog.Infof("Found good reflection around col %d", col)
				reflectsOn = col
			} else {
				glog.Infof("Ignoring additional reflection around col %d", col)
			}
			ok = true
		} else {
			nearMatches = append(nearMatches, nm...)
		}
	}
	return
}

func (m MirrorValley) ReflectsHorizontal(g Grid, ignoreRow int) (reflectsOn int, nearMatches []Mismatch, ok bool) {
	reflectsOn = -1
	for row := 2; row <= g.maxrow; row++ {
		if row == ignoreRow {
			continue
		}
		mismatched := m.rowMismatch(g, row, row-1)
		if len(mismatched) != 0 {
			continue
		}
		nm, rowOk := m.checkRowReflects(g, row)
		if rowOk {
			if reflectsOn == -1 {
				glog.Infof("Found good reflection around row %d", row)
				reflectsOn = row
			} else {
				glog.Infof("Ignoring additional reflection around row %d", row)
			}
			ok = true
		} else {
			nearMatches = append(nearMatches, nm...)
		}
	}
	return
}

func (m MirrorValley) colMismatch(g Grid, c1, c2 int) (rv []Mismatch) {
	for r := 1; r <= g.maxrow; r++ {
		p1, p2 := Pos{r, c1}, Pos{r, c2}
		v1, v2 := g.C(p1), g.C(p2)
		if v1.(BaseCell).Symbol != v2.(BaseCell).Symbol {
			glog.V(1).Infof("mismatch between %s (%s) and %s (%s)", p1, v1, p2, v2)
			rv = append(rv, Mismatch{a: c1, b: c2, idx: r})
		}
	}
	glog.V(1).Infof("found matching cols %d and %d", c1, c2)
	return
}

func (m MirrorValley) rowMismatch(g Grid, r1, r2 int) (rv []Mismatch) {
	for c := 1; c <= g.maxcol; c++ {
		p1, p2 := Pos{r1, c}, Pos{r2, c}
		v1, v2 := g.C(p1), g.C(p2)
		if v1.(BaseCell).Symbol != v2.(BaseCell).Symbol {
			glog.V(1).Infof("mismatch between %s (%s) and %s (%s)", p1, v1, p2, v2)
			rv = append(rv, Mismatch{a: r1, b: r2, idx: c})
		}
	}
	glog.V(1).Infof("found matching rows %d and %d", r1, r2)
	return
}

type Mismatch struct {
	row int
	col int

	a   int
	b   int
	idx int
}

func (nm Mismatch) String() string {
	var desc string
	var id int
	if nm.row != 0 {
		desc = "row"
		id = nm.row
	} else if nm.col != 0 {
		desc = "col"
		id = nm.col
	} else {
		glog.Fatalf("invalid Mismatch, neither row nor col set: %v", nm)
	}
	return fmt.Sprintf("Mismatch between %ss %d,%d idx %d (mirroring around %s %d)", desc, nm.a, nm.b, nm.idx, desc, id)
}

// checks the reflection point before row matches the rest of the field.
func (m MirrorValley) checkRowReflects(g Grid, row int) (rv []Mismatch, ok bool) {
	ok = true
	step := 1
	for {
		l := row - 1 - step
		r := row + step
		if l < 1 || r > g.maxrow {
			break
		}
		for _, mm := range m.rowMismatch(g, l, r) {
			mm.row = row
			rv = append(rv, mm)
			ok = false
		}
		step++
	}
	return
}

// checks the reflection point before col matches the rest of the field.
func (m MirrorValley) checkColReflects(g Grid, col int) (rv []Mismatch, ok bool) {
	ok = true
	step := 1
	for {
		t := col - 1 - step
		b := col + step
		if t < 1 || b > g.maxcol {
			break
		}
		for _, mm := range m.colMismatch(g, t, b) {
			mm.col = col
			rv = append(rv, mm)
			ok = false
		}
		step++
	}
	return
}

func (m MirrorValley) Score(g Grid) int {
	return m.ScoreIgnoring(g, -1, -1)
}

func (m MirrorValley) ScoreIgnoring(g Grid, iR, iC int) int {
	r, _, ok := m.ReflectsHorizontal(g, iR)
	if ok {
		return (r - 1) * 100
	}
	c, _, ok := m.ReflectsVertical(g, iC)
	if ok {
		return c - 1
	}
	return 0
}

var flip = map[string]string{
	"#": ".",
	".": "#",
}

func (m MirrorValley) FlipCell(c Cell) BaseCell {
	return BaseCell{
		id:     c.(BaseCell).id,
		Symbol: flip[c.(BaseCell).Symbol],
	}
}

func (m MirrorValley) ScoreSmudged(g Grid) int {
	var mm []Mismatch
	r, mmH, _ := m.ReflectsHorizontal(g, -1)
	mm = append(mm, mmH...)
	c, mmV, _ := m.ReflectsVertical(g, -1)
	mm = append(mm, mmV...)
	glog.Infof("Found %d mismatches to try", len(mm))
	scores := []int{}
	for _, t := range mm {
		var aPos, bPos Pos
		if t.col > 0 {
			aPos = Pos{t.idx, t.a}
			bPos = Pos{t.idx, t.b}
		} else if t.row > 0 {
			aPos = Pos{t.a, t.idx}
			bPos = Pos{t.b, t.idx}
		} else {
			glog.Fatalf("invalid mismatch: %#v", t)
		}
		for _, p := range []Pos{aPos, bPos} {
			g2 := g.Copy()
			g2.SetC(p, m.FlipCell(g.C(p)))
			glog.Infof("Flipping %s", p)
			score := m.ScoreIgnoring(g2, r, c)
			if score != 0 {
				glog.Infof(" - Got score of %d", score)
				scores = append(scores, score)
			} else {
				glog.Infof(" -  doesn't give a reflection")
			}
		}
	}
	if len(scores) == 0 {
		glog.Errorf("Found no possible reflections!")
		return 0
	} else if len(scores) > 1 {
		glog.Errorf("Got %d possible scores!, returning first", len(scores))
	}
	return scores[0]
}

func (m MirrorValley) ScoreSum() (rv int) {
	for _, g := range m.Mirrors {
		rv += m.Score(g)
	}
	return
}

func (m MirrorValley) ScoreSmudgedSum() (rv int) {
	for _, g := range m.Mirrors {
		rv += m.ScoreSmudged(g)
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

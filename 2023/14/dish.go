package day14

import (
	"fmt"
	"os"

	"github.com/golang/glog"
)

type DishCell struct {
	BaseCell
}

var _ Cell = &DishCell{}

type Dish struct {
	Grid
}

func NewDish(filename string) (dish Dish, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	dish.Grid = NewGrid[DishCell](f)
	return
}

func (d *Dish) Tilt(dir CardinalDirection) {
	var a, b, aInc, bInc int
	var aEnd func(int) bool
	var bEnd func(int) bool
	var makePos func(int, int) Pos
	switch dir {
	case NORTH:
		a, aInc = 1, 1
		aEnd = func(i int) bool { return i <= d.maxcol }
		b, bInc = 1, 1
		bEnd = func(i int) bool { return i <= d.maxrow }
		makePos = func(row, col int) Pos { return Pos{row, col} }
	case WEST:
		a, aInc = 1, 1
		aEnd = func(i int) bool { return i <= d.maxrow }
		b, bInc = 1, 1
		bEnd = func(i int) bool { return i <= d.maxcol }
		makePos = func(col, row int) Pos { return Pos{row, col} }
	case SOUTH:
		a, aInc = d.maxcol, -1
		aEnd = func(i int) bool { return i > 0 }
		b, bInc = d.maxrow, -1
		bEnd = func(i int) bool { return i > 0 }
		makePos = func(row, col int) Pos { return Pos{row, col} }
	case EAST:
		a, aInc = d.maxrow, -1
		aEnd = func(i int) bool { return i > 0 }
		b, bInc = d.maxcol, -1
		bEnd = func(i int) bool { return i > 0 }
		makePos = func(col, row int) Pos { return Pos{row, col} }
	default:
		glog.Fatalf("cannot tilt %s yet", dir)
	}
	for o := a; aEnd(o); o += aInc {
		next := a + bInc
	FILL:
		for fill := b; bEnd(fill); {
			cell := d.C(makePos(fill, o)).(BaseCell)
			if cell.Symbol == "O" || cell.Symbol == "#" {
				fill += bInc
				if bInc > 0 {
					next = Max(fill+bInc, next)
				} else {
					next = Min(fill+bInc, next)
				}
				continue
			}
			for check := next; bEnd(check); check += bInc {
				cell = d.C(makePos(check, o)).(BaseCell)
				if cell.Symbol == "." {
					next = check
					continue
				}

				if cell.Symbol == "#" {
					fill = check + bInc
					next = fill + bInc
				} else {
					glog.V(1).Infof("%d rolls to %d", makePos(check, o), makePos(fill, o))
					d.SetC(makePos(fill, o), BaseCell{id: makePos(fill, o), Symbol: "O"})
					d.SetC(makePos(check, o), BaseCell{id: makePos(check, o), Symbol: "."})
					next = check + bInc
				}
				continue FILL
			}

			fill += bInc // only hit if column was entirely space to the end.
		}
	}
}

func (d *Dish) Cycle(print ...bool) {
	d.Tilt(NORTH)
	if len(print) > 0 && print[0] {
		fmt.Println("Tilted North")
		d.PrintNumbered()
	}
	d.Tilt(WEST)
	if len(print) > 0 && print[0] {
		fmt.Println("Tilted West")
		d.PrintNumbered()
	}
	d.Tilt(SOUTH)
	if len(print) > 0 && print[0] {
		fmt.Println("Tilted South")
		d.PrintNumbered()
	}
	d.Tilt(EAST)
	if len(print) > 0 && print[0] {
		fmt.Println("Tilted East")
		d.PrintNumbered()
	}
}

func (d *Dish) Load() (sum int) {
	for r := 1; r <= d.maxrow; r++ {
		load := d.maxrow - r + 1 // 1 based, vs zero
		for c := 1; c <= d.maxcol; c++ {
			cell := d.C(Pos{r, c}).(BaseCell)
			if cell.Symbol == "O" {
				sum += load
			}
		}
	}
	return
}

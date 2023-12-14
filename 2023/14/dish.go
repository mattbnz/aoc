package day14

import (
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

func (d *Dish) Tilt() {
	for c := 1; c <= d.maxcol; c++ {
		next := 2
	FILL:
		for fill := 1; fill <= d.maxrow; {
			cell := d.C(Pos{fill, c}).(BaseCell)
			if cell.Symbol == "O" || cell.Symbol == "#" {
				fill++
				next = Max(fill+1, next)
				continue
			}
			for check := next; check <= d.maxrow; check++ {
				cell = d.C(Pos{check, c}).(BaseCell)
				if cell.Symbol == "." {
					next = check
					continue
				}

				if cell.Symbol == "#" {
					fill = check + 1
					next = fill + 1
				} else {
					glog.V(1).Infof("%d falls to %d", Pos{check, c}, Pos{fill, c})
					d.SetC(Pos{fill, c}, BaseCell{id: Pos{fill, c}, Symbol: "O"})
					d.SetC(Pos{check, c}, BaseCell{id: Pos{check, c}, Symbol: "."})
					next = check + 1
				}
				continue FILL
			}

			fill++ // only hit if column was entirely space to the end.
		}
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

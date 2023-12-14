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

type Key struct {
	GridKey string
	Tilt    CardinalDirection
}

type Dish struct {
	Grid

	Cache     map[Key]Grid
	hit, miss int
	checks    int
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

// TODO - move this to Grid?
func (d *Dish) Key() (key string) {
	d.Each(func(_ Pos, c Cell) bool {
		key += c.(BaseCell).Symbol
		return true
	})
	return
}

func (d *Dish) RowKey(row int) (key string) {
	for c := 0; c <= d.maxcol; c++ {
		key += d.C(Pos{row, c}).(BaseCell).Symbol
	}
	return
}
func (d *Dish) ColKey(col int) (key string) {
	for r := 0; r <= d.maxrow; r++ {
		key += d.C(Pos{r, col}).(BaseCell).Symbol
	}
	return
}

func (d *Dish) SliceKey(a, b, bInc int, bEnd func(int) bool, makePos func(int, int) Pos) (key string) {
	for i := b; bEnd(i); i += bInc {
		key += d.C(makePos(i, a)).(BaseCell).Symbol
	}
	return
}

func (d *Dish) Tilt(dir CardinalDirection) {
	/*key := Key{d.Key(), dir}
	if result, found := d.Cache[key]; found {
		d.Grid = result.Copy()
		glog.V(1).Infof("Tilt(%s) from %s found in cache", dir, key)
		d.hit++
		return
	}
	d.miss++*/

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
		key := Key{d.SliceKey(o, b, bInc, bEnd, makePos), NO_DIRECTION}
		if _, found := d.Cache[key]; found {
			d.hit++
		}
		d.miss++
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
				d.checks++
				cell = d.C(makePos(check, o)).(BaseCell)
				if cell.Symbol == "." {
					next = check
					continue
				}

				if cell.Symbol == "#" {
					fill = check + bInc
					next = fill + bInc
				} else {
					glog.V(2).Infof("%d rolls to %d", makePos(check, o), makePos(fill, o))
					d.SetC(makePos(fill, o), BaseCell{id: makePos(fill, o), Symbol: "O"})
					d.SetC(makePos(check, o), BaseCell{id: makePos(check, o), Symbol: "."})
					next = check + bInc
				}
				continue FILL
			}

			fill += bInc // only hit if column was entirely space to the end.
		}
	}

	/*if d.Cache == nil {
		d.Cache = make(map[Key]Grid)
	}
	glog.V(1).Infof("Tilt(%s) from %s added to cache", dir, key)
	d.Cache[key] = d.Grid.Copy()*/
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

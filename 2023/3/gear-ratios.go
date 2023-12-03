// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 3.
// Gear Ratios

package day3

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

type Pos struct {
	row, col int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d,%d", p.row, p.col)
}

type Cell struct {
	id   Pos
	grid *Grid

	Value  int
	Symbol rune
	Prev   *Cell
	Next   *Cell
}

func (c Cell) String() string {
	return fmt.Sprintf("%s", c.id)
}

func (c Cell) NumberStarts() (v int, hasSymbol bool) {
	if c.Prev != nil {
		return -1, false
	}
	return c.Number()
}

func (c Cell) Number() (v int, hasSymbol bool) {
	if c.Value == -1 || c.Symbol != 0 {
		return -1, false
	}
	if c.Prev != nil {
		return c.Prev.Number()
	}
	v, _, hasSymbol = c.getDigit()
	return
}

func (c Cell) getDigit() (value, pos int, hasSymbol bool) {
	if c.Value == -1 || c.Symbol != 0 {
		return -1, -1, false
	}
	if c.Next == nil {
		return c.Value, 1, c.hasSymbol()
	}
	v, pos, hasSymbol := c.Next.getDigit()
	if !hasSymbol {
		hasSymbol = c.hasSymbol()
	}
	return (int(math.Pow(10, float64(pos))) * c.Value) + v, pos + 1, hasSymbol
}

func (c Cell) hasSymbol() bool {
	for _, rD := range []int{-1, 0, 1} {
		for _, cD := range []int{-1, 0, 1} {
			if rD == 0 && cD == 0 {
				continue
			}
			oCell := c.grid.Cell(Pos{c.id.row + rD, c.id.col + cD})
			if oCell == nil {
				continue
			}
			if oCell.Symbol != 0 {
				return true
			}
		}
	}
	return false
}

func (c Cell) Numbers() (rv []int) {
	touches := map[int]bool{}
	for _, rD := range []int{-1, 0, 1} {
		for _, cD := range []int{-1, 0, 1} {
			if rD == 0 && cD == 0 {
				continue
			}
			oCell := c.grid.Cell(Pos{c.id.row + rD, c.id.col + cD})
			if oCell == nil {
				continue
			}
			if n, _ := oCell.Number(); n != -1 {
				touches[n] = true
			}
		}
	}
	for n, _ := range touches {
		rv = append(rv, n)
	}
	return
}

type Grid struct {
	cells [][]*Cell
}

func (g *Grid) ParseRow(r string) {
	var lastDigit *Cell
	row := []*Cell{}
	for i, cS := range r {
		c := &Cell{id: Pos{len(g.cells), i}, grid: g}
		if cS == '.' {
			c.Value = -1
			lastDigit = nil
		} else if cS >= '0' && cS <= '9' {
			n, err := strconv.Atoi(string(cS))
			if err != nil {
				log.Fatalf("Invalid cell value at (%d, %d) - %c: %v", len(g.cells), i, cS, err)
			}
			c.Value = n
			c.Prev = lastDigit
			if c.Prev != nil {
				c.Prev.Next = c
			}
			lastDigit = c
		} else {
			c.Symbol = cS
			lastDigit = nil
		}
		row = append(row, c)
	}
	g.cells = append(g.cells, row)
}

func (g *Grid) Print() {
	for _, cols := range g.cells {
		for _, cell := range cols {
			if cell.Symbol != 0 {
				fmt.Printf("%c", cell.Symbol)
			} else if cell.Value != -1 {
				fmt.Printf("%d", cell.Value)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
}

func (g *Grid) FindPartNumbers() (rv []int) {
	for _, cols := range g.cells {
		for _, cell := range cols {
			if n, hasSymbol := cell.NumberStarts(); n != -1 && hasSymbol {
				log.Printf("%s: %d", cell, n)
				rv = append(rv, n)
			}
		}
	}
	return
}

func (g *Grid) FindGearRatios() (rv []int) {
	for _, cols := range g.cells {
		for _, cell := range cols {
			if cell.Symbol == '*' {
				n := cell.Numbers()
				if len(n) == 2 {
					rv = append(rv, n[0]*n[1])
				}
			}
		}
	}
	return
}

// returns a cell from the grid (or nil, if pos is invalid)
func (g *Grid) Cell(p Pos) *Cell {
	if p.row < 0 || p.row > len(g.cells)-1 || p.col < 0 || p.col > len(g.cells[0])-1 {
		return nil
	}
	return g.cells[p.row][p.col]
}

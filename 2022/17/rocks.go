// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 17, Puzzle 1.
// Pyroclastic Flow - tetris?

package main

import (
	"errors"
	"fmt"
	"log"
)

const WIDTH = 7

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// row 0 == bottom
type Pos struct {
	row, col int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d,%d", p.row, p.col)
}

const AIR = 0
const ROCK = 1
const FALLING_ROCK = 2

const LEFT = -1
const RIGHT = 1

type Column struct {
	// last final state, doesn't include current piece
	c      map[Pos]int
	toprow int // topmost row in c

	// the piece currently falling
	piece Piece
}

func NewColumn() Column {
	return Column{c: map[Pos]int{}}
}

func (c Column) Print() {
	c.print()
}

func (c Column) print() {
	top := c.toprow
	if c.piece.top != 0 {
		top = c.piece.top
	}

	for row := top; row >= 0; row-- {
		for col := 0; col < WIDTH+2; col++ {
			// borders
			if row == 0 {
				if col == 0 || col == WIDTH+1 {
					fmt.Printf("+")
				} else {
					fmt.Printf("-")
				}
				continue
			}
			if col == 0 || col == WIDTH+1 {
				fmt.Printf("|")
				continue
			}
			// playing space
			t := c.C(Pos{row, col})
			if t == ROCK {
				fmt.Printf("#")
			} else if t == FALLING_ROCK {
				fmt.Printf("@")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
}

func (c Column) C(p Pos) int {
	pxl := AIR
	if c.piece.top != 0 {
		t, err := c.piece.C(p)
		if err == nil {
			pxl = t
		}
	}
	if pxl != AIR {
		if c.c[p] != AIR {
			log.Fatal("Unexpected collision at ", p, " with ", c.piece)
		}
		return pxl
	}
	return c.c[p]
}

func (c *Column) SetC(p Pos, v int) {
	if c.c[p] != AIR {
		log.Fatal("Overwriting AIR at ", p, " with", v)
	}
	c.c[p] = v
	c.toprow = Max(c.toprow, p.row)
}

func (c *Column) Push(dir int) {
	np, err := c.piece.Copy(dir)
	if err != nil {
		return
	}
	// TODO check for rock intersections
	c.piece = np
}

func (c *Column) Drop() bool {
	oldtop := c.piece.top
	c.piece.top--
	// TODO check for rock intersections
	if c.piece.top-c.piece.height == 0 {
		// hit the floor!
		c.piece.top = oldtop // move back up one
		c.Merge()
		return true
	}
	return false
}

func (c *Column) Merge() {
	for _, pos := range c.piece.AbsPixels() {
		c.SetC(pos, ROCK)
	}
	c.piece = Piece{}
}

// Inserts a new piece at it's default location
func (c *Column) New(p Piece) {
	p.top = c.toprow + 4 + p.height
	c.piece = p
}

// Positions within piece are relative to the piece bottom-left (0,0), but arguments coming
// in will be in column position, and need to be trasnformed.
type Piece struct {
	name   string      // convenience
	top    int         // row of column this rock's top is at
	height int         // height of piece
	pixels map[Pos]int // set of pixels that make up this rock, relative to bottom,left(0,0)
}

var OutsidePiece = errors.New("outside piece")
var InvalidMove = errors.New("invalid move")

func (p Piece) String() string {
	return fmt.Sprintf("%s@%d", p.name, p.top)
}

// Convert an absolute (column) position, into a relative position in this item
func (p Piece) AbsToRel(pos Pos) (Pos, error) {
	if pos.row >= (p.top-p.height) && pos.row <= p.top {
		return Pos{row: (pos.row - p.top) + p.height, col: pos.col - 1}, nil
	}
	return Pos{}, OutsidePiece
}

// Convert an relative position in this item into a absolute (column) posotion
func (p Piece) RelToAbs(pos Pos) Pos {
	return Pos{row: (pos.row + p.top) - p.height, col: pos.col + 1}
}

func (p *Piece) C(pos Pos) (int, error) {
	rpos, err := p.AbsToRel(pos)
	if err != nil {
		return -1, err
	}
	return p.pixels[rpos], nil
}

func (p Piece) AbsPixels() []Pos {
	rv := []Pos{}
	for pos := range p.pixels {
		rv = append(rv, p.RelToAbs(pos))
	}
	return rv
}

// Returns a copy of iteslf moved col places left -/right +, or err if
// that movement would go outside the bounds
func (p Piece) Copy(col int) (Piece, error) {
	np := Piece{name: p.name, height: p.height, top: p.top}
	nPxls := map[Pos]int{}
	for pos, v := range p.pixels {
		pos.col += col
		if pos.col < 0 || pos.col > WIDTH {
			return Piece{}, InvalidMove
		}
		nPxls[pos] = v
	}
	np.pixels = nPxls
	return np, nil
}

func NewPiece(n string, p []Pos) Piece {
	np := Piece{name: n}
	pixels := map[Pos]int{}
	h := 0
	for _, t := range p {
		pixels[t] = FALLING_ROCK
		h = Max(h, t.row)
	}
	np.pixels = pixels
	np.height = h
	return np
}

func main() {
	pieces := []Piece{
		NewPiece("hline", []Pos{{0, 2}, {0, 3}, {0, 4}, {0, 5}}),
		NewPiece("plus", []Pos{{0, 3}, {1, 2}, {1, 3}, {1, 4}, {2, 3}}),
		NewPiece("badl", []Pos{{0, 2}, {0, 3}, {0, 4}, {1, 4}, {2, 4}}),
		NewPiece("vline", []Pos{{0, 2}, {1, 2}, {2, 2}, {3, 2}}),
		NewPiece("square", []Pos{{0, 2}, {0, 3}, {1, 2}, {1, 3}}),
	}

	for _, p := range pieces {
		col := NewColumn()
		fmt.Println(p)
		col.New(p)
		col.Print()
		col.Push(RIGHT)
		col.Drop()
		col.Merge()
		col.Print()
		fmt.Println("")
	}
}

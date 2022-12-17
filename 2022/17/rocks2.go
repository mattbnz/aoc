// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 17, Puzzle 2.
// Pyroclastic Flow - giant tetris...

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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

	for row := top; row >= Max(c.toprow-10, 0); row-- {
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
			t := c.C(Pos{row, col}, true)
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
	fmt.Println()
}

func (c Column) C(p Pos, usePiece bool) int {
	pxl := AIR
	if usePiece {
		if c.piece.top != 0 {
			t, err := c.piece.C(p)
			if err == nil {
				pxl = t
			}
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
	if c.TouchesRock(np) {
		return
	}
	c.piece = np
}

func (c *Column) Drop() bool {
	np := c.piece
	np.top--
	if c.TouchesRock(np) {
		// can't drop, merge where is and stop
		c.Merge()
		return true
	}
	if np.top-np.height == 0 {
		// hit the floor!
		c.Merge()
		return true
	}
	c.piece = np
	return false
}

func (c *Column) TouchesRock(p Piece) bool {
	for _, pos := range p.AbsPixels() {
		if c.C(pos, false) != AIR {
			return true
		}
	}
	return false
}

func (c *Column) Merge() {
	//c.Print()
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
		if pos.col < 0 || pos.col >= WIDTH {
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

func ParseJets(in string) []int {
	rv := []int{}
	for n, c := range strings.TrimSpace(in) {
		if c == '<' {
			rv = append(rv, LEFT)
		} else if c == '>' {
			rv = append(rv, RIGHT)
		} else {
			log.Fatalf("bad jet at char %d: %c", n, c)
		}
	}
	return rv
}

func main() {
	pieces := []Piece{
		NewPiece("hline", []Pos{{0, 2}, {0, 3}, {0, 4}, {0, 5}}),
		NewPiece("plus", []Pos{{0, 3}, {1, 2}, {1, 3}, {1, 4}, {2, 3}}),
		NewPiece("badl", []Pos{{0, 2}, {0, 3}, {0, 4}, {1, 4}, {2, 4}}),
		NewPiece("vline", []Pos{{0, 2}, {1, 2}, {2, 2}, {3, 2}}),
		NewPiece("square", []Pos{{0, 2}, {0, 3}, {1, 2}, {1, 3}}),
	}

	reader := bufio.NewReader(os.Stdin)
	jetS, _ := reader.ReadString('\n')
	jets := ParseJets(jetS)
	if len(jets) < 1 {
		log.Fatal("Didn't get jets")
	}

	wanted := 1000000000000

	col := NewColumn()
	move := 0
	heights := []int{}      // height after each rock
	patternsizes := []int{} // sizes we're still evaluating for patterns
	firstmatch := map[int]int{}
	matches := map[int]int{}
	lastrock := 0

rocks:
	for rock := 0; rock < wanted; rock++ {
		if rock > len(pieces) {
			patternsizes = append(patternsizes, rock)
			firstmatch[rock] = -1
		}
		col.New(pieces[rock%len(pieces)])
		//col.Print()
		for {
			col.Push(jets[move%len(jets)])
			move++
			if col.Drop() {
				break
			}
			//col.Print()
		}
		//fmt.Println()
		heights = append(heights, col.toprow)
		lastrock = rock
		// Check for any pattern
		newsizes := []int{}

		//patterns:
		for _, l := range patternsizes {
			if rock < l*2 {
				newsizes = append(newsizes, l)
				continue
			}
			// Drop anything that's a multiple of a pattern that's been found
			/*for fl, at := range firstmatch {
				if at == -1 {
					continue
				}
				if l > fl && l%fl == 0 {
					fmt.Printf("%d: len=%d DROPPING as multiple of found pattern %d\n", rock, l, fl)
					continue patterns
				}
			}*/
			h1 := rock
			h2 := rock - l
			h3 := rock - (l * 2)
			delta1 := heights[h1] - heights[h2]
			delta2 := heights[h2] - heights[h3]
			if delta1 == delta2 {
				fmt.Printf("%d: len=%d,delta=%d FOUND h%d=%d +%d = h%d=%d %d = h%d=%d\n",
					rock, l, delta1, h3, heights[h3], delta2, h2, heights[h2], delta1, h1, heights[h1])
				if firstmatch[l] == -1 {
					firstmatch[l] = h3
				}
				matches[l]++
				if matches[l] > 200 {
					// 200 consecutive matches, enough to be the winner?
					patternsizes = []int{l}
					break rocks
				}
				newsizes = append(newsizes, l)
				continue
			}
			if firstmatch[l] != -1 {
				// This pattern had matched, but now doesn't, so can't be the answer
				fmt.Printf("%d: len=%d DROPPING didn't match after previous match\n", rock, l)
				firstmatch[l] = -1
				continue
			}
			// keep considering.
			newsizes = append(newsizes, l)
		}
		patternsizes = newsizes
		/*if len(patternsizes) == 1 {
			fmt.Println("Bailing with one remaining pattern len=", patternsizes[0])
			break
		}*/
	}

	if len(patternsizes) != 1 {
		fmt.Println("top row is ", col.toprow)
	} else {
		pLen := patternsizes[0]
		pDelta := heights[lastrock] - heights[lastrock-pLen]
		baseRow := firstmatch[pLen]
		fmt.Println("baseRow ", baseRow)
		gap := (wanted - 1) - baseRow
		fmt.Println("gap, ncycles ", gap, gap/pLen)
		cycleRows := pLen * (gap / pLen)
		fmt.Println("cycleHeight ", cycleRows)
		baseHeightRow := (wanted - 1) - cycleRows
		fmt.Println("baseheightrow", baseHeightRow)
		baseHeight := heights[baseHeightRow]
		fmt.Println("baseheight", baseHeight)
		expected := baseHeight + (pDelta * (gap / pLen))
		fmt.Printf("Given cycle of len=%d,delta=%d, expected height after %d is: %d\n",
			pLen, pDelta, wanted, expected)
	}

}

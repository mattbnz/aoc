// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 24, Puzzle 1.
// Blizzard Basin. Another search...

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func Max(a, b int) int {
	if a == -1 {
		return b
	}
	if b == -1 {
		return a
	}
	if a > b {
		return a
	}
	return b
}

func Int(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("%s is not an int: %v", s, err)
	}
	return v
}

type Pos struct {
	row, col int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d,%d", p.row, p.col)
}

type Content int

const (
	OPEN = iota
	WALL
	BLIZZARD_LEFT
	BLIZZARD_RIGHT
	BLIZZARD_UP
	BLIZZARD_DOWN
)

var StrToContent = map[string]Content{
	"#": WALL,
	".": OPEN,
	">": BLIZZARD_RIGHT,
	"<": BLIZZARD_LEFT,
	"^": BLIZZARD_UP,
	"v": BLIZZARD_DOWN,
}
var ContentToStr = map[Content]string{}

func init() {
	for s, c := range StrToContent {
		ContentToStr[c] = s
	}
}

type Grid struct {
	c              map[Pos]Content
	maxrow, maxcol int
}

func (g Grid) Print() {
	for row := 1; row <= g.maxrow; row++ {
		for col := 1; col <= g.maxcol; col++ {
			fmt.Printf(ContentToStr[g.C(Pos{row, col})])
		}
		fmt.Println()
	}
}

func (g Grid) C(p Pos) Content {
	return g.c[p] // default val if missing == VOID
}

func (g *Grid) SetC(p Pos, v Content) {
	g.c[p] = v
	g.maxrow = Max(g.maxrow, p.row)
	g.maxcol = Max(g.maxcol, p.col)
}

func NewGrid() Grid {
	g := Grid{c: map[Pos]Content{}}
	g.maxrow = -1
	g.maxcol = -1
	return g
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	g := NewGrid()

	row := 1
	for s.Scan() {
		for col, c := range s.Text() {
			p := Pos{row, col + 1}
			g.SetC(p, StrToContent[string(c)])
		}
		row++
	}

	g.Print()
}

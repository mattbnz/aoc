// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 22, Puzzle 1.
// Monkey Map.  Path following.

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

const VOID = 0
const OPEN = 1
const WALL = 2

type Grid struct {
	c              map[Pos]int
	maxrow, maxcol int
}

func (g Grid) Print() {
	for row := 1; row <= g.maxrow; row++ {
		for col := 1; col <= g.maxcol; col++ {
			c := g.C(Pos{row, col})
			if c == VOID {
				fmt.Printf(" ")
			} else if c == OPEN {
				fmt.Printf(".")
			} else if c == WALL {
				fmt.Printf("#")
			}
		}
		fmt.Println()
	}
}

func (g Grid) C(p Pos) int {
	return g.c[p] // default val if missing == VOID
}

func (g *Grid) SetC(p Pos, v int) {
	g.c[p] = v
	g.maxrow = Max(g.maxrow, p.row)
	g.maxcol = Max(g.maxcol, p.col)
}

func NewGrid() Grid {
	g := Grid{c: map[Pos]int{}}
	g.maxrow = -1
	g.maxcol = -1
	return g
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	g := NewGrid()

	var firstPos *Pos
	row := 1
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		for col, c := range s.Text() {
			if c == ' ' {
				continue
			}
			if c == '.' {
				p := Pos{row, col + 1}
				g.SetC(p, OPEN)
				if firstPos == nil {
					firstPos = &p
				}
			} else if c == '#' {
				g.SetC(Pos{row, col + 1}, WALL)
			} else {
				log.Fatalf("Bad input at %d, unknwon char '%c': %s", col, c, s.Text())
			}
		}
		row++
	}
	if !s.Scan() {
		log.Fatal("Couldn't read instruction line")
	}
	instructions := s.Text()
	g.Print()

	fmt.Println(instructions)
	fmt.Println(firstPos)
}

// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 22, Puzzle 1.
// Monkey Map.  Path following.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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

const S_LEFT = "L"
const S_RIGHT = "R"

var DIR_S = map[Direction]string{
	D_RIGHT: ">",
	D_DOWN:  "v",
	D_LEFT:  ">",
	D_UP:    "^",
}

var TURN_MAP = map[Direction]map[string]Direction{
	D_RIGHT: {S_LEFT: D_UP, S_RIGHT: D_DOWN},
	D_DOWN:  {S_LEFT: D_RIGHT, S_RIGHT: D_LEFT},
	D_LEFT:  {S_LEFT: D_DOWN, S_RIGHT: D_UP},
	D_UP:    {S_LEFT: D_LEFT, S_RIGHT: D_RIGHT},
}

type Direction int

const (
	D_RIGHT Direction = iota
	D_DOWN
	D_LEFT
	D_UP
)

func (d Direction) String() string {
	return DIR_S[d]
}

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

func (g *Grid) Next(from Pos, heading Direction) Pos {
	if heading == D_RIGHT {
		return Pos{from.row, from.col + 1}
	} else if heading == D_DOWN {
		return Pos{from.row + 1, from.col}
	} else if heading == D_LEFT {
		return Pos{from.row, from.col - 1}
	} else if heading == D_UP {
		return Pos{from.row - 1, from.col}
	}
	return Pos{}
}

func (g *Grid) Wrap(from Pos, heading Direction) Pos {
	back := TURN_MAP[TURN_MAP[heading][S_LEFT]][S_LEFT] // left twice == reverse.
	at := from
	for g.C(at) != VOID {
		at = g.Next(at, back)
	}
	// at this point we're off the map on the otherside, so return one back in the other direction
	// for the first map spot
	return g.Next(at, heading)
}

func (g *Grid) Nav(from Pos, heading Direction, steps int) Pos {
	at := from
	for i := 0; i < steps; i++ {
		next := g.Next(at, heading)
		if g.C(next) == OPEN {
			at = next
			continue
		}
		if g.C(next) == WALL {
			return at
		}
		if g.C(next) == VOID {
			next = g.Wrap(at, heading)
			if g.C(next) == OPEN {
				at = next
				continue
			}
			if g.C(next) == WALL {
				return at
			}
		}
	}
	return at
}

func NewGrid() Grid {
	g := Grid{c: map[Pos]int{}}
	g.maxrow = -1
	g.maxcol = -1
	return g
}

var I_RE = regexp.MustCompile(`([LR])?(\d+)`)

func main() {
	s := bufio.NewScanner(os.Stdin)

	g := NewGrid()

	var firstPos Pos
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
				if firstPos.row == 0 {
					firstPos = p
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
	fmt.Println()
	fmt.Println(instructions)
	fmt.Println()

	pos := firstPos
	heading := D_RIGHT
	for step, i := range I_RE.FindAllStringSubmatch(instructions, -1) {
		fmt.Printf("* Step %d, at %s going %s\n", step, pos, heading)
		if i[1] != "" {
			heading = TURN_MAP[heading][i[1]]
			fmt.Printf("  - Turn %s, Heading: %s\n", i[1], heading)
		}
		fmt.Printf("  - Moving %s steps\n", i[2])
		pos = g.Nav(pos, heading, Int(i[2]))
		fmt.Println()
	}
	fmt.Printf("Finished at: %s going %s\n", pos, heading)
	fmt.Printf("Password: %d * 1000 + %d * 4 + %d = %d\n", pos.row, pos.col, heading, pos.row*1000+pos.col*4+int(heading))
}

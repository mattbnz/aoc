// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 18, Puzzle 2.
// Boiling Boulders - external surface area...

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
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
	x, y, z int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d,%dm%d", p.x, p.y, p.z)
}

type Side struct {
	covered bool
	checked bool
}

type Grid map[Pos]*Cube

// side constants, for indexing into sides slice of a Cube
const S_LEFT = 0
const S_RIGHT = 1

const S_FRONT = 2
const S_BACK = 3

const S_BOTTOM = 4
const S_TOP = 5

// axis constants for indexing an axis in the side map below
const A_X = 0
const A_Y = 1
const A_Z = 2

// Map for each axis of an offset, to the side it represents
var sidemap = map[int]map[int]int{
	A_X: {-1: S_LEFT, 1: S_RIGHT},
	A_Y: {-1: S_FRONT, 1: S_BACK},
	A_Z: {-1: S_BOTTOM, 1: S_TOP},
}

type Cube struct {
	pos Pos

	// state of sides, indexed by SIDE_ consts above
	sides map[int]*Side
}

// Wrapper around a grid lookup that looks for the cube at specified axis offset
func (c *Cube) Neighbour(g *Grid, axis, dir int) (*Cube, bool) {
	var n *Cube
	var found bool
	if axis == A_X {
		n, found = (*g)[Pos{c.pos.x + dir, c.pos.y, c.pos.z}]
	} else if axis == A_Y {
		n, found = (*g)[Pos{c.pos.x, c.pos.y + dir, c.pos.z}]
	} else if axis == A_Z {
		n, found = (*g)[Pos{c.pos.x, c.pos.y, c.pos.z + dir}]
	}
	return n, found
}

func (c *Cube) CheckSide(g *Grid, axis, dir int) {
	sI := sidemap[axis][dir]

	if c.sides[sI].checked {
		return // neighbour already marked this side, early exit.
	}
	n, found := c.Neighbour(g, axis, dir)
	if !found {
		// empty space == default, so just return after marking checked
		c.sides[sI].checked = true
		return
	}
	// mark us and neighbours opposite side as covered
	c.sides[sI].covered = true
	c.sides[sI].checked = true
	nI := sidemap[axis][dir*-1]
	n.sides[nI].covered = true
	n.sides[nI].checked = true
}

func (c *Cube) CheckSides(g *Grid) {
	for axis, sides := range sidemap {
		for dir, _ := range sides {
			c.CheckSide(g, axis, dir)
		}
	}
}

func (c *Cube) SurfaceArea() int {
	a := 0
	for _, side := range c.sides {
		if !side.covered {
			a++
		}
	}
	return a
}

func NewCube(cS string) *Cube {
	parts := strings.Split(cS, ",")
	if len(parts) != 3 {
		log.Fatal("bad cube dimensions", cS)
	}
	p := Pos{Int(parts[0]), Int(parts[1]), Int(parts[2])}
	c := Cube{pos: p, sides: map[int]*Side{}}
	for _, a := range sidemap {
		for _, s := range a {
			c.sides[s] = &Side{}
		}
	}
	return &c
}

func Print(g *Grid, xM, yM, zM int) {
	for z := 1; z <= zM; z++ {
		for x := 1; x <= xM; x++ {
			for y := 1; y <= yM; y++ {
				_, present := (*g)[Pos{x, y, z}]
				if present {
					fmt.Printf("#")
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func PrintSurface(s map[Pos]int, xM, yM, zM int) {
	for x := 1; x <= xM; x++ {
		for y := 1; y <= yM; y++ {
			_, present := s[Pos{x, y, 0}]
			if present {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// THis is buggy, doesn't consider the top/bottoms exposed
// e.g. at 1,1,2 and 1,3,3
func ProjectZ(g *Grid, xM, yM, zM int) int {
	surface := map[Pos]int{}
	for z := 1; z <= zM; z++ {
		for x := 1; x <= xM; x++ {
			for y := 1; y <= yM; y++ {
				_, present := (*g)[Pos{x, y, z}]
				if present {
					surface[Pos{x, y, 0}] = 1
				}
			}
		}
		PrintSurface(surface, xM, yM, zM)
	}
	return len(surface)
}

func ProjectY(g *Grid, xM, yM, zM int) int {
	sum := 0
	for z := 1; z <= zM; z++ {
		surface := map[Pos]int{}
		for x := 1; x <= xM; x++ {
			for y := 1; y <= yM; y++ {
				_, present := (*g)[Pos{x, y, z}]
				if present {
					surface[Pos{x, 1, 0}] = 1
				}
			}
		}
		PrintSurface(surface, xM, yM, zM)
		sum += len(surface)
	}
	return sum
}

func ProjectX(g *Grid, xM, yM, zM int) int {
	sum := 0
	for z := 1; z <= zM; z++ {
		surface := map[Pos]int{}
		for x := 1; x <= xM; x++ {
			for y := 1; y <= yM; y++ {
				_, present := (*g)[Pos{x, y, z}]
				if present {
					surface[Pos{1, y, 0}] = 1
				}
			}
		}
		PrintSurface(surface, xM, yM, zM)
		sum += len(surface)
	}
	return sum
}

func main() {
	grid := Grid{}

	s := bufio.NewScanner(os.Stdin)

	xM, yM, zM := 1, 1, 1
	for s.Scan() {
		c := NewCube(s.Text())
		grid[c.pos] = c
		xM = Max(xM, c.pos.x)
		yM = Max(yM, c.pos.y)
		zM = Max(zM, c.pos.z)
	}
	fmt.Printf("Loaded %d cubes into grid\n", len(grid))
	fmt.Println(xM, yM, zM)
	fmt.Println("Grid:")
	Print(&grid, xM, yM, zM)

	fmt.Println("Bottom")
	b := ProjectZ(&grid, xM, yM, zM)
	fmt.Println(b)
	fmt.Println("Left")
	l := ProjectY(&grid, xM, yM, zM)
	fmt.Println(l)
	fmt.Println("Front")
	f := ProjectX(&grid, xM, yM, zM)
	fmt.Println(f)
	// missing 8 on sample due to bug above!
	fmt.Println((b * 2) + (l * 2) + (f * 2))
}

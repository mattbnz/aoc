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
	A_Y: {-1: S_LEFT, 1: S_RIGHT},
	A_X: {-1: S_FRONT, 1: S_BACK},
	A_Z: {-1: S_BOTTOM, 1: S_TOP},
}

const AIR = 0
const LAVA = 1
const STEAM = 2

type Cube struct {
	pos Pos
	is  int

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

func (c *Cube) SideTouching(g *Grid, axis, dir int, what int) bool {
	n, found := c.Neighbour(g, axis, dir)
	if !found {
		return false
	}
	return n.is == what
}

func (c *Cube) ExpandIfTouching(g *Grid, what int) {
	for axis, sides := range sidemap {
		for dir, _ := range sides {
			if c.SideTouching(g, axis, dir, what) {
				c.is = what
				return
			}
		}
	}
	return
}

func (c *Cube) ExternalSurfaceArea(g *Grid) int {
	a := 0
	for axis, sides := range sidemap {
		for dir, _ := range sides {
			if c.SideTouching(g, axis, dir, STEAM) {
				a++
			}
		}
	}
	return a
}

func NewCubeString(cS string, is int) *Cube {
	parts := strings.Split(cS, ",")
	if len(parts) != 3 {
		log.Fatal("bad cube dimensions", cS)
	}
	p := Pos{Int(parts[0]), Int(parts[1]), Int(parts[2])}
	return NewCube(p, is)
}

func NewCube(p Pos, is int) *Cube {
	c := Cube{pos: p, is: is, sides: map[int]*Side{}}
	for _, a := range sidemap {
		for _, s := range a {
			c.sides[s] = &Side{}
		}
	}
	return &c
}

func Print(g *Grid, min, xM, yM, zM int) {
	for z := min; z <= zM; z++ {
		PrintLayer(g, min, xM, yM, z)
	}
}
func PrintLayer(g *Grid, min, xM, yM, z int) {
	for x := min; x <= xM; x++ {
		for y := min; y <= yM; y++ {
			c, present := (*g)[Pos{x, y, z}]
			if present && c.is == LAVA {
				fmt.Printf("%d", c.ExternalSurfaceArea(g))
			} else if present && c.is == STEAM {
				fmt.Printf("*")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// Sets up the steam 'border' (using co-ords 0, aX+1 in all dimensions) to ensure we have the
// starting steam to expand from. Also convert any 'air' cell (internal or external) to an
// explicit cube (e.g. grid becomes complete, not sparse) of AIR to make later computations easier
func AddSteam(g *Grid, xM, yM, zM int) {
	for z := 0; z <= zM+1; z++ {
		for x := 0; x <= xM+1; x++ {
			for y := 0; y <= yM+1; y++ {
				p := Pos{x, y, z}
				if z == 0 || z == zM+1 || x == 0 || x == xM+1 || y == 0 || y == yM+1 {
					c := NewCube(p, STEAM)
					(*g)[c.pos] = c
				} else if _, present := (*g)[p]; !present {
					c := NewCube(p, AIR)
					(*g)[c.pos] = c
				}
			}
		}
	}
}

func ExpandSteam(g *Grid, xS, yS, zS, step, xE, yE, zE int) {
	for z := zS; z != zE; z += step {
		for x := xS; x != xE; x += step {
			for y := yS; y != yE; y += step {
				p := Pos{x, y, z}
				c := (*g)[p]
				if c.is != AIR {
					continue
				}
				c.ExpandIfTouching(g, STEAM)
			}
		}
		PrintLayer(g, 0, Max(xS, xE), Max(yS, yE), z)
	}
}

func main() {
	grid := Grid{}

	s := bufio.NewScanner(os.Stdin)

	xM, yM, zM := 1, 1, 1
	for s.Scan() {
		c := NewCubeString(s.Text(), LAVA)
		grid[c.pos] = c
		xM = Max(xM, c.pos.x)
		yM = Max(yM, c.pos.y)
		zM = Max(zM, c.pos.z)
	}
	fmt.Printf("Loaded %d cubes into grid\n", len(grid))
	fmt.Println(xM, yM, zM)
	fmt.Println("Grid:")
	Print(&grid, 1, xM, yM, zM)
	fmt.Println()
	AddSteam(&grid, xM, yM, zM)

	fmt.Println("Expanding positively")
	ExpandSteam(&grid, 1, 1, 1, 1, xM+1, yM+1, zM+1)
	fmt.Println("Expanding negatively")
	ExpandSteam(&grid, xM+1, yM+1, zM+1, -1, 0, 0, 0)

	fmt.Println("Steamy Grid:")
	Print(&grid, 0, xM+1, yM+1, zM+1)

	sum := 0
	for _, c := range grid {
		if c.is == LAVA {
			sum += c.ExternalSurfaceArea(&grid)
		}
	}
	fmt.Println(sum)
}

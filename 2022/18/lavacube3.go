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
	return fmt.Sprintf("%d,%d,%d", p.x, p.y, p.z)
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
const LAVA = 10
const STEAM = 20

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
	if !found || !n.Is(LAVA) {
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

func (c *Cube) Is(what int) bool {
	return (c.is/10)*10 == (what/10)*10
}

func (c *Cube) SideTouching(g *Grid, axis, dir int, what int) bool {
	n, found := c.Neighbour(g, axis, dir)
	if !found {
		//fmt.Println("Missing neighbour for ", c.pos, axis, dir)
		return false
	}
	return n.Is(what)
}

func (c *Cube) ExpandIfTouching(g *Grid, what int) bool {
	for axis, sides := range sidemap {
		for dir := range sides {
			if c.SideTouching(g, axis, dir, what) {
				c.is = what
				return true
			}
		}
	}
	return false
}

func (c *Cube) SurfaceAreaTouching(g *Grid, what int) int {
	a := 0
	for axis, sides := range sidemap {
		for dir := range sides {
			if c.SideTouching(g, axis, dir, what) {
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
	// add 1 to avoid zeros in co-ords, since we used them for the border.
	p := Pos{Int(parts[0]) + 1, Int(parts[1]) + 1, Int(parts[2]) + 1}
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

const P_LAVA = 1
const P_STEAM = 2

func Print(g *Grid, mode int, min, xM, yM, zM int) {
	for z := min; z <= zM; z++ {
		PrintLayer(g, mode, min, xM, yM, z)
	}
}
func PrintLayer(g *Grid, mode int, min, xM, yM, z int) {
	fmt.Printf("  ")
	for y := 0; y <= yM; y++ {
		fmt.Printf("%d", y%9)
	}
	fmt.Println()
	for x := min; x <= xM; x++ {
		fmt.Printf("%d ", x%9)
		for y := min; y <= yM; y++ {
			c, present := (*g)[Pos{x, y, z}]
			if present && c.Is(LAVA) {
				if mode == P_LAVA {
					fmt.Printf("%d", c.SurfaceAreaTouching(g, STEAM))
				} else {
					fmt.Print("#")
				}
			} else if present && c.Is(STEAM) {
				if mode == P_STEAM {
					fmt.Printf("%d", c.is%STEAM)
				} else {
					fmt.Printf("*")
				}
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// ORIGINAL BUG WAS HERE - precisely 5 cubes in my input data had a co-ordinate of 0 and were being
// ovewritten here... I thought I had checked everything was 1-based co-ordinateds in the input data
// when choosing to use zero here, but obviously I missed the 5 "edge-cases", literally...
//
// Fixed above, by simply adding 1 to all incoming co-ords.
//
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

func ExpandSteam(g *Grid, xM, yM, zM int) {
	r := 1
	for {
		e := 0
		for z := 0; z <= zM; z++ {
			for x := 0; x <= xM; x++ {
				for y := 0; y <= yM; y++ {
					c := (*g)[Pos{x, y, z}]
					if !c.Is(AIR) {
						continue
					}
					if c.ExpandIfTouching(g, STEAM+r) {
						e++
					}
				}
			}
		}
		fmt.Printf("Expansion round %d: steam growth=%d\n", r, e)
		if e == 0 {
			break
		}
		//Print(g, P_STEAM, 0, xM, yM, zM)
		r++
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
	oC := len(grid)
	fmt.Printf("Loaded %d cubes into grid\n", len(grid))
	for _, c := range grid {
		c.CheckSides(&grid)
	}
	tSA := 0
	for _, c := range grid {
		if c.Is(LAVA) {
			//fmt.Println("L1: ", c.pos)
			tSA += c.SurfaceArea()
		} else {
			log.Fatal("not lava in initial grid!")
		}
	}
	fmt.Println("Total Surface Area: ", tSA)

	fmt.Println(xM, yM, zM)
	fmt.Println("Grid:")
	//Print(&grid, P_LAVA, 1, xM, yM, zM)
	fmt.Println()
	AddSteam(&grid, xM, yM, zM)

	check1 := 0
	nLava := 0
	for _, c := range grid {
		if c.Is(LAVA) {
			//fmt.Println("L2: ", c.pos)
			check1 += c.SurfaceArea()
			nLava++
		}
	}
	fmt.Println("Total Surface Area (after adding Steam): ", check1)
	if check1 != tSA {
		fmt.Println(oC, nLava)
		log.Fatal("Adding steam border changed lava surface area!")
	}

	fmt.Println("Expanding steam")
	ExpandSteam(&grid, xM+1, yM+1, zM+1)

	fmt.Println("Steamy Grid:")
	//Print(&grid, P_LAVA, 0, xM+1, yM+1, zM+1)

	fmt.Println("Calculating external surface area (via steam)...")
	eSA := 0
	for _, c := range grid {
		if c.Is(LAVA) {
			eSA += c.SurfaceAreaTouching(&grid, STEAM)
		}
	}
	fmt.Println("External Surface Area: ", eSA)

	fmt.Println("Calculating external surface area (via total - air)...")
	tSA2 := 0
	for _, c := range grid {
		if c.Is(LAVA) {
			//fmt.Println(c.pos, " is lava", c.SurfaceArea())
			tSA2 += c.SurfaceArea()
		}
	}
	fmt.Println("Total Surface Area (round 2)", tSA2, " vs ", tSA)
	tA := 0
	tL := 0
	tS := 0
	tC := 0
	for _, c := range grid {
		if c.Is(LAVA) {
			tC++
			tA += c.SurfaceAreaTouching(&grid, AIR)
			tL += c.SurfaceAreaTouching(&grid, LAVA)
			tS += c.SurfaceAreaTouching(&grid, STEAM)
		}
	}
	fmt.Println("Grid contains: ", tC, " blocks of lava, vs ", oC, " total sides=", tC*6)
	fmt.Println("Touching:")
	fmt.Println(" * AIR: ", tA)
	fmt.Println(" * STEAM: ", tS)
	fmt.Println(" * LAVA: ", tL)
	fmt.Println(" = ", tA+tS+tL, " vs ", tC*6)
	fmt.Println()
	fmt.Println("External Surface Area (via steam): ", tS)
	fmt.Println("External Surface Area (tSA2 - air): ", tSA2-tA)
	fmt.Println("External Surface Area (tSA - air): ", tSA-tA)

}

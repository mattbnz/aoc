// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 9, Puzzle 2.
// Rope Bridge, planck stuff, many knots!

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Knot struct {
	symbol string
	x, y   int

	minx, miny int
	maxx, maxy int

	visited map[string]bool
}

func (k Knot) String() string {
	return k.symbol
}

func (k Knot) Pos() string {
	return fmt.Sprintf("%d,%d", k.x, k.y)
}

func (k Knot) At(x, y int) bool {
	if k.x == x && k.y == y {
		return true
	}
	return false
}

func (k *Knot) Record() {
	if k.visited == nil {
		k.visited = map[string]bool{}
	}
	k.visited[k.Pos()] = true
}

func (k *Knot) Move(dir string) {
	if dir == "R" {
		k.x++
		k.maxx = max(k.maxx, k.x)
	} else if dir == "L" {
		k.x--
		k.minx = min(k.minx, k.x)
	} else if dir == "U" {
		k.y++
		k.maxy = max(k.maxy, k.y)
	} else if dir == "D" {
		k.y--
		k.miny = min(k.miny, k.y)
	} else {
		log.Fatalf("Bad Move: %s cannot move %s", k, dir)
	}
}

func abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

var DBG = true

func DPrintln(a ...any) {
	if DBG {
		fmt.Println(a...)
	}
}

func DPrintf(f string, args ...any) {
	if DBG {
		fmt.Printf(f, args...)
	}
}

func DPrint(a ...any) {
	if DBG {
		fmt.Print(a...)
	}
}

func (k *Knot) MoveX(dir int) {
	if dir > 0 {
		DPrintln("L")
		k.Move("L")
	} else {
		DPrintln("R")
		k.Move("R")
	}
}
func (k *Knot) MoveY(dir int) {
	if dir > 0 {
		DPrintln("D")
		k.Move("D")
	} else {
		DPrintln("U")
		k.Move("U")
	}
}

func (k *Knot) Chase(o Knot) {
	if k.x == o.x && k.y == o.y {
		return // together, nothing to do
	}
	DPrintf("%s => %s: ", k.Pos(), o.Pos())
	// Straight moves
	if k.x == o.x {
		d := k.y - o.y
		if abs(d) > 1 {
			k.MoveY(d)
		} else {
			DPrintln("ok")
		}
		return
	}
	if k.y == o.y {
		d := k.x - o.x
		if abs(d) > 1 {
			k.MoveX(d)
		} else {
			DPrintln("ok")
		}
		return
	}
	// Diagonal moves
	dX := k.x - o.x
	dY := k.y - o.y
	DPrint("dx=", dX, " dy=", dY)
	mX, mY := false, false
	if dX > 1 || dX < -1 {
		mX = true
	}
	if dY > 1 || dY < -1 {
		mY = true
	}
	DPrintln(" mX=", mX, " mY=", mY)
	if mX || mY {
		k.MoveX(dX)
		k.MoveY(dY)
	}
}

func max(v ...int) int {
	m := v[0]
	for _, n := range v[1:] {
		if n > m {
			m = n
		}
	}
	return m
}

func min(v ...int) int {
	m := v[0]
	for _, n := range v[1:] {
		if n < m {
			m = n
		}
	}
	return m
}

func extents(knots []*Knot) (int, int, int, int) {
	xMax := 5
	yMax := 4
	xMin := 0
	yMin := 0
	for _, k := range knots {
		if k.maxx > xMax {
			xMax = k.maxx
		}
		if k.minx < xMin {
			xMin = k.minx
		}
		if k.maxy > yMax {
			yMax = k.maxy
		}
		if k.miny < yMin {
			yMin = k.miny
		}
	}
	return xMin, xMax, yMin, yMax
}
func Print(knots []*Knot) {
	xMin, xMax, yMin, yMax := extents(knots)

	for y := yMax; y >= yMin; y-- {
	X:
		for x := xMin; x <= xMax; x++ {
			for _, k := range knots {
				if k.At(x, y) {
					fmt.Print(k)
					continue X
				}
			}
			if x == 0 && y == 0 {
				fmt.Print("s")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func main() {
	_, DBG = os.LookupEnv("DEBUG")

	s := bufio.NewScanner(os.Stdin)

	knots := []*Knot{{symbol: "H"}}
	for i := 1; i <= 9; i++ {
		knots = append(knots, &Knot{symbol: fmt.Sprintf("%d", i)})
	}

	fmt.Println("== Initial ==")
	Print(knots)

	// Read map and populate L/T heights
	for s.Scan() {
		dir, countS, found := strings.Cut(s.Text(), " ")
		if !found {
			log.Fatal("Bad instruction: ", s.Text())
		}
		count, err := strconv.Atoi(countS)
		if err != nil {
			log.Fatal("Bad count: ", s.Text())
		}
		fmt.Printf("== %s %d ==\n", dir, count)
		for i := 0; i < count; i++ {
			knots[0].Move(dir)
			for i, k := range knots {
				if i == 0 {
					continue
				}
				k.Chase(*knots[i-1])
				if i == 9 {
					k.Record()
				}
			}
		}
		//Print(knots)
	}

	xMin, xMax, yMin, yMax := extents(knots)
	count := 0
	for y := yMax; y >= yMin; y-- {
		for x := xMin; x <= xMax; x++ {
			pos := fmt.Sprintf("%d,%d", x, y)
			if x == 0 && y == 0 {
				fmt.Print("s")
				count++
			} else if _, visited := knots[9].visited[pos]; visited {
				fmt.Print("#")
				count++
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
	fmt.Println(count)
}

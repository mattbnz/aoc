// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 9, Puzzle 1.
// Rope Bridge, planck stuff.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type End struct {
	symbol string
	x, y   int

	minx, miny int
	maxx, maxy int

	visited map[string]bool
}

func (e End) String() string {
	return e.symbol
}

func (e End) Pos() string {
	return fmt.Sprintf("%d,%d", e.x, e.y)
}

func (e End) At(x, y int) bool {
	if e.x == x && e.y == y {
		return true
	}
	return false
}

func (e *End) Record() {
	if e.visited == nil {
		e.visited = map[string]bool{}
	}
	e.visited[e.Pos()] = true
}

func (e *End) Move(dir string) {
	if dir == "R" {
		e.x++
		e.maxx = max(e.maxx, e.x)
	} else if dir == "L" {
		e.x--
		e.minx = min(e.minx, e.x)
	} else if dir == "U" {
		e.y++
		e.maxy = max(e.maxy, e.y)
	} else if dir == "D" {
		e.y--
		e.miny = min(e.miny, e.y)
	} else {
		log.Fatalf("Bad Move: %s cannot move %s", e, dir)
	}
}

func abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

func (e *End) MoveX(dir int) {
	if dir > 0 {
		fmt.Println("L")
		e.Move("L")
	} else {
		fmt.Println("R")
		e.Move("R")
	}
}
func (e *End) MoveY(dir int) {
	if dir > 0 {
		fmt.Println("D")
		e.Move("D")
	} else {
		fmt.Println("U")
		e.Move("U")
	}
}

func (e *End) Chase(o End) {
	if e.x == o.x && e.y == o.y {
		return // together, nothing to do
	}
	fmt.Printf("%s => %s: ", e.Pos(), o.Pos())
	// Straight moves
	if e.x == o.x {
		d := e.y - o.y
		if abs(d) > 1 {
			e.MoveY(d)
		} else {
			fmt.Println("ok")
		}
		return
	}
	if e.y == o.y {
		d := e.x - o.x
		if abs(d) > 1 {
			e.MoveX(d)
		} else {
			fmt.Println("ok")
		}
		return
	}
	// Diagonal moves
	dX := e.x - o.x
	dY := e.y - o.y
	fmt.Print("dx=", dX, " dy=", dY)
	mX, mY := false, false
	if dX > 1 || dX < -1 {
		mX = true
	}
	if dY > 1 || dY < -1 {
		mY = true
	}
	fmt.Println(" mX=", mX, " mY=", mY)
	if mX || mY {
		e.MoveX(dX)
		e.MoveY(dY)
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

func Print(head, tail End) {
	xMax := max(head.maxx, tail.maxx, 5)
	yMax := max(head.maxy, tail.maxy, 4)
	xMin := min(head.minx, tail.minx, 0)
	yMin := min(head.miny, tail.miny, 0)

	for y := yMax; y >= yMin; y-- {
		for x := xMin; x <= xMax; x++ {
			if head.At(x, y) {
				fmt.Print(head)
			} else if tail.At(x, y) {
				fmt.Print(tail)
			} else if x == 0 && y == 0 {
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
	s := bufio.NewScanner(os.Stdin)

	head := End{symbol: "H"}
	tail := End{symbol: "T"}
	fmt.Println("== Initial ==")
	Print(head, tail)

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
			head.Move(dir)
			tail.Chase(head)
			tail.Record()
			//Print(head, tail)
		}
	}

	xMax := max(head.maxx, tail.maxx, 5)
	yMax := max(head.maxy, tail.maxy, 4)
	xMin := min(head.minx, tail.minx, 0)
	yMin := min(head.miny, tail.miny, 0)

	count := 0
	for y := yMax; y >= yMin; y-- {
		for x := xMin; x <= xMax; x++ {
			pos := fmt.Sprintf("%d,%d", x, y)
			if x == 0 && y == 0 {
				fmt.Print("s")
				count++
			} else if _, visited := tail.visited[pos]; visited {
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

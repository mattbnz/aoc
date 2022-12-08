// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 8, Puzzle 1.
//

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type Cell struct {
	height int
	// max height to dir before
	maxL, maxR int
	maxT, maxB int
}

func (c Cell) visible() bool {
	h := []int{c.maxL, c.maxR, c.maxT, c.maxB}
	sort.Ints(h)
	if h[0] < c.height {
		return true
	}
	return false
}

func ParseRow(row string, prev []*Cell) []*Cell {
	rv := []*Cell{}
	max := -1
	for i, c := range row {
		if c < 48 || c > 57 {
			log.Fatalf("%v is not a valid height in %s", c, row)
		}
		h := int(c - 48)
		tMax := -1
		if len(prev) > 0 {
			tMax = prev[i].maxT
			if prev[i].height > tMax {
				tMax = prev[i].height
			}
		}
		rv = append(rv, &Cell{height: h, maxL: max, maxT: tMax})
		if h > max {
			max = h
		}
	}
	return rv
}

func Print(heights [][]*Cell) {
	for _, row := range heights {
		for _, col := range row {
			fmt.Printf("%d (%v, %d, %d, %d, %d)  ", col.height, col.visible(), col.maxL, col.maxT, col.maxR, col.maxB)
		}
		fmt.Println("")
	}
}

func View(heights [][]int, xS, yS, xE, yE int) {
	xI, yI := 1, 1
	if xS > xE {
		xI = -1
	}
	if yS > yE {
		yI = -1
	}

	for x := xS; x != xE; x += xI {
		for y := yS; y != yE; y += yI {

		}
	}
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	heights := [][]*Cell{}
	prev := []*Cell{}

	// Read map and populate L/T heights
	for s.Scan() {
		row := ParseRow(s.Text(), prev)
		heights = append(heights, row)
		prev = row
	}
	cols := len(heights[0])
	// Now in reverse and populate R/B heights
	prev = []*Cell{}
	for y := len(heights) - 1; y >= 0; y-- {
		max := -1
		for x := cols - 1; x >= 0; x-- {
			bMax := -1
			if len(prev) > 0 {
				bMax = prev[x].maxB
				if prev[x].height > bMax {
					bMax = prev[x].height
				}
			}
			heights[y][x].maxB = bMax
			heights[y][x].maxR = max
			if heights[y][x].height > max {
				max = heights[y][x].height
			}
		}
		prev = heights[y]
	}
	//Print(heights)

	visible := 0
	for y := 0; y < len(heights); y++ {
		for x := 0; x < cols; x++ {
			if heights[y][x].visible() {
				visible++
			}
		}
	}
	fmt.Println(visible)

}

// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 8, Puzzle 2.
//

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Cell struct {
	height int
	// spaces clear in each dir
	clearL, clearR int
	clearT, clearB int
}

func ParseRow(row string, rowN int, blockT *[]map[int]int) []*Cell {
	rv := []*Cell{}
	blockL := map[int]int{}
	for i, c := range row {
		if c < 48 || c > 57 {
			log.Fatalf("%v is not a valid height in %s", c, row)
		}
		if len(*blockT) < i+1 {
			*blockT = append(*blockT, map[int]int{})
		}
		h := int(c - 48)
		rv = append(rv, &Cell{height: h, clearL: i - blockL[h], clearT: rowN - (*blockT)[i][h]})
		for b := h; b >= 0; b-- {
			blockL[b] = i
			(*blockT)[i][b] = rowN
		}
	}
	return rv
}

func Print(heights [][]*Cell) {
	for _, row := range heights {
		for _, col := range row {
			fmt.Printf("%d (%d, %d, %d, %d)  ", col.height, col.clearL, col.clearT, col.clearR, col.clearB)
		}
		fmt.Println("")
	}
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	heights := [][]*Cell{}
	blockT := []map[int]int{}

	// Read map and populate L/T heights
	rows := 0
	for s.Scan() {
		row := ParseRow(s.Text(), rows, &blockT)
		heights = append(heights, row)
		rows++
	}
	cols := len(heights[0])
	// Now in reverse and populate R/B heights
	blockB := []map[int]int{}
	for c := 0; c < cols; c++ {
		b := map[int]int{}
		for i := 0; i <= 9; i++ {
			b[i] = rows - 1
		}
		blockB = append(blockB, b)
	}
	for y := rows - 1; y >= 0; y-- {
		blockR := map[int]int{}
		for i := 0; i <= 9; i++ {
			blockR[i] = cols - 1
		}
		for x := cols - 1; x >= 0; x-- {
			h := heights[y][x].height
			heights[y][x].clearB = blockB[x][h] - y
			heights[y][x].clearR = blockR[h] - x
			for b := h; b >= 0; b-- {
				blockR[b] = x
				blockB[x][b] = y
			}
		}
	}
	//Print(heights)

	best := 0
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			c := heights[y][x]
			s := c.clearL * c.clearT * c.clearR * c.clearB
			if s > best {
				best = s
			}
		}
	}
	fmt.Println(best)

}

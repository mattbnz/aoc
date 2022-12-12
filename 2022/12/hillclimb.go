// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 12, Puzzle 1.
// Hill Climbing Algorithm.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func ParseRow(r string) ([]int, int, int) {
	row := []int{}
	start := -1
	end := -1
	for i, c := range r {
		if c == 'S' {
			start = i
			row = append(row, 0)
		} else if c == 'E' {
			end = i
			row = append(row, 25)
		} else if c >= 'a' && c <= 'z' {
			row = append(row, int(c-'a'))
		} else {
			log.Fatal("bad row: %s", r)
		}
	}
	return row, start, end
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	grid := [][]int{}
	startrow, startcol := -1, -1
	endrow, endcol := -1, -1

	row := 0
	for s.Scan() {
		r, s, e := ParseRow(s.Text())
		grid = append(grid, r)
		if s != -1 {
			startrow = row
			startcol = s
		}
		if e != -1 {
			endrow = row
			endcol = e
		}
		row++
	}

	for r, cols := range grid {
		for c, h := range cols {
			if r == startrow && c == startcol {
				fmt.Printf("S")
			} else if r == endrow && c == endcol {
				fmt.Printf("E")
			} else {
				fmt.Printf("%c", h+'a')
			}
		}
		fmt.Println("")
	}
}

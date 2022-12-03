// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 3, Puzzle 2.
// Rucksack mix-up.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func itemPriority(c rune) int {
	n := int(c)
	if n > int('Z') {
		return n - int('`')
	}
	return n - int('@') + 26
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	sum := 0
	member := 0
	items := map[rune]int{}
	for s.Scan() {
		for _, c := range s.Text() {
			items[c] |= 1 << member
		}
		member++
		if member == 3 {
			for c, count := range items {
				if count == 7 {
					sum += itemPriority(c)
					break
				}
			}
			items = map[rune]int{}
			member = 0
		}
	}

	fmt.Println(sum)
}

// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 3, Puzzle 1.
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
	for s.Scan() {
		items := map[rune]bool{}
		for i, c := range s.Text() {
			if i < len(s.Text())/2 {
				items[c] = true
			} else {
				if items[c] {
					sum += itemPriority(c)
					break
				}
			}
		}
	}

	fmt.Println(sum)
}

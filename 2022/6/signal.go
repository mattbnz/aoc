// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 6, Puzzle 1.
// Start of Frame Marker.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func allUnique(chars []rune) bool {
	t := map[rune]bool{}
	for _, c := range chars {
		if _, found := t[c]; found {
			return false
		}
		t[c] = true
	}
	return true
}

func main() {
	s := bufio.NewScanner(os.Stdin)

LINES:
	for s.Scan() {
		last4 := make([]rune, 4)
		idx := 0
		for i, c := range s.Text() {
			last4[idx] = c
			if i > 3 && allUnique(last4) {
				fmt.Println(i + 1)
				continue LINES
			}
			idx = (idx + 1) % 4
		}
		log.Fatal("No start sequence found in: ", s.Text())
	}
}

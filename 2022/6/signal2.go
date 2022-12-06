// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 6, Puzzle 2.
// Start of Message Marker.

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

const uniqLen = 14

func main() {
	s := bufio.NewScanner(os.Stdin)

LINES:
	for s.Scan() {
		last := make([]rune, uniqLen)
		idx := 0
		for i, c := range s.Text() {
			last[idx] = c
			if i > 3 && allUnique(last) {
				fmt.Println(i + 1)
				continue LINES
			}
			idx = (idx + 1) % uniqLen
		}
		log.Fatal("No start sequence found in: ", s.Text())
	}
}

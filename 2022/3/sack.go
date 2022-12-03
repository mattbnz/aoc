// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 3, Puzzle 1.
// Rucksack mix-up.

package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	Undefined int = iota
	Rock
	Paper
	Scissors
)

var Codes map[string]int = map[string]int{"A": Rock, "B": Paper, "C": Scissors, "X": Rock, "Y": Paper, "Z": Scissors}
var Names map[int]string = map[int]string{Rock: "Rock", Paper: "Paper", Scissors: "Scissors"}

const Loss = 0
const Draw = 3
const Win = 6

type BeatTable map[int]map[int]int

var Beats BeatTable = BeatTable{
	Rock:     {Rock: Draw, Paper: Win, Scissors: Loss},
	Paper:    {Rock: Loss, Paper: Draw, Scissors: Win},
	Scissors: {Rock: Win, Paper: Loss, Scissors: Draw},
}

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

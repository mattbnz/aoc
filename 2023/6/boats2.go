// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 6, Puzzle 1
// Wait For It

package main

import (
	"log"
)

type Race struct {
	Time         int
	BestDistance int
}

var Sample = Race{71530, 940200}
var Input = Race{48938595, 296192812361391}

func do(race Race) {
	beats := 0
	for c := 1; c < race.Time; c++ {
		dist := (race.Time - c) * c
		if dist > race.BestDistance {
			beats++
		}
	}
	log.Printf("Race has %d winning strategies", beats)
}

func main() {
	log.Printf("Sample")
	do(Sample)
	log.Printf("Input")
	do(Input)
}

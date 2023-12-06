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

var Sample = []Race{
	{7, 9}, {15, 40}, {30, 200},
}
var Input = []Race{
	{48, 296}, {93, 1928}, {85, 1236}, {95, 1391},
}

func do(races []Race) {
	product := 1

	for n, race := range races {
		beats := 0
		for c := 1; c < race.Time; c++ {
			dist := (race.Time - c) * c
			if dist > race.BestDistance {
				beats++
			}
		}
		log.Printf("Race %d: %d winning strategies", n, beats)
		product *= beats
	}

	log.Printf("%d", product)
}

func main() {
	log.Printf("Sample")
	do(Sample)
	log.Printf("Input")
	do(Input)
}

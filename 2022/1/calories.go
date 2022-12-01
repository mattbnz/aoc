// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 1, Puzzle 1
// Who has the most calories?

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	maxCal := 0

	// Read list from stdin
	s := bufio.NewScanner(os.Stdin)
	sum := 0
	for s.Scan() {
		if s.Text() == "" {
			if sum > maxCal {
				maxCal = sum
			}
			sum = 0
			continue
		}
		i, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Fatal("cannot take non-numeric input %s", s.Text())
		}
		sum += i
	}

	fmt.Println(maxCal)
}

// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 1, Puzzle 1
// Who has the most calories?

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {

	// Read list from stdin
	s := bufio.NewScanner(os.Stdin)
	sums := []int{}
	sum := 0
	for s.Scan() {
		if s.Text() == "" {
			sums = append(sums, sum)
			sum = 0
			continue
		}
		i, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Fatal("cannot take non-numeric input %s", s.Text())
		}
		sum += i
	}
	sums = append(sums, sum)
	sort.Ints(sums)

	sum = 0
	for i := len(sums) - 1; i > 0 && i > len(sums)-4; i-- {
		sum += sums[i]
	}
	fmt.Println(sum)
}

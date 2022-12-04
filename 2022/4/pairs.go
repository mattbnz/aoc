// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 4, Puzzle 1.
// Section overlaps.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func splitAssignment(a string) (min int, max int) {
	p := strings.Split(a, "-")
	min, err := strconv.Atoi(p[0])
	if err != nil {
		log.Fatal("bad input in pair", a)
	}
	max, err = strconv.Atoi(p[1])
	if err != nil {
		log.Fatal("bad input in pair", a)
	}
	return
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	count := 0
	for s.Scan() {
		a := strings.Split(s.Text(), ",")
		min1, max1 := splitAssignment(a[0])
		min2, max2 := splitAssignment(a[1])
		if (min1 <= min2 && max1 >= max2) || (min2 <= min1 && max2 >= max1) {
			count++
		}
	}

	fmt.Println(count)
}

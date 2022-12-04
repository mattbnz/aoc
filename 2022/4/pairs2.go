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
		sectors := map[int]int{}
		min1, max1 := splitAssignment(a[0])
		for i := min1; i <= max1; i++ {
			sectors[i]++
		}
		min2, max2 := splitAssignment(a[1])
		for i := min2; i <= max2; i++ {
			sectors[i]++
		}
		for _, c := range sectors {
			if c > 1 {
				count++
				break
			}
		}
	}

	fmt.Println(count)
}

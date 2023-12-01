// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 1, Puzzle 1
// Calibration Codes

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	sum := 0
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		var first, last string
		for _, c := range s.Text() {
			switch {
			case c >= '0' && c <= '9':
				if first == "" {
					first += string(c)
				} else {
					last = string(c)
				}
			}
		}
		if last == "" {
			last = first
		}
		num := first + last
		if len(num) != 2 {
			log.Fatalf("Did not find 2 digits (%s) from %s", num, s.Text())
		}
		i, err := strconv.Atoi(num)
		if err != nil {
			log.Fatalf("found num %s from %s was not a number", num, s.Text())
		}
		sum += i
	}

	fmt.Println(sum)
}

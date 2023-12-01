// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 1, Puzzle 2
// Calibration Codes

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var digits = map[rune]string{
	'1': "one",
	'2': "two",
	'3': "three",
	'4': "four",
	'5': "five",
	'6': "six",
	'7': "seven",
	'8': "eight",
	'9': "nine",
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	sum := 0
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		var first, last string
		l := s.Text()
	LINE:
		for n := 0; n < len(l); n++ {
			//fmt.Println(n, l, l[n], l[n:])
			var digit rune
			if l[n] >= '0' && l[n] <= '9' {
				digit = rune(l[n])
			} else {
			WORDS:
				for r, word := range digits {
					if idx := strings.Index(l[n:], word); idx == 0 {
						digit = r
						//fmt.Println(" ", l[n:], "==", word)
						// n += len(word) - 1
						break WORDS
					} else {
						//fmt.Println(" ", l[n:], "!=", word, idx)
					}
				}
			}
			if digit == 0 {
				//fmt.Println(" no digit found at ")
				continue LINE
			}
			if first == "" {
				first += string(digit)
			} else {
				last = string(digit)
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
		fmt.Println(l, "=>", i)
		sum += i
	}

	fmt.Println(sum)
}

// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 4, Puzzle 1
// Scratchcards

package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func numberList(list string) (rv []int) {
	for _, nStr := range strings.Split(strings.TrimSpace(list), " ") {
		nStr = strings.TrimSpace(nStr)
		if nStr == "" {
			continue
		}
		i, err := strconv.Atoi(nStr)
		if err != nil {
			log.Fatalf("Bad number '%s' (from %s): %v", nStr, list, err)
		}
		rv = append(rv, i)
	}
	return
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	sum := 0
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		wStr, nStr, found := strings.Cut(s.Text(), "|")
		if !found {
			log.Fatalf("Bad game line: %s", s.Text())
		}
		card, wStr2, found := strings.Cut(wStr, ":")
		if !found {
			log.Fatalf("Bad game line: %s", s.Text())
		}
		winners := numberList(wStr2)
		numbers := numberList(nStr)
		score := 0
	N:
		for _, n := range numbers {
			for _, wN := range winners {
				if n == wN {
					if score == 0 {
						score = 1
					} else {
						score *= 2
					}
					continue N
				}
			}
		}
		log.Printf("%s scores %d", card, score)
		sum += score
	}

	log.Printf("%d", sum)
}

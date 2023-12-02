// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 2, Puzzle 1
// Cube Conundrum

package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var maxCubes = map[string]int{
	"red":   12,
	"green": 13,
	"blue":  14,
}

var regCube = regexp.MustCompile(`\s+(\d+) (red|green|blue)\s*`)

func main() {
	s := bufio.NewScanner(os.Stdin)
	sum := 0
LINES:
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		game, draws, found := strings.Cut(s.Text(), ":")
		if !found {
			log.Fatalf("Bad game line: %s", s.Text())
		}
		id, err := strconv.Atoi(strings.TrimSpace(game[5:]))
		if err != nil {
			log.Fatalf("Bad game ID from (%s) for: %s", game[5:], s.Text())
		}
		for _, set := range strings.Split(draws, ";") {
			for _, nCube := range strings.Split(set, ",") {
				m := regCube.FindStringSubmatch(nCube)
				if m == nil {
					log.Fatalf("Bad cube input (%s) for line: %s", nCube, s.Text())
				}
				n, err := strconv.Atoi(m[1])
				if err != nil {
					log.Fatalf("Bad cube count from (%s) for line: %s", nCube, s.Text())
				}
				if n > maxCubes[m[2]] {
					log.Printf("Game %d is infeasible (%d %s cube > %d)", id, n, m[2], maxCubes[m[2]])
					continue LINES
				}
			}
		}
		log.Printf("Game %d  is feasible!", id)
		sum += id
	}

	log.Printf("%d", sum)
}

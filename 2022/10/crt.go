// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 10, Puzzle 1.
// Cathode-Ray Tube CPU.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Op int

func (o Op) String() string {
	if o == Undefined {
		return "UNDEFINED"
	} else if o == Noop {
		return "noop"
	} else if o == Load {
		return "load"
	} else if o == StoreAndLoad {
		return "storeandload"
	}
	return "ERROR"
}

const (
	Undefined Op = iota
	Noop
	Load
	StoreAndLoad
)

func main() {
	watches := map[int]int{
		20:  0,
		60:  0,
		100: 0,
		140: 0,
		180: 0,
		220: 0,
	}
	s := bufio.NewScanner(os.Stdin)

	tick := 1
	at := tick
	nextOp := Load
	inc := 0
	x := 1
	for {
		op := Undefined
		if tick == at {
			op = nextOp
		} else {
			op = Noop
		}
		if op == StoreAndLoad {
			x += inc
			op = Load
		}
		if _, watch := watches[tick]; watch {
			fmt.Printf("% 3d: WP x=%d\n", tick, x)
			watches[tick] = x
		}
		if op == Load {
			if !s.Scan() {
				fmt.Printf("% 3d: %s x=%d, EOF\n", tick, op, x)
				break
			}
			op, nS, found := strings.Cut(s.Text(), " ")
			if !found && op != "noop" {
				log.Fatal("Bad instruction: ", s.Text())
			}
			if op == "addx" {
				var err error
				inc, err = strconv.Atoi(nS)
				if err != nil {
					log.Fatal("Bad count: ", s.Text())
				}
				at = tick + 2
				nextOp = StoreAndLoad
			} else {
				nextOp = Load
				at = tick + 1
			}
		}
		fmt.Printf("% 3d: %s x=%d, %s(%d)@%d\n", tick, op, x, nextOp, inc, at)
		tick++
	}

	sum := 0
	for tick, watch := range watches {
		sum += tick * watch
	}
	fmt.Println(sum)
}

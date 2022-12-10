// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 10, Puzzle 2.
// Cathode-Ray Tube CPU with sprites.

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

type CRT [240]int

func (c *CRT) Print() {
	for i, v := range c {
		if i > 0 && i%40 == 0 {
			fmt.Println("")
		}
		if v == 1 {
			fmt.Printf("#")
		} else {
			fmt.Printf(".")
		}
	}
	fmt.Printf("\n\n")
}

func (c *CRT) Render(tick int, x int) {
	p := tick % 240
	rp := p % 40
	if x-1 <= rp && x+1 >= rp {
		c[p] = 1
	} else {
		c[p] = 0
	}
	//fmt.Printf("% 3d: p=%d rp=%d x=%d, v=%d\n", tick, p, rp, x, c[p])
	//c.Print()
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	pixels := CRT{}
	//pixels.Print()

	tick := 0
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
		pixels.Render(tick, x)
		if op == Load {
			if !s.Scan() {
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
		tick++
	}

	pixels.Print()
}

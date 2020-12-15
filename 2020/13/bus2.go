// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 13, Puzzle 2
// Bus schedule optimization.

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
)

func ProductExcludeBus(minutes []int, bus int) int {
    m := 1
    for _, b := range minutes {
        if b == -1 || b == bus {
            continue
        }
        m *= b
    }
    return m
}

func main() {
    r := bufio.NewReader(os.Stdin)
    _, err := r.ReadString('\n')
    if err != nil {
        log.Fatal("couldn't read departure!")
    }
    b, err := r.ReadString('\n')
    if err != nil {
        log.Fatal("couldn't read schedule! ")
    }
    minutes := []int{}
    for _, bus := range strings.Split(strings.TrimSpace(b), ",") {
        if bus == "x" {
            minutes = append(minutes, -1)
        } else {
            b, err := strconv.Atoi(bus)
            if err != nil {
                log.Fatal("bad bus", bus)
            }
            minutes = append(minutes, b)
        }
    }
    // stolen from https://github.com/lizthegrey/adventofcode/blob/main/2020/day13.go
    minValue := 0
	runningProduct := 1
    for min, bus := range minutes {
        if bus == -1 {
            continue
        }
		for (minValue+min)%bus != 0 {
			minValue += runningProduct
		}
		runningProduct *= bus
		fmt.Printf("%d + %d === 0 mod %d\n", minValue, min, bus)
		fmt.Printf("Sum so far: %d, product so far: %d\n", minValue, runningProduct)
	}
	fmt.Println(minValue)
}

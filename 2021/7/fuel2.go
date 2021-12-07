// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 7, Puzzle 2
// Crab Sub Fuel Optimisation part 2.

package main

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
    "strconv"
    "strings"
)

func main() {
    s := bufio.NewScanner(os.Stdin)
    s.Scan()
    subs := s.Text()

    pos := make(map[int]int)
    min, max := -1, -1
    for _, a := range strings.Split(subs, ",") {
        n, err := strconv.Atoi(a)
        if err != nil {
            log.Fatal("Bad sub position: ", a)
        }
        pos[n]++
        if min == -1 || n < min {
            min = n
        }
        if max == -1 || n > max {
            max = n
        }
    }

    used := make(map[int]int)
    for t := min; t<=max; t++ {
        sum := 0
        spos := t
        for p, c:= range pos {
            if p == spos{
                continue
            }
            mm := 0
            for i := 1; i<=int(math.Abs(float64(spos-p))); i++ {
                mm += i * c
            }
            sum += mm
        }
        used[spos] = sum
    }

    best := -1
    min = -1
    for pos, fuel := range used {
        if min == -1 || fuel < min {
            min  = fuel
            best = pos
        }
    }
    fmt.Printf("Move to %d using %d fuel\n", best, min)
}

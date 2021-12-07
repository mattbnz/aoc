// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 7, Puzzle 1
// Crab Sub Fuel Optimisation.

package main

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
    "sort"
    "strconv"
    "strings"
)

func main() {
    s := bufio.NewScanner(os.Stdin)
    s.Scan()
    subs := s.Text()

    pos := make(map[int]int)
    for _, a := range strings.Split(subs, ",") {
        n, err := strconv.Atoi(a)
        if err != nil {
            log.Fatal("Bad sub position: ", a)
        }
        pos[n]++
    }
    ordered := [][]int{}
    for pos, count  := range pos {
        ordered = append(ordered, []int{pos, count})
    }
    sort.Slice(ordered, func(a, b int) bool {
        return ordered[a][1] > ordered[b][1]
    })

    used := make(map[int]int)
    for _, t := range ordered {
        sum := 0
        spos := t[0]
        for _, pos := range ordered {
            if pos[0] == spos{
                continue
            }
            sum += int(math.Abs(float64(spos-pos[0])) * float64(pos[1]))
        }
        used[spos] = sum
    }

    best := -1
    min := -1
    for pos, fuel := range used {
        if min == -1 || fuel < min {
            min  = fuel
            best = pos
        }
    }
    fmt.Printf("Move to %d using %d fuel\n", best, min)
}

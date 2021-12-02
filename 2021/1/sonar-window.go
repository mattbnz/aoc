// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 1, Puzzle 2
// Count number of increases in depth using a 3-slot sliding window.

package main

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
    "strconv"
)

func sum(in []int) int {
    r := 0
    for _, v := range in {
        r += v
    }
    return r
}

func main() {
    var w []int
    inc := 0

    // Read list from stdin
    s := bufio.NewScanner(os.Stdin)
    last := math.MaxInt32
    for s.Scan() {
        i, err := strconv.Atoi(s.Text())
        if err != nil {
            log.Fatal("cannot take non-numeric input %s", s.Text())
        }
        w = append(w, i)
        if len(w) == 3 {
            t := sum(w)
            if t > last {
                inc++
            }
            last = t
            w = w[1:]
        }
    }

    // Walk through the list, look for the complement
    fmt.Println(inc)
}

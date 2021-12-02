// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 1, Puzzle 1
// Count number of increases in depth.

package main

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
    "strconv"
)

func main() {
     inc := 0

    // Read list from stdin
    s := bufio.NewScanner(os.Stdin)
    last := math.MaxInt32
    for s.Scan() {
        i, err := strconv.Atoi(s.Text())
        if err != nil {
            log.Fatal("cannot take non-numeric input %s", s.Text())
        }
        if i > last {
            inc++
        }
        last = i
    }

    // Walk through the list, look for the complement
    fmt.Println(inc)
}

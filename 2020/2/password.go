// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 2, Puzzle 1
// Check password validity.

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
)

func main() {
    valid := 0
    // Loop over passwords. 
    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
        fields := strings.Fields(s.Text())

        minmax := strings.Split(fields[0], "-")
        min, err := strconv.Atoi(minmax[0])
        if err != nil {
            log.Fatal("invalid input, min not numeric: %s", s.Text())
        }
        max, err := strconv.Atoi(minmax[1])
        if err != nil {
            log.Fatal("invalid input, max not numeric: %s", s.Text())
        }
        if min > max {
            log.Fatal("invalid input, min > max: %s", s.Text())
        }

        char := strings.TrimRight(fields[1], ":")
        if len(char) != 1 {
            log.Fatal("invalid input, didn't find letter: %s", s.Text())
        }

        count := strings.Count(fields[2], char)
        if count >= min && count <= max {
            valid++
        }
    }

    fmt.Println(valid);
}

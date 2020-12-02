// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 2, Puzzle 2
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

        ab := strings.Split(fields[0], "-")
        a, err := strconv.Atoi(ab[0])
        if err != nil {
            log.Fatal("invalid input, a not numeric: %s", s.Text())
        }
        b, err := strconv.Atoi(ab[1])
        if err != nil {
            log.Fatal("invalid input, b not numeric: %s", s.Text())
        }

        char := strings.TrimRight(fields[1], ":")
        if len(char) != 1 {
            log.Fatal("invalid input, didn't find letter: %s", s.Text())
        }
        pos_a := string(fields[2][a-1])
        pos_b := string(fields[2][b-1])

        if (pos_a != pos_b) && ((pos_a == char) || (pos_b == char)) {
            valid++
        }
    }

    fmt.Println(valid);
}

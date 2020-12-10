// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 9, Puzzle 1
// XMAS cipher

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
)

func main() {
    // Keep reading
    s := bufio.NewScanner(os.Stdin)
    buf := []int{}
    for s.Scan() {
        n, err := strconv.Atoi(s.Text())
        if err != nil {
            log.Fatal("Cannot parse: ", s.Text())
        }
        if len(buf) == 25 {
            // new input, find the sum.
            found := false
            a := 0
            for !found && a < len(buf) {
                b := 0
                for !found && b < len(buf) {
                    if buf[a] + buf[b] == n {
                        found = true
                    }
                    b++
                }
                a++
            }
            if (!found) {
                fmt.Println(n)
                os.Exit(0)
            }
            buf = append(buf[1:], n)
        } else {
            buf = append(buf, n)
        }
    }
    fmt.Println("failed")
}

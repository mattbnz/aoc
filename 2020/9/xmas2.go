// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 9, Puzzle 2
// XMAS cipher

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
)

var PREAMBLE = 25

func main() {
    // Keep reading
    s := bufio.NewScanner(os.Stdin)
    buf := []int{}
    l := 0
    key := -1
    for s.Scan() {
        n, err := strconv.Atoi(s.Text())
        if err != nil {
            log.Fatal("Cannot parse: ", s.Text())
        }
        if l >= PREAMBLE && key == -1{
            // new input, find the sum.
            found := false
            a := 0
            for !found && a < PREAMBLE {
                b := 0
                for !found && b < PREAMBLE {
                    base := l-PREAMBLE
                    if a != b && buf[base+a] + buf[base+b] == n {
                        found = true
                    }
                    b++
                }
                a++
            }
            if (!found) {
                key = n
            }
        }
        buf = append(buf, n)
        l++
    }
    if key == -1 {
        log.Fatal("didn't find a key!")
    }

    a := 0
    for a < len(buf) {
        b := a + 1
        sum := buf[a]
        min := buf[a]
        max := buf[a]
        for b < len(buf) {
            sum += buf[b]
            if buf[b] < min {
                min = buf[b]
            }
            if buf[b] > max {
                max = buf[b]
            }
            if sum == key {
                fmt.Println(min+max)
                os.Exit(0)
            }
            b++
        }
        a++
    }
    log.Fatal("failed")
}

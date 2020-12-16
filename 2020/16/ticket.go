// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 16, Puzzle 1
// Ticket Validity

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"
)

func main() {
    valid := make(map[int]bool)

    s := bufio.NewScanner(os.Stdin)
    // Read valid values for each field.
    for s.Scan() {
        if strings.TrimSpace(s.Text()) == "" {
            break
        }
        split := strings.Split(strings.TrimSpace(s.Text()), ":")
        pairs := strings.Split(split[1], " or ")
        for _, pair := range pairs {
            minmax := strings.Split(pair, "-")
            min, err := strconv.Atoi(strings.TrimSpace(minmax[0]))
            if err != nil {
                log.Fatal("Bad min for ", minmax[0], err)
            }
            max, err := strconv.Atoi(strings.TrimSpace(minmax[1]))
            if err != nil {
                log.Fatal("Bad max for ", minmax[1], err)
            }
            for i:=min; i<=max; i++ {
                valid[i] = true
            }
        }
    }

    // Skip 5 lines (your ticket/nearby ticket headers, blank lines)
    for i:=0; i<4; i++ {
        s.Scan()
    }

    // Now read and check nearby tickets
    sum := 0
    for s.Scan() {
        values := strings.Split(strings.TrimSpace(s.Text()), ",")
        for _, vstr := range values {
            v, err := strconv.Atoi(vstr)
            if err != nil {
                log.Fatal("Bad ticket value for ", s.Scan())
            }
            _, found := valid[v]
            if !found {
                sum += v
            }
        }
    }
    fmt.Println(sum)
}

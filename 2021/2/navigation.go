// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 2, Puzzle 1
// Follow navigation instructions

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
)

func sum(in []int) int {
    r := 0
    for _, v := range in {
        r += v
    }
    return r
}

func main() {
    depth := 0
    pos := 0

    // Read commands from stdin
    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
        fields := strings.Fields(s.Text())
        if len(fields) != 2 {
            log.Fatal("expected 2 fields in the command:  ", s.Text())
        }
        n, err := strconv.Atoi(fields[1])
        if err != nil {
            log.Fatal("command had non-numeric input: ", s.Text())
        }
        if fields[0] == "forward" {
            pos += n
        } else if fields[0] == "down" {
            depth += n
        } else if fields[0] == "up" {
            depth -= n
        } else {
            log.Fatal("unkonwn command: ", s.Text())
        }
    }

    // Walk through the list, look for the complement
    fmt.Println("Final position product: ", depth*pos)
}

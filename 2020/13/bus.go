// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 12, Puzzle 1
// Ship navigation.

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
    r := bufio.NewReader(os.Stdin)
    d, err := r.ReadString('\n')
    if err != nil {
        log.Fatal("couldn't read departure1")
    }
    b, err := r.ReadString('\n')
    if err != nil {
        log.Fatal("couldn't read schedule! ")
    }
    depart, err := strconv.Atoi(strings.TrimSpace(d))
    if err != nil {
        log.Fatal("departure time not int! ", d, err)
    }
    busses := []int{}
    for _, bus := range strings.Split(strings.TrimSpace(b), ",") {
        if bus == "x" {
            continue
        }
        b, err := strconv.Atoi(bus)
        if err != nil {
            log.Fatal("bad bus", bus)
        }
        busses = append(busses, b)
    }
    first := -1
    firstbus := -1
    for _, bus := range busses {
        next := ((depart/bus)+1)*bus
        if next < first || first == -1 {
            first = next
            firstbus = bus
        }
    }
    fmt.Println(firstbus * (first - depart))
}

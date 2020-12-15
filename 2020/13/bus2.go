// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 13, Puzzle 2
// Bus schedule optimization.

package main

import (
    "bufio"
    "fmt"
    "log"
    "math/big"
    "os"
    "strconv"
    "strings"
)

func ProductExcludeBus(minutes []int, bus int) int {
    m := 1
    for _, b := range minutes {
        if b == -1 || b == bus {
            continue
        }
        m *= b
    }
    return m
}

func main() {
    r := bufio.NewReader(os.Stdin)
    _, err := r.ReadString('\n')
    if err != nil {
        log.Fatal("couldn't read departure!")
    }
    b, err := r.ReadString('\n')
    if err != nil {
        log.Fatal("couldn't read schedule! ")
    }
    minutes := []int{}
    for _, bus := range strings.Split(strings.TrimSpace(b), ",") {
        if bus == "x" {
            minutes = append(minutes, -1)
        } else {
            b, err := strconv.Atoi(bus)
            if err != nil {
                log.Fatal("bad bus", bus)
            }
            minutes = append(minutes, b)
        }
    }
    // Chinese Remainder Theorum
    // per http://homepages.math.uic.edu/~leon/mcs425-s08/handouts/chinese_remainder.pdf
    m := 1
    z := make([]int, len(minutes))
    for o, bus := range minutes {
        if bus == -1 {
            continue
        }
        m *= bus
        z[o] = ProductExcludeBus(minutes, bus)
        fmt.Printf("z%d=%d\n", o, z[o])
    }
    fmt.Printf("m=%d\n", m)
    y := make([]int, len(minutes))
    for o, bus := range minutes {
        if bus == -1 {
            continue
        }
        b := big.NewInt(int64(z[o]))
        y[o] = int(b.ModInverse(b, big.NewInt(int64(bus))).Int64())
        fmt.Printf("y%d=%d\n", o, y[o])
    }
    w := make([]int, len(minutes))
    for o, bus := range minutes {
        if bus == -1 {
            continue
        }
        w[o] = y[o]*z[o] % m
        fmt.Printf("w%d=%d\n", o, w[o])
    }
    result := 0
    for o, bus := range minutes {
        if bus == -1 {
            continue
        }
        result += o * w[o]
    }
    fmt.Println(result)
    fmt.Println(result % m)
}

// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 6, Puzzle 2
// Exponential lanternfish.

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
    s := bufio.NewScanner(os.Stdin)
    s.Scan()
    fish := s.Text()

    ages := []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
    for _, a := range strings.Split(fish, ",") {
        n, err := strconv.Atoi(a)
        if err != nil {
            log.Fatal("Bad initial fish: ", a)
        }
        ages[n]++
    }
    fmt.Printf("Day %d: %v total %d\n",0, ages, sum(ages))

    for i:=0; i<256; i++ {
        next := []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
        next[6] = ages[0]
        next[8] = ages[0]
        for i:=0;  i<8; i++ {
            next[i] += ages[i+1]
        }
        ages = next
        fmt.Printf("Day %d: %v total %d\n",i+1, ages, sum(ages))
    }
    fmt.Printf("There are now %d fish\n", sum(ages))

}

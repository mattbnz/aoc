// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 6, Puzzle 1
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

func main() {
    s := bufio.NewScanner(os.Stdin)
    s.Scan()
    fish := s.Text()

    ages := []int{}
    for _, a := range strings.Split(fish, ",") {
        n, err := strconv.Atoi(a)
        if err != nil {
            log.Fatal("Bad initial fish: ", a)
        }
        ages = append(ages, n)
    }

    for i:=0; i<80; i++ {
        new := 0
        for n, a := range ages {
            if a == 0 {
                ages[n] = 6
                new++
            } else {
                ages[n]--
            }
        }
        for i:=0; i<new; i++ {
            ages = append(ages, 8)
        }
    }
    fmt.Printf("There are now %d fish\n", len(ages))

}

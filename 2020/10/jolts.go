// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 10, Puzzle 1
// Jolt differences.

package main

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
    "strconv"
)

func Find(list []int, jolt int) (int,int) {
    min := -1
    minj := -1
    for idx, j := range list {
        diff := int(math.Abs(float64(jolt - j)))
        if diff <= 3 {
            if minj==-1 || diff<minj {
                min = idx
                minj= diff
            }
        }
        //fmt.Printf("Find %d at %d: %d => %d\n", jolt, idx ,j, min)
    }
    return min, minj
}

func main() {
    s := bufio.NewScanner(os.Stdin)
    adaptors := []int{}
    for s.Scan() {
        n, err := strconv.Atoi(s.Text())
        if err != nil {
            log.Fatal("Cannot parse: ", s.Text())
        }
        adaptors = append(adaptors, n)
    }
    j1 := 0
    j3 := 0
    jolt := 0
    for len(adaptors) > 0 {
        c, diff := Find(adaptors, jolt)
        if c == -1 {
            log.Fatal("none found!")
        }
        if (diff == 1) {
            j1++
        } else if (diff == 3) {
            j3++
        }
        jolt = adaptors[c]
        fmt.Printf("Using %d (%dJ), diff of %d; j1=%d, j3=%d\n", c, adaptors[c], diff, j1, j3)
        adaptors = append(adaptors[:c], adaptors[c+1:]...)
    }
    fmt.Println(j1 * (j3+1))
}

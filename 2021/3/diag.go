// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 3, Puzzle 1
// Decode binary diagnostic output. 

package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
)

func main() {
    var ones []int
    lines := 0

    // Read commands from stdin
    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
        // Initialize ones first time through, now we know length of input.
        if len(ones) == 0 {
            for range s.Text() {
                ones = append(ones, 0)
            }
        }
        // Now iterate input and count ones.
        for n, b := range s.Text() {
            if b == '1' {
                ones[n]++
            }
        }
        lines++
    }

    gamma := 0.0
    epsilon := 0.0
    for n, c := range ones {
        v := math.Pow(2, float64(len(ones)-n-1))
        if c > lines/2 {
            // One most common, Zero least common 
            fmt.Println(n, " gamma", gamma, "+=", v)
            gamma += v
        } else {
            // Zero most common, One least common
            fmt.Println(n, " epsilon", epsilon, "+=", v)
            epsilon += v
        }
    }

    fmt.Println("Gamma: ", gamma)
    fmt.Println("Epsilon: ", epsilon)
    // Walk through the list, look for the complement
    fmt.Printf("Power consumption: %f\n", gamma*epsilon)
}

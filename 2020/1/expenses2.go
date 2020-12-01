// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 1, Puzzle 2
// Find three numbers in a list that sum to 2020 and return their product.

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
)

func main() {
    var expenses []int

    // Read list from stdin
    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
        i, err := strconv.Atoi(s.Text())
        if err != nil {
            log.Fatal("cannot take non-numeric input %s", s.Text())
        }
        expenses = append(expenses, i)
    }

    // Walk through the list and try all the possibilities.
    // This is slow, but the input is only len(200), so meh.
    a := 0
    for a < len(expenses) {
        b := 0
        for b < len(expenses) {
            c := 0
            for c < len(expenses) {
                sum := expenses[a] + expenses [b] + expenses[c]
                if sum == 2020 {
                    fmt.Println(expenses[a] * expenses[b] * expenses[c]);
                    os.Exit(0);
                }
                c++
            }
            b++
        }
        a++
    }
    log.Fatal("we failed :(")
}

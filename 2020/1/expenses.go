// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 1
// Find two numbers in a list that sum to 2020 and return their product.

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

    // Walk through the list, look for the complement
    for i, expense := range expenses {
        want := 2020 - expense
        for _, other := range expenses[i+1:] {
            if other == want {
                fmt.Println(expense * other);
                os.Exit(0);
            }
        }
    }
    log.Fatal("we failed :(")
}

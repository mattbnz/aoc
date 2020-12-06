// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 6, Puzzle 1
// Collate customs answers. 

package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)


func main() {
    s := bufio.NewScanner(os.Stdin)
    groups := [][]int{}
    questions := make([]int, 26)
    for s.Scan() {
        if strings.TrimSpace(s.Text()) == "" {
            groups = append(groups, questions)
            questions = make([]int, 26)
        }
        for _, char := range s.Text() {
            idx := char - 97
            questions[idx] = 1
        }
    }
    groups = append(groups, questions)

    sum := 0
    for _, questions := range groups {
        for _, q := range questions {
            sum += q
        }
    }

    fmt.Println(sum);
}

// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 6, Puzzle 2
// Collate customs answers for everyone. 

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
    questions := make([]int, 27)
    questions[26] = 0
    for s.Scan() {
        if strings.TrimSpace(s.Text()) == "" {
            groups = append(groups, questions)
            questions = make([]int, 27)
            questions[26] = 0
            continue
        }
        for _, char := range s.Text() {
            idx := char - 97
            questions[idx] += 1
        }
        questions[26]++
    }
    groups = append(groups, questions)

    sum := 0
    for _, questions := range groups {
        for _, q := range questions[:26] {
            if q == questions[26] {
                sum += 1
            }
        }
    }

    fmt.Println(sum);
}

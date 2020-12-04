// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 3, Puzzle 1
// Count trees in a grid.

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
)

func main() {
    treeMap := [][]int{}

    // Read map from stdin
    row := 0
    nCols := -1
    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
        l := []int{}
        for _, c := range s.Text() {
            if c == '.' {
                l = append(l, 0)
            } else if c == '#' {
                l = append(l, 1)
            } else {
                log.Fatal("unexpected map input: %s", s.Text())
            }
        }
        if nCols != -1 && nCols != len(l) {
            log.Fatal("map line not %d chars: %s", nCols, s.Text())
        } else {
            nCols = len(l)
        }
        treeMap = append(treeMap, l)
    }

    // Iterate and check.
    row = 1
    col := 3
    treesHit := 0
    for row < len(treeMap) {
        // Check for tree
        if treeMap[row][col % nCols] == 1 {
            treesHit++
        }
        // Move to new position
        col += 3
        row += 1
    }
    fmt.Println(treesHit)
}

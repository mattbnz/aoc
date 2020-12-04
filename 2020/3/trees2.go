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
    slopes := make([][]int, 2)
    slopes[0] = []int{1, 3, 5, 7, 1}
    slopes[1] = []int{1, 1, 1, 1, 2}

    // Read map from stdin
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

    product := 1

    // Iterate slopes.
    slope := 0
    for slope < len(slopes[0]) {
        // Iterate down the slope for trees.
        col := slopes[0][slope]
        row := slopes[1][slope]
        treesHit := 0
        for row < len(treeMap) {
            // Check for tree
            if treeMap[row][col % nCols] == 1 {
                treesHit++
            }
            //fmt.Println("Slope ", slope, ": at ", row, ", ", col, ", hit ", treesHit)
            // Move to new position
            col += slopes[0][slope]
            row += slopes[1][slope]
        }
        product *= treesHit
        slope++
    }
    fmt.Println(product)
}

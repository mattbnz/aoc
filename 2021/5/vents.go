// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 5, Puzzle 1
// Avoid the vents.

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
)

type Grid [][]int  // [[r1c1, r1c2, ...], [r2c1, r2c2, ...], ...]]
type Line [][]int  // [[x1,y1],[x2,y2]]

func ParsePoint(in string) []int {
    point := []int {0,0}
    p := strings.Split(in, ",")
    if len(p) != 2 {
        log.Fatal("Bad point: ", in)
    }
    for i, t := range p {
        n, err := strconv.Atoi(t)
        if err != nil {
            log.Fatal("Non numeric value in point: ", in)
        }
        point[i] = n
    }
    return point
}

func main() {
    s := bufio.NewScanner(os.Stdin)

    var lines []Line
    xmax, ymax := 0, 0

    // Load line defs, keep track of max dimensions
    for s.Scan() {
        t := strings.Fields(s.Text())
        if len(t) !=3 {
            log.Fatal("Bad iput line: ", t)
        }
        start := ParsePoint(t[0])
        end := ParsePoint(t[2])
        lines = append(lines, [][]int {start, end})
        if start[0] > xmax {
            xmax = start[0]
        }
        if end[0] > xmax {
            xmax = end[0]
        }
        if start[1] > ymax {
            ymax = start[1]
        }
        if end[1] > ymax {
            ymax = end[1]
        }
    }

    // Create grid of appropriate size
    var grid Grid
    fmt.Printf("Creating %dx%d grid\n", xmax+1, ymax+1)
    for i:=0; i<=xmax; i++ {
        row := make([]int, ymax+1)
        grid = append(grid, row)
    }
    //fmt.Printf("Grid: %v\n", grid)

    // Iterate through the lines, and increment cells
    for _, line := range lines {
        point := line[0]
        incIndex := 0
        if line[0][0] == line[1][0] {
            incIndex = 1
        } else if line[0][1] == line[1][1] {
            incIndex = 0
        } else {
            // Diagonal
            continue
        }
        inc := 1
        //fmt.Printf("Processing line: %v, incIndex=%d, inc=%d\n", line, incIndex, inc)
        if point[incIndex] >= line[1][incIndex] {
            inc = -1
        }
        for {
            grid[point[0]][point[1]]++
            //fmt.Printf("Incrementing %v now %d\n", point, grid[point[0]][point[1]])
            if point[incIndex] == line[1][incIndex] {
                break
            }
            point[incIndex] += inc
        }
        //fmt.Printf("New Grid: %v\n", grid)
    }

    points := 0
    for _, row := range grid {
        for _, cell := range row {
            if cell >=2 {
                points++
            }
        }
    }
    fmt.Printf("There are %d points with >2 lines passing through\n", points)

}

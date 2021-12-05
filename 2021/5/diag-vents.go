// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 5, Puzzle 2
// Avoid the diagonal vents too!

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
        //fmt.Printf("Processing %v\n", line)
        point := line[0]
        incX := 0
        incY := 0
        if line[0][0] > line[1][0] {
            incX = -1
        } else if line[0][0] < line[1][0] {
            incX = 1
        }
        if line[0][1] > line[1][1] {
           incY = -1
        } else if line[0][1] < line[1][1] {
           incY = 1
        }
        for {
            //fmt.Printf("Incrementing %v, incX=%d, incY=%d\n", point, incX, incY)
            grid[point[0]][point[1]]++
            if point[0] == line[1][0] && point[1] == line[1][1] {
                break
            }
            point[0] += incX
            point[1] += incY
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

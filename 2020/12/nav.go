// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 12, Puzzle 1
// Ship navigation.

package main

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
    "strconv"
)

func FixHeading(heading int) int {
    if heading < 0 {
        heading += 360
    }
    if heading > 360 {
        heading -= 360
    }
    if heading == 360 {
        heading = 0
    }
    return heading
}

func main() {
    s := bufio.NewScanner(os.Stdin)
    lat := 0
    long := 0
    heading := 90
    for s.Scan() {
        c := s.Text()[0]
        v, err := strconv.Atoi(s.Text()[1:])
        if err != nil {
            log.Fatal("Bad instruction ", s.Text())
        }
        if c == 'L' {
            heading = FixHeading(heading - v)
        } else if c == 'R'  {
            heading = FixHeading(heading + v)
        } else {
            dir := c
            if c == 'F' {
                if heading == 0 {
                    dir = 'N'
                }else if heading == 90 {
                    dir = 'E'
                } else if heading == 180 {
                    dir ='S'
                } else if heading == 270 {
                    dir = 'W'
                } else {
                    log.Fatal("heading is not cardinal! ", heading)
                }
            }
            if dir == 'N' {
                lat -= v
            } else if dir == 'S' {
                lat += v
            } else if dir == 'W' {
                long -= v
            } else if dir == 'E' {
                long += v
            }
        }
        fmt.Printf("%s: lat=%d, long=%d, heading=%d\n", s.Text(), lat, long, heading)
    }
    fmt.Println(math.Abs(float64(lat)) + math.Abs(float64(long)))
}

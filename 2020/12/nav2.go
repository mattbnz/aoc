// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 12, Puzzle 2
// Ship navigation with waypoint.

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
    off_long := 10
    off_lat := -1
    for s.Scan() {
        c := s.Text()[0]
        v, err := strconv.Atoi(s.Text()[1:])
        if err != nil {
            log.Fatal("Bad instruction ", s.Text())
        }
        if c == 'F' {
            lat += off_lat * v
            long += off_long * v
        } else if c == 'N' {
            off_lat -= v
        } else if c == 'S' {
            off_lat += v
        } else if c == 'W' {
            off_long -= v
        } else if c == 'E' {
            off_long += v
        } else if c == 'L' || c == 'R' {
            if (v == 180) {
                off_long *= -1
                off_lat *= -1
            } else if (c == 'L' && v == 90) || (c == 'R' && v == 270) {
                off_long, off_lat = off_lat, off_long*-1
            } else {
                off_long, off_lat = off_lat*-1, off_long
            }
            heading = FixHeading(heading + v)
        }
        fmt.Printf("%s: lat=%d, long=%d, heading=%d\n", s.Text(), lat, long, heading)
    }
    fmt.Println(math.Abs(float64(lat)) + math.Abs(float64(long)))
}


//     0 90 180 270
// C:  N E  S   W
// AC: N W  S   E

//     R90 L90 R270 L270
// N   E   W   W    E 
// E   S   N   N    S 
// S   W   E   E    W
// W   N   S   S    N


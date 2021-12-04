// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 3, Puzzle 2
// Decode binary diagnostic output. 

package main

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
)

// Count bits in pos for each line, and return most or least common rune
func CountBits(lines []string, pos int, most bool) rune {
    ones := 0
    for _, l := range lines {
        if l[pos] == '1' {
            ones++
        }
    }
    fmt.Printf("%d %t %d\n", pos, most, ones)
    zeros := len(lines) - ones
    if ones == zeros {
        // equal numbers
        if most { return '1'} else { return '0' }
    } else if ones > zeros {
        // 1 is most common
        if most {
            return '1'
        } else {
            return '0'
        }
    } else {
        // 0 is most common
        if most {
            return '0'
        } else {
            return '1'
        }
    }
}

// Return a subset of line with a matching rune in the specified pos
func RemoveUnmatched(lines []string, pos int, match rune) []string {
    var r []string
    for _, l := range lines {
        if []rune(l)[pos] == match {
            r = append(r, l)
            fmt.Printf("%d (%c) #%d: keeping line %s\n", pos, match, len(lines), l)
        }
    }
    return r
}

// Convert a string of binary bits to a decimal
func ToDecimal(bits string) int {
    r := 0.0
    for n, b := range bits {
        if b == '1' {
            r += math.Pow(2, float64(len(bits)-n-1))
        }
    }
    fmt.Printf("%s is %d\n", bits, int(r))
    return int(r)
}


func Filter(lines []string, most bool) int {
    blen := len(lines[0])
    for b := 0; b<blen; b++ {
        mcb := CountBits(lines, b, most)
        lines = RemoveUnmatched(lines, b, mcb)
        fmt.Printf("%d: mcb=%c, now %d lines\n", b, mcb, len(lines))
        if len(lines) == 1 {
            return ToDecimal(lines[0])
        }
    }
    log.Fatal("failed to get to single line in filter with most=", most)
    return -1
}

func main() {
    var lines []string

    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
        lines = append(lines, s.Text())
    }

    o := Filter(lines, true)
    co2 := Filter(lines, false)


    fmt.Println("Oxygen Generator Rating: ", o)
    fmt.Println("CO2 Scrubber Rating: ", co2)
    fmt.Printf("Life Support Rating: %d\n", o*co2)
}

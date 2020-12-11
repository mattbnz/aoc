// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 11, Puzzle 1
// Seat layout.

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
)

// Return 1 if the seat is occupied, 0 otherwise (including invalid)`
func SeatOccupied(rows [][]int, row, seat int) int {
    if row < 0 || row >= len(rows) {
        return 0
    }
    if seat <0 || seat >= len(rows[0]) {
        return 0
    }
    if rows[row][seat] == 1 {
        return 1
    }
    return 0
}

func AllocateSeats(rows [][]int) (int, [][]int) {
    newrows := make([][]int, len(rows))
    changed := 0
    for r, seats := range rows {
        newrows[r] = make([]int, len(seats))
        pr := r-1
        nr := r+1
        for s, state := range seats {
            if state == -1 {
                newrows[r][s] = -1
                continue
            }
            occupied := 0
            ps := s-1
            ns := s+1
            // Previous Row
            occupied += SeatOccupied(rows, pr, ps)
            occupied += SeatOccupied(rows, pr, s)
            occupied += SeatOccupied(rows, pr, ns)
            // This Row
            occupied += SeatOccupied(rows, r, ps)
            occupied += SeatOccupied(rows, r, ns)
            // Next Row
            occupied += SeatOccupied(rows, nr, ps)
            occupied += SeatOccupied(rows, nr, s)
            occupied += SeatOccupied(rows, nr, ns)
            // Then
            newrows[r][s] = rows[r][s]
            if state == 0 {
                if occupied == 0 {
                    newrows[r][s] = 1
                    changed++
                }
            } else if state == 1 {
                if occupied >= 4 {
                    newrows[r][s] = 0
                    changed++
                }
            } else {
                log.Fatal("bad state in rows!", r, s, rows)
            }
        }
    }
    return changed, newrows
}

func main() {
    s := bufio.NewScanner(os.Stdin)
    rows := [][]int{}  // -1=Floor; 0=Empty; 1=Occupied
    for s.Scan() {
        seats := []int{}
        for _, r := range s.Text() {
            if r == '.' {
                seats = append(seats, -1)
            } else if r == 'L' {
                seats = append(seats, 0)
            } else if r == '#' {
                seats = append(seats, 1)
            } else {
                log.Fatal("Bad seat state: ", s.Text())
            }
        }
        rows = append(rows, seats)
    }
    // Run until stable.
    changed := -1
    for changed != 0 {
        changed, rows = AllocateSeats(rows)
    }
    // Count
    occupied := 0
    for _, seats := range rows {
        for _, state := range seats {
            if state == 1 {
                occupied++
            }
        }
    }
    fmt.Println(occupied)
}

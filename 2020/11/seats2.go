// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 11, Puzzle 2
// Seat layout.

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
)

// Return 1 if the seat is occupied, 0 otherwise (including invalid)`
func SeatOccupied(rows [][]int, row, seat, rd, sd int) int {
    for {
        row += rd
        seat += sd
        if row < 0 || row >= len(rows) {
            return 0
        }
        if seat <0 || seat >= len(rows[0]) {
            return 0
        }
        if rows[row][seat] == -1 {
            continue
        }
        if rows[row][seat] == 1 {
            return 1
        } else {
            return 0
        }
    }
    log.Fatal("Seat occupied failed!", rows, row, seat, rd, sd)
    return 0
}

func AllocateSeats(rows [][]int) (int, [][]int) {
    newrows := make([][]int, len(rows))
    changed := 0
    for r, seats := range rows {
        newrows[r] = make([]int, len(seats))
        for s, state := range seats {
            if state == -1 {
                newrows[r][s] = -1
                continue
            }
            occupied := 0
            occupied += SeatOccupied(rows, r, s, 0, -1)  // Left
            occupied += SeatOccupied(rows, r, s, 0, 1)   // Right
            occupied += SeatOccupied(rows, r, s, -1, 0)  // Up
            occupied += SeatOccupied(rows, r, s, 1, 0)   // Down
            occupied += SeatOccupied(rows, r, s, -1, -1) // Diag Up Left
            occupied += SeatOccupied(rows, r, s, -1, 1)  // Diag Up Right
            occupied += SeatOccupied(rows, r, s, 1, -1)  // Diag Down Left
            occupied += SeatOccupied(rows, r, s, 1, 1)   // Diag Down Right
            // Then
            newrows[r][s] = rows[r][s]
            if state == 0 {
                if occupied == 0 {
                    newrows[r][s] = 1
                    changed++
                }
            } else if state == 1 {
                if occupied >= 5 {
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

// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 4, Puzzle 1
// Bingo!

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
)

type Cell struct {
    val string
    hit bool
}
type Board [][]Cell

func LoadBoards(s *bufio.Scanner) []Board {
    var boards []Board

    b := Board{}
    for s.Scan() {
        if strings.TrimSpace(s.Text()) == "" {
            boards = append(boards, b)
            b = Board{}
            continue
        }
        r := []Cell{}
        for _, n := range strings.Fields(s.Text()) {
            r = append(r, Cell{n, false})
        }
        b = append(b, r)
    }
    if len(b) > 0 {
        boards = append(boards, b)
    }
    return boards
}

func UpdateBoard(b Board, n string) Board {
    for r, cells := range b {
        for c, cell := range cells {
            if cell.val == n  {
                b[r][c].hit = true
            }
        }
    }
    return b
}

func CheckBoard(b Board) bool {
    cols := make([]int, len(b))
    for i := range cols { cols[i] = 0 }

    for _, cells := range b {
        hits := 0
        for c, cell := range cells {
            if cell.hit {
                //fmt.Printf("  ! on %s @ (%d, %d)\n", cell.val, r, c)
                hits++
                cols[c]++
            }
        }
        // return immediately if a row is all hit.
        if hits == len(cells) {
            return true
        }
    }

    // now check if any columns are all hit.
    for _, n := range cols {
        if n == len(b) {
            return true
        }
    }
    return false;
}

func SumBoard(b Board) int {
    r := 0
    for _, cells := range b {
        for _, cell := range cells {
            if cell.hit {
                continue
            }
            v, err := strconv.Atoi(cell.val)
            if err != nil {
                log.Fatal("non-numeric value in board! ", cell.val)
            }
            r += v
        }
    }
    return r
}

func main() {
    s := bufio.NewScanner(os.Stdin)
    s.Scan()
    draw := s.Text()
    s.Scan() // consume blank line after draw.

    // Load boards
    boards := LoadBoards(s)
    fmt.Println("Num Boards: ", len(boards))
    //fmt.Printf("%v\n", boards)
    fmt.Println("Draw: ", draw)

    // Run through the draw
    for _, n := range strings.Split(draw, ",") {
        //fmt.Println("Playing ", n)
        for i, b := range boards {
            boards[i] = UpdateBoard(b, n)
            //fmt.Println(" + Board ", i)
            if CheckBoard(boards[i]) {
                v, err := strconv.Atoi(n)
                if err != nil {
                    log.Fatal("winning number is not a number! ", n)
                }
                fmt.Println("Board ", i, " wins with a score of ", SumBoard(boards[i])*v)
                os.Exit(0)
            }
        }
    }
    log.Fatal("No Winner!")
}

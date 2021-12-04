// Copyright (C) 2021 Matt Brown

// Advent of Code 2021 - Day 4, Puzzle 2
// Losing Bingo!

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
type Squares [][]Cell
type Board struct {
    num int
    board Squares
}

func LoadBoards(s *bufio.Scanner) []Board {
    var boards []Board

    c := 1
    b := Board{c, Squares{}}
    for s.Scan() {
        if strings.TrimSpace(s.Text()) == "" {
            boards = append(boards, b)
            c++
            b = Board{c, Squares{}}
            continue
        }
        r := []Cell{}
        for _, n := range strings.Fields(s.Text()) {
            r = append(r, Cell{n, false})
        }
        b.board = append(b.board, r)
    }
    if len(b.board) > 0 {
        boards = append(boards, b)
    }
    return boards
}

func UpdateSquares(b Squares, n string) Squares {
    for r, cells := range b {
        for c, cell := range cells {
            if cell.val == n  {
                b[r][c].hit = true
            }
        }
    }
    return b
}

func CheckSquares(b Squares) bool {
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

func SumSquares(b Squares) int {
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
        can_win := false
        if len(boards) == 1 {
            can_win = true
        }
        new_boards := []Board{}
        for _, b := range boards {
            new_b := UpdateSquares(b.board, n)
            //fmt.Println(" + Squares ", i)
            if !CheckSquares(new_b) {
                b.board = new_b
                new_boards = append(new_boards, b)
            } else if can_win {
                v, err := strconv.Atoi(n)
                if err != nil {
                    log.Fatal("winning number is not a number! ", n)
                }
                fmt.Println("Board ", b.num, " wins last with a score of ", SumSquares(new_b)*v)
                os.Exit(0)
            }
        }
        boards = new_boards
    }
    log.Fatal("No Winner!")
}

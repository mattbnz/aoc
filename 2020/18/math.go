// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 18, Puzzle 1
// Math operator precedence

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"
)

func Atoi(in string) int {
    v, err := strconv.Atoi(in)
    if err != nil {
        log.Fatal("Couldn't convert ", in)
    }
    return v
}

func Evaluate(tokens []string) (int, int) {
    rv := -1
    op := ""
    skip := 0
    done := 0
    for n, tok := range tokens {
        if skip > 0 {
            skip--
            done++
            continue
        }
        end := false
        if rv == -1 || op != "" {
            // Get a value
            v := 0
            if tok == "(" {
                skip, v = Evaluate(tokens[n+1:])
            } else {
                v = Atoi(tok)
            }
            if rv == -1 {
                rv = v  // First iteration
            } else if op == "+" {
                rv += v
            } else if op == "*" {
                rv *= v
            }
            op = ""
        } else if tok == ")" {
            end = true
        } else if op == "" {
            op = tok
        }
        done++
        //fmt.Printf("%d: %s; rv=%d, skip=%d, done=%d\n", n, tok, rv, skip, done)
        if end {
            break
        }
    }
    return done, rv
}

func main() {

    s := bufio.NewScanner(os.Stdin)

    // Evaluate each line and add to the sum
    sum := 0
    for s.Scan() {
        t := strings.Replace(s.Text(), "(", "( ", -1)
        t = strings.Replace(t, ")", " )", -1)
        tokens := strings.Split(t, " ")
        _, v := Evaluate(tokens)
        //fmt.Printf("%s = %d\n", t, v)
        sum += v
    }
    fmt.Println(sum)
}

// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 19, Puzzle 1
// Monster Message matching.

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"
)

type Rule struct {
    rules [][]int
    match string
}

func Atoi(in string) int {
    v, err := strconv.Atoi(in)
    if err != nil {
        log.Fatal("Couldn't convert ", in)
    }
    return v
}

func MatchRule(rules map[int]Rule, rule int, message string, at int) (int, bool) {
    r := rules[rule]
    // Simple case of string match first.
    if len(r.rules) == 0 {
        if message[at:at+len(r.match)] != r.match {
            return -1, false
        } else {
            return len(r.match), true
        }
    }
    // Then deal with rule references
    for _, option := range r.rules {
        pos := at
        good := true
        for _, rn := range option {
            consumed, matched := MatchRule(rules, rn, message, pos)
            if !matched {
                good = false
                break
            } else {
                pos += consumed
            }
        }
        if good {
            return pos-at, true
        }
    }
    // No match found
    return -1, false
}

func main() {
    s := bufio.NewScanner(os.Stdin)

    // Read in rules
    rules := make(map[int]Rule)
    for s.Scan() {
        if s.Text() == "" {
            break
        }
        t := strings.Split(s.Text(), ": ")
        n, err := strconv.Atoi(t[0])
        if err != nil {
            log.Fatal("Bad rule number!")
        }
        r := Rule{}
        if t[1][0] == '"' {
            r.match = strings.Trim(t[1], "\"")
        } else {
            for _, m := range strings.Split(t[1], " | ") {
                l := []int{}
                for _, p := range strings.Split(m, " ") {
                    l = append(l, Atoi(p))
                }
                r.rules = append(r.rules, l)
            }
        }
        rules[n] = r
    }

    // Match Strings to rules
    sum := 0
    for s.Scan() {
        used, matched := MatchRule(rules, 0, s.Text(), 0)
        if matched && used==len(s.Text()) {
            sum++
        }
    }

    fmt.Println(sum)
}

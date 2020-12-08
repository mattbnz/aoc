// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 8, Puzzle 1
// Solve the infinite loop...

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "regexp"
    "strconv"
    "strings"
)

type BagMap map[string][]ContainedBag
type ContainedBag struct {
    color string
    num int
}
var BagRe = regexp.MustCompile(`(\d)+\s(.*)`)
var StripRe = regexp.MustCompile(` +bags?.?$`)

func StripName(input string) string {
    return StripRe.ReplaceAllString(input, "")
}

func FindBagsFor(want string, rules BagMap) []string {
    rv := []string{}
    for bag, contains := range rules {
        for _, inner := range contains {
            if inner.color == want {
                rv = append(rv, bag)
                break
            }
        }
    }
    return rv
}

func main() {
    // Read in the instructions
    code := []string{}
    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
        code = append(code, s.Text())
    }

    // Execute 
    visited := make([]bool, len(code))
    ip := 0
    acc := 0
    for {
        if visited[ip] {
            break
        }
        visited[ip] = true
        parts := strings.Fields(code[ip])
        v, err := strconv.Atoi(strings.TrimLeft(parts[1], "+"))
        if err != nil {
            log.Fatal("Invalid instruction at ", ip)
        }
        if (parts[0] == "acc") {
            acc += v
        } else if (parts[0] == "jmp") {
            ip += v
            continue
        }
        ip++
    }
    fmt.Println(acc)
}

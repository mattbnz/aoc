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
        log.Fatal(fmt.Sprintf("Couldn't convert %s", in))
    }
    return v
}

func RemoveEmpty(in []string) []string {
    rv := []string{}
    for _, t := range in {
        g := strings.TrimSpace(t)
        if g != "" {
            rv = append(rv, g)
        }
    }
    return rv
}

func Evaluate(tokens []string) (int, int) {
    //fmt.Printf("in: %v\n", tokens)
    // Collapse brackets
    working := make([]string, len(tokens))
    copy(working, tokens)
    while := true
    done := 0
    endat := -1
    for while {
        while = false
        for i:=0; i<len(working); i++ {
            if working[i] == ")" {
                endat = i
                done++
                break
            }
            if working[i] != "(" {
                continue
            }
            skip, v := Evaluate(working[i+1:])
            //fmt.Printf(" b%d: %v %d %d %d\n", i, working, len(working), skip, v)
            if i == 0 {
                working[i] = fmt.Sprintf("%d", v)
                working = append(working[:1], working[i+1+skip:]...)
            } else if i+skip >= len(working)-1 {
                working = append(working[:i], fmt.Sprintf("%d", v))
            } else {
                working = append(working[:i], working[i+skip:]...)
                working[i] = fmt.Sprintf("%d", v)
            }
            //fmt.Printf(" a%d: %v %d %d\n", i, working, len(working), skip)
            done += skip
            //fmt.Printf(" !: %v (consumed %d)\n", working, done)
            while = true
            break
        }
    }
    if endat == -1 {
        endat = len(working)
    }
    //fmt.Printf(" (: %v (consumed %d)\n", working[:endat], done)
    // Addition
    while = true
    for while {
        while = false
        for i:=1; i<endat; i++ {
            if working[i] != "+" {
                continue
            }
            working[i+1] = fmt.Sprintf("%d", Atoi(working[i-1]) + Atoi(working[i+1]))
            working = append(working[:i-1], working[i+1:]...)
            done += 2
            endat -= 2
            while = true
            break
        }
    }
    //fmt.Printf(" +: %v (consumed %d)\n", working[:endat], done)
    // Multiplication
    rv := Atoi(working[0])
    done++
    for i := 1; i<endat; i++ {
        done++
        if working[i] == "*" {
            continue
        }
        rv *= Atoi(working[i])
    }
    //fmt.Printf("*: %v (consumed %d)\n", working[:endat], done)
    return done, rv
}

func main() {

    s := bufio.NewScanner(os.Stdin)

    // Evaluate each line and add to the sum
    sum := 0
    for s.Scan() {
        t := strings.Replace(s.Text(), "(", "( ", -1)
        t = strings.Replace(t, ")", " )", -1)
        tokens := RemoveEmpty(strings.Split(t, " "))
        _, v := Evaluate(tokens)
        fmt.Printf("%s = %d\n", t, v)
        sum += v
    }
    fmt.Println(sum)
}

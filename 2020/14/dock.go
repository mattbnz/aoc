// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 14, Puzzle 1
// Docking bitmasks.

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "regexp"
    "strconv"
)
var MemRe = regexp.MustCompile(`^mem\[(.*)\] = (\d+)`)

func ApplyMask(mask string, v int) int {
    var bpos uint
    bpos = 0
    rv := 0
    for i := len(mask)-1; i>= 0; i-- {
        bit := 1 << bpos
        if mask[i] == 'X' {
            rv |= (v & bit)
        } else if mask[i] == '1' {
            rv |= bit
        }
        bpos++
    }
    return rv
}

func main() {
    s := bufio.NewScanner(os.Stdin)
    mem := make(map[int]int)
    mask := ""
    for s.Scan() {
        if s.Text()[:4] == "mask" {
            mask = s.Text()[7:]
            continue
        }
        match := MemRe.FindStringSubmatch(s.Text())
        r, err := strconv.Atoi(match[1])
        if err != nil {
            log.Fatal("Bad instruction ", s.Text())
        }
        v, err := strconv.Atoi(match[2])
        if err != nil {
            log.Fatal("Bad instruction ", s.Text())
        }
        //fmt.Printf("%036b\n%s\n", v, mask)
        mem[r] = ApplyMask(mask, v)
        //fmt.Printf("%036b\n\n", mem[r])
    }
    sum := 0
    for i, v := range mem {
        //fmt.Printf("%d: %d\n", i, v)
        sum += v
    }
    fmt.Println(sum)
}

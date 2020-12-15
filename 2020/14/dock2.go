// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 14, Puzzle 2
// Docking bitmasks with floating values.

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

func GenerateAddresses(a int, mask string, mpos int, addr int) []int {
    //fmt.Printf("%d:\n%036b\n%s\n%036b\n", mpos, a, mask, addr)
    if mpos == -1 {
        return []int{addr}
    }
    rv := []int{}
    for i := mpos; i>= 0; i-- {
        var bpos uint
        bpos = uint(len(mask)-1-i)
        bit := 1 << bpos
        if mask[i] == 'X' {
            addr |= bit
            rv = append(rv, GenerateAddresses(a, mask, i-1, addr)...)
            addr ^= bit
            rv = append(rv, GenerateAddresses(a, mask, i-1, addr)...)
            return rv
        } else if mask[i] == '1' {
            addr |= bit
        } else if mask[i] == '0' {
            addr |= (a & bit)
        }
    }
    rv = append(rv, addr)
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
        for _, addr := range GenerateAddresses(r, mask, len(mask)-1, 0) {
            //fmt.Printf("%036b: %d\n", addr, v)
            mem[addr] = v
        }
        //fmt.Printf("%036b\n%s\n", v, mask)
        //fmt.Printf("%036b\n\n", mem[r])
    }
    sum := 0
    for _, v := range mem {
        //fmt.Printf("%d: %d\n", i, v)
        sum += v
    }
    fmt.Println(sum)
}

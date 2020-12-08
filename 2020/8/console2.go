// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 8, Puzzle 2
// Solve the infinite loop...

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
)

func main() {
    // Read in the instructions
    code := []string{}
    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
        code = append(code, s.Text())
    }

    // Build a list of jumps and nops
    jumps := []int{}
    nops := []int{}
    for idx, instr := range code {
        if instr[:3] == "jmp" {
            jumps = append(jumps, idx)
        } else if instr[:3] == "nop" {
            nops = append(nops, idx)
        }
    }

    // Execute 
    skipjump := -1
    jmpnop := -1
    acc := 0
    for {
        visited := make([]bool, len(code))
        ip := 0
        acc = 0
        ok := true
        fmt.Printf("trying with skipjump=%d and jmpnop=%d\n", skipjump, jmpnop)
        for ip < len(code) {
            fmt.Printf(">%d (%s): acc=%d\n", ip, code[ip], acc)
            if visited[ip] {
                if skipjump < len(jumps) {
                    skipjump++
                } else if jmpnop < len(jumps) {
                    jmpnop++
                } else {
                    log.Fatal("run out of options :(")
                }
                ok = false
                break
            }
            visited[ip] = true
            parts := strings.Fields(code[ip])
            v, err := strconv.Atoi(strings.TrimLeft(parts[1], "+"))
            if err != nil {
                log.Fatal("Invalid instruction at ", ip)
            }
            if skipjump != -1 && ip == jumps[skipjump] {
                parts[0] = "nop"
            } else if jmpnop != -1 && ip == nops[jmpnop] {
                parts[0] = "jmp"
            }
            if (parts[0] == "acc") {
                acc += v
            } else if (parts[0] == "jmp") {
                ip += v
                continue
            }
            fmt.Printf("<%d (%s): acc=%d\n", ip, code[ip], acc)
            ip++
        }
        if ok {
            break
        }
    }
    fmt.Println(acc)
}

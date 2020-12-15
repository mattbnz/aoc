// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 15, Puzzle 1
// Rambunctious Recitation.

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"
)

type Spoken map[int][]int

func Speak(spoken Spoken, v int, turn int) Spoken {
    //fmt.Printf("%d: said %d\n", turn, v)
    last, found := spoken[v]
    if !found {
        // Not spoken before
        spoken[v] = []int{turn,0}
        return spoken
    }
    spoken[v] = []int{turn, last[0]}
    return spoken
}

func main() {
    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
        starters := strings.Split(strings.TrimSpace(s.Text()), ",")
        spoken := make(Spoken)
        last := 0
        for turn := 1; turn<=30000000; turn++ {
            speak := 0
            if turn-1 < len(starters) {
                v, err := strconv.Atoi(starters[turn-1])
                if err != nil {
                    log.Fatal("Bad starting number! ", starters[turn-1])
                }
                speak = v
            } else {
                history := spoken[last]
                if history[1] != 0 {
                    speak = history[0] - history[1]
                }
                //fmt.Printf("%d %v %v\n", last, found, history)
            }
            spoken = Speak(spoken, speak, turn)
            last = speak
        }
        fmt.Println(last)
    }
}

// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 1, Puzzle 1
// Who has the most calories?

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
  Undefined int = iota
  Rock
  Paper
  Scissors
)

var Codes map[string]int = map[string]int{"A": Rock, "B": Paper, "C": Scissors, "X": Rock, "Y": Paper, "Z": Scissors}

const Draw = 3
const Win = 6

func main() {
	s := bufio.NewScanner(os.Stdin)
	sum := 0
	for s.Scan() {
    plays := strings.Split(s.Text(), " ")
    them := Codes[plays[0]]
    me := Codes[plays[1]]
    if them == 0 || me == 0 {
      log.Fatal("unknown play %s", s.Text())
    }
    if me > them {
      fmt.Print(plays, " win ", me, " ", Win)
      sum += me + Win
    } else if me == them {
      fmt.Print(plays, " draw ", me, " ", Draw)
      sum += me + Draw
    } else {
      fmt.Print(plays, " loss ", me)
      sum += me
    }
    fmt.Println(" sum ", sum)
	}

	fmt.Println(sum)
}

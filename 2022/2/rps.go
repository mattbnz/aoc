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
var Names map[int]string = map[int]string{Rock: "Rock", Paper: "Paper", Scissors: "Scissors"}

const Loss = 0
const Draw = 3
const Win = 6

type BeatTable map[int]map[int]int

var Beats BeatTable = BeatTable{
	Rock:     {Rock: Draw, Paper: Win, Scissors: Loss},
	Paper:    {Rock: Loss, Paper: Draw, Scissors: Win},
	Scissors: {Rock: Win, Paper: Loss, Scissors: Draw},
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	sum := 0
	for s.Scan() {
		plays := strings.Split(s.Text(), " ")
		them := Codes[plays[0]]
		me := Codes[plays[1]]
		if them == 0 || me == 0 {
			log.Fatalf("unknown play %s", s.Text())
		}
		sum += me + Beats[them][me]
		fmt.Printf("%v %s vs %s, +%d +%d = %d\n", plays, Names[them], Names[me], me, Beats[them][me], sum)
	}

	fmt.Println(sum)
}

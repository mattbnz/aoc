// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 5, Puzzle 1.
// Stack manipulation

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

func printStacks(stacks map[int][]rune) {
	stackNos := []int{}
	for i := range stacks {
		stackNos = append(stackNos, i)
	}
	sort.Ints(stackNos)
	for _, i := range stackNos {
		fmt.Printf(" %d ", i)
	}
	fmt.Println("")
	height := 0
	for true {
		any := false
		for _, i := range stackNos {
			if len(stacks[i]) > height {
				fmt.Printf(" %c ", stacks[i][height])
				any = true
			} else {
				fmt.Printf("   ")
			}
		}
		fmt.Println("")
		height++
		if !any {
			break
		}
	}

}

func makeStack(depth int) []rune {
	s := make([]rune, depth)
	for i := 0; i < depth; i++ {
		s = append(s, ' ')
	}
	return s
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	stacksRe := regexp.MustCompile(`(\s\d+)+`)

	stacks := map[int][]rune{}
	depth := 0
	for s.Scan() {
		if stacksRe.MatchString(s.Text()) {
			break
		}
		stackNo := 1
		for i, c := range s.Text() {
			if i == 0 {
				continue
			}
			if i == ((stackNo-1)*4)+1 {
				stack, exists := stacks[stackNo]
				if !exists {
					stack = makeStack(depth)
				}
				stack = append(stack, c)
				stacks[stackNo] = stack
				stackNo++
			}
		}
	}
	printStacks(stacks)

}

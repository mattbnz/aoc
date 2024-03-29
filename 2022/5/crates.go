// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 5, Puzzle 1.
// Stack manipulation

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func stackIdx(stacks map[int][]rune) ([]int, int) {
	stackNos := []int{}
	maxHeight := 0
	for i := range stacks {
		stackNos = append(stackNos, i)
		if len(stacks[i]) > maxHeight {
			maxHeight = len(stacks[i])
		}
	}
	sort.Ints(stackNos)
	return stackNos, maxHeight
}

func printStacks(stacks map[int][]rune) {
	stackNos, maxHeight := stackIdx(stacks)
	for _, i := range stackNos {
		fmt.Printf(" %d ", i)
	}
	fmt.Println("")
	for height := 0; height < maxHeight; height++ {
		for _, i := range stackNos {
			diff := maxHeight - len(stacks[i])
			idx := height
			if diff > 0 && i < diff {
				fmt.Printf("   ")
				continue
			} else if diff > 0 && height >= diff {
				idx = height - diff
			}
			fmt.Printf(" %c ", stacks[i][idx])
		}
		fmt.Println("")
	}
}

func makeStack(depth int) []rune {
	s := make([]rune, depth)
	for i := 0; i < depth; i++ {
		s = append(s, ' ')
	}
	return s
}

func asInt(c string) int {
	v, err := strconv.Atoi(c)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func move(stacks map[int][]rune, src, dst int) {
	item := ' '
	for i := 0; i < len(stacks[src]); i++ {
		if stacks[src][i] != ' ' {
			item = stacks[src][i]
			stacks[src][i] = ' '
			fmt.Printf("Found item '%c' from %d at %d\n", item, src, i)
			break
		}
	}
	if item == ' ' {
		printStacks(stacks)
		log.Fatalf("Could not find an item to move in stack %d", src)
	}

	for i := len(stacks[dst]) - 1; i >= 0; i-- {
		if stacks[dst][i] != ' ' {
			fmt.Printf("Can't add to %d at %d, contains '%c'\n", dst, i, stacks[dst][i])
			continue
		}
		stacks[dst][i] = item
		fmt.Printf("Added '%c' to %d at %d\n", item, dst, i)
		item = ' '
		break
	}
	if item != ' ' {
		fmt.Printf("Added '%c' to %d at top\n", item, dst)
		stacks[dst] = append([]rune{item}, stacks[dst]...)
	}
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

	moveRe := regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)
	for s.Scan() {
		m := moveRe.FindStringSubmatch(s.Text())
		if len(m) == 0 {
			continue
		}
		fmt.Println(s.Text())
		for i := 0; i < asInt(m[1]); i++ {
			move(stacks, asInt(m[2]), asInt(m[3]))
			printStacks(stacks)
		}
	}

	tops := ""
	stackNos, _ := stackIdx(stacks)
	for _, i := range stackNos {
		for h := 0; h <= len(stacks[i]); h++ {
			if stacks[i][h] != ' ' {
				tops += string(stacks[i][h])
				break
			}
		}
	}
	fmt.Println(tops)

}

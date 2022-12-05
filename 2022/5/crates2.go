// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 5, Puzzle 2.
// Stack manipulation with CrateMover 9001

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

func move(stacks map[int][]rune, src, dst, count int) {
	items := make([]rune, count)
	for i := 0; i < len(stacks[src]); i++ {
		if stacks[src][i] != ' ' {
			if i+count > len(stacks[src]) {
				log.Fatalf("Can't take %d from %d starting at %d (only %d)",
					count, src, i, len(stacks[src]))
			}
			copy(items, stacks[src][i:i+count])
			for j := i; j < i+count; j++ {
				stacks[src][j] = ' '
			}
			fmt.Printf("Found %d items '%v' from %d at %d\n", count, items, src, i)
			break
		}
	}
	if len(items) == 0 {
		printStacks(stacks)
		log.Fatalf("Could not find an item to move in stack %d", src)
	}

	top := 0
	for i := 0; i < len(stacks[dst]); i++ {
		if stacks[dst][i] != ' ' {
			top = i
			break
		}
	}
	stacks[dst] = append(items, stacks[dst][top:]...)
	fmt.Printf("Added %d items '%v' to %d from %d\n", len(items), items, dst, top)
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
		move(stacks, asInt(m[2]), asInt(m[3]), asInt(m[1]))
		printStacks(stacks)
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

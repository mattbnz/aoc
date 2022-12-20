// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 19, Puzzle 1.
// Not Enough Minerals. Robots to the rescue.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Int(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("%s is not an int: %v", s, err)
	}
	return v
}

func Abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

type Num struct {
	StartPos int // maybe useless? keep JIC.
	Value    int

	Next *Num
	Prev *Num
}

func (n *Num) String() string {
	return fmt.Sprintf("%d", n.Value)
}
func (n *Num) Go() *Num {
	return n.GoN(n.Value)
}

func (n *Num) GoN(amount int) *Num {
	num := n
	for i := Abs(amount); i > 0; i-- {
		if amount > 0 {
			num = num.Next
		} else {
			num = num.Prev
		}
	}
	return num
}

func PrintList(start *Num, expectedLen int) {
	item := start
	s := []string{}
	for {
		s = append(s, item.String())
		item = item.Next
		if item == start {
			break
		}
	}
	if len(s) != expectedLen {
		fmt.Printf("WARNING: List is not expected length (%d vs %d)\n", len(s), expectedLen)
	}
	fmt.Println(strings.Join(s, ", "))
	fmt.Println()
}

func Validate(start *Num, expectedItems []*Num) bool {
	item := start
	ok := true
	present := map[int]bool{}
	for {
		if _, exists := present[item.StartPos]; exists {
			fmt.Printf("ERROR: %s from starting position %d has already been seen in the list!!\n", item, item.StartPos)
			ok = false
		}
		present[item.StartPos] = true
		item = item.Next
		if item == start {
			break
		}
	}
	if len(present) != len(expectedItems) {
		fmt.Printf("ERROR: List is not expected length (%d vs %d)\n", len(present), len(expectedItems))
		ok = false
	}
	for _, i := range expectedItems {
		if _, exists := present[i.StartPos]; !exists {
			fmt.Printf("ERROR: Item %s from starting position %d was not found in the list!!\n", i, i.StartPos)
			ok = false
		}
	}
	return ok
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	order := []*Num{}

	// Build a doubly-linked list of numbers
	pos := 0
	var last *Num
	var zero *Num
	for s.Scan() {
		n := &Num{StartPos: pos, Value: Int(s.Text())}
		pos++
		if last != nil {
			last.Next = n
			n.Prev = last
		}
		order = append(order, n)
		last = n
		// Keep track of zero
		if n.Value == 0 {
			if zero != nil {
				log.Fatal("Two Zeros!")
			}
			zero = n
		}
	}

	// Complete the list to be a loop (yowsers!)
	last.Next = order[0]
	order[0].Prev = last
	PrintList(order[0], len(order))
	fmt.Println("List validity: ", Validate(order[0], order))
	fmt.Println()

	// Now do the mixing
	start := order[0]
	for _, n := range order {
		// Go in the direction our value dictates
		o := n.Go()
		if o == n {
			fmt.Println("Found ourself (expected for zero): ", n)
			continue
		}
		if n.Next == o || n.Prev == o {
			fmt.Println(n, " and ", o, " are adjacent!")
		}
		// Take ourselves out of our current position
		if n == start {
			start = n.Next
		}
		n.Prev.Next = n.Next
		n.Next.Prev = n.Prev
		// Put ourselves "after" the item we ended up on.
		if n.Value > 0 {
			if o.Next == start {
				fmt.Println("we're now the start (A)!", n, n.StartPos)
				start = n
			}
			n.Next = o.Next
			o.Next.Prev = n
			n.Prev = o
			o.Next = n
		} else {
			if o.Prev == start {
				fmt.Println("we're now the start (B)!", n, n.StartPos)
				start = n
			}
			n.Prev = o.Prev
			o.Prev.Next = n
			n.Next = o
			o.Prev = n
		}
		if os.Getenv("VALIDATE") == "1" && !Validate(order[0], order) {
			fmt.Println("Validity broken after processing item: ", n)
			PrintList(start, len(order))
		}

	}

	PrintList(start, len(order))

	oneThou := zero.GoN(1000)
	twoThou := oneThou.GoN(1000)
	theThou := twoThou.GoN(1000)

	fmt.Println("Answer via list traversal:")
	fmt.Println(oneThou)
	fmt.Println(twoThou)
	fmt.Println(theThou)
	fmt.Println(oneThou.Value + twoThou.Value + theThou.Value)
	fmt.Println()

	fmt.Println("Cross-check answer via list indexing:")
	end := []*Num{}
	item := start
	n := 0
	var endZero *Num
	endZeroAt := -1
	for {
		end = append(end, item)
		if item.Value == 0 {
			if endZero != nil {
				log.Fatal("two zeros at end!")
			}
			endZero = item
			endZeroAt = n
		}
		n++
		item = item.Next
		if item == start {
			break
		}
	}
	if len(end) != len(order) {
		log.Fatal("end list is not the right lenght")
	}
	if endZero != zero {
		log.Fatal("zero has changed!")
	}
	fmt.Println(len(end), " Items", n)
	fmt.Println("End Zero is item #", endZeroAt)
	i1 := (endZeroAt + 1000) % len(end)
	i2 := (endZeroAt + 2000) % len(end)
	i3 := (endZeroAt + 3000) % len(end)
	fmt.Printf("1000th value is #%d (aka %d) == %d\n", endZeroAt+1000, i1, end[i1].Value)
	fmt.Printf("2000th value is #%d (aka %d) == %d\n", endZeroAt+2000, i2, end[i2].Value)
	fmt.Printf("3000th value is #%d (aka %d) == %d\n", endZeroAt+3000, i3, end[i3].Value)
	fmt.Println(end[i1].Value + end[i2].Value + end[i3].Value)
}

// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 11, Puzzle 1.
// Monkey in the Middle - Monkey Business!

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Monkey struct {
	Queue     []int
	MultBy    int
	AddBy     int
	Divisor   int
	DestTrue  int
	DestFalse int
}

func ParseQueue(l string) []int {
	rv := []int{}
	items := strings.Split(l[16:], ", ")
	for _, i := range items {
		v, err := strconv.Atoi(i)
		if err != nil {
			log.Fatalf("Bad items (%s): %v", l, err)
		}
		rv = append(rv, v)
	}
	return rv
}

var lastnumRE = regexp.MustCompile(`(\d+)$`)
var endoldRE = regexp.MustCompile(` old$`)

func LastNumber(l string) int {
	m := lastnumRE.FindStringSubmatch(l)
	if m == nil {
		log.Fatalf("No last number in %s", l)
	}
	v, err := strconv.Atoi(m[0])
	if err != nil {
		log.Fatalf("Bad last number in %s: %v", l, err)
	}
	return v
}

func ParseOp(l string) (int, int) {
	m := 0
	a := 0
	v := 0
	if endoldRE.MatchString(l) {
		v = -1
	} else {
		v = LastNumber(l)
	}
	if l[21] == '*' {
		m = v
	} else if l[21] == '+' {
		a = v
	} else {
		log.Fatalf("Couldn't parse %s: %c is not a known operation", l, l[22])
	}
	return m, a
}

func PrintItems(q []int) string {
	s := []string{}
	for _, i := range q {
		s = append(s, fmt.Sprintf("%d", i))
	}
	return strings.Join(s, ", ")
}

func ValFor(v int) string {
	if v == -1 {
		return "old"
	}
	return fmt.Sprintf("%d", v)
}

func PrintOp(m *Monkey) string {
	if m.MultBy != 0 {
		return fmt.Sprintf("* %s", ValFor(m.MultBy))
	} else if m.AddBy != 0 {
		return fmt.Sprintf("+ %s", ValFor(m.AddBy))
	}
	log.Fatal("Monkey has bad state")
	return ""
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	monkeys := []*Monkey{}
	current := -1

	for s.Scan() {
		l := strings.TrimSpace(s.Text())
		if strings.HasPrefix(l, "Monkey ") {
			monkeys = append(monkeys, &Monkey{})
			current++
		} else if strings.HasPrefix(l, "Starting") {
			monkeys[current].Queue = ParseQueue(l)
		} else if strings.HasPrefix(l, "Operation") {
			monkeys[current].MultBy, monkeys[current].AddBy = ParseOp(l)
		} else if strings.HasPrefix(l, "Test") {
			monkeys[current].Divisor = LastNumber(l)
		} else if strings.HasPrefix(l, "If true") {
			monkeys[current].DestTrue = LastNumber(l)
		} else if strings.HasPrefix(l, "If false") {
			monkeys[current].DestFalse = LastNumber(l)
		}
	}

	for i, m := range monkeys {
		fmt.Printf("Monkey %d\n", i)
		fmt.Printf("  Starting items: %s\n", PrintItems(m.Queue))
		fmt.Printf("  Operation: new = old %s\n", PrintOp(m))
		fmt.Printf("  Test: divisible by %d\n", m.Divisor)
		fmt.Printf("    If true: throw to monkey %d\n", m.DestTrue)
		fmt.Printf("    If false: throw to monkey %d\n", m.DestFalse)
		fmt.Println("")
	}

}

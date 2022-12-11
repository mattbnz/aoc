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
	"sort"
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

	Inspections int
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

var DBG = true

func DPrintln(a ...any) {
	if DBG {
		fmt.Println(a...)
	}
}

func DPrintf(f string, args ...any) {
	if DBG {
		fmt.Printf(f, args...)
	}
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
		DPrintf("Monkey %d\n", i)
		DPrintf("  Starting items: %s\n", PrintItems(m.Queue))
		DPrintf("  Operation: new = old %s\n", PrintOp(m))
		DPrintf("  Test: divisible by %d\n", m.Divisor)
		DPrintf("    If true: throw to monkey %d\n", m.DestTrue)
		DPrintf("    If false: throw to monkey %d\n", m.DestFalse)
		DPrintln("")
	}

	for i := 0; i < 20; i++ {
		for n, m := range monkeys {
			DPrintf("Monkey %d:\n", n)
			for _, w := range m.Queue {
				m.Inspections++
				DPrintf("  Monkey inspects an item with worry level of %d.\n", w)
				if m.MultBy > 0 {
					w *= m.MultBy
				} else if m.AddBy > 0 {
					w += m.AddBy
				} else if m.MultBy == -1 {
					w *= w
				} else if m.AddBy == -1 {
					w += w
				}
				DPrintf("    Worry level is now %d.\n", w)
				w /= 3
				DPrintf("    Monkey gets bored with item. Worry level is divided by 3 to %d.\n", w)
				next := -1
				if w%m.Divisor == 0 {
					DPrintf("    Current worry level is divisible by %d\n", m.Divisor)
					next = m.DestTrue
				} else {
					DPrintf("    Current worry level is not divisible by %d\n", m.Divisor)
					next = m.DestFalse
				}
				nextM := monkeys[next]
				nextM.Queue = append(nextM.Queue, w)
				DPrintf("    Item with worry level %d is thrown to monkey %d.\n", w, next)
			}
			m.Queue = []int{}
		}

		for n, m := range monkeys {
			fmt.Printf("Monkey %d: %s\n", n, PrintItems(m.Queue))
		}
	}

	counts := []int{}
	for n, m := range monkeys {
		fmt.Printf("Monkey %d inspected items %d times.\n", n, m.Inspections)
		counts = append(counts, m.Inspections)
	}
	sort.Ints(counts)
	fmt.Println(counts[len(counts)-2] * counts[len(counts)-1])
}

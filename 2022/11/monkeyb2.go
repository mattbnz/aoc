// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 11, Puzzle 1.
// Monkey in the Middle - Monkey Business!
//
// This was a failed/naive attempt at using
// bigints to solve the issue - completely
// wrong, does not work.

package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var BIGZERO big.Int
var BIGM1 *big.Int = big.NewInt(-1)

type Monkey struct {
	Queue     []big.Int
	MultBy    big.Int
	AddBy     big.Int
	Divisor   big.Int
	DestTrue  int
	DestFalse int

	Inspections int
}

func ParseQueue(l string) []big.Int {
	rv := []big.Int{}
	if len(l) < 16 {
		return rv
	}
	items := strings.Split(l[16:], ", ")
	for _, i := range items {
		v, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			log.Fatalf("Bad items (%s): %v", l, err)
		}
		b := big.Int{}
		b.SetInt64(v)
		rv = append(rv, b)
	}
	return rv
}

var lastnumRE = regexp.MustCompile(`(\d+)$`)
var endoldRE = regexp.MustCompile(` old$`)

func LastNumber(l string) big.Int {
	m := lastnumRE.FindStringSubmatch(l)
	if m == nil {
		log.Fatalf("No last number in %s", l)
	}
	v, err := strconv.ParseInt(m[0], 10, 64)
	if err != nil {
		log.Fatalf("Bad last number in %s: %v", l, err)
	}
	b := big.Int{}
	b.SetInt64(v)
	return b
}

func ParseOp(l string) (big.Int, big.Int) {
	m := big.Int{}
	a := big.Int{}
	v := big.Int{}
	if endoldRE.MatchString(l) {
		v.SetInt64(-1)
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

func PrintItems(q []big.Int) string {
	s := []string{}
	for _, i := range q {
		s = append(s, i.String())
	}
	return strings.Join(s, ", ")
}

func ValFor(v big.Int) string {
	if v.Int64() == -1 {
		return "old"
	}
	return v.String()
}

func PrintOp(m *Monkey) string {
	if m.MultBy.Int64() != 0 {
		return fmt.Sprintf("* %s", ValFor(m.MultBy))
	} else if m.AddBy.Int64() != 0 {
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
			t := LastNumber(l)
			monkeys[current].DestTrue = int(t.Int64())
		} else if strings.HasPrefix(l, "If false") {
			t := LastNumber(l)
			monkeys[current].DestFalse = int(t.Int64())
		}
	}

	for i, m := range monkeys {
		DPrintf("Monkey %d\n", i)
		DPrintf("  Starting items: %s\n", PrintItems(m.Queue))
		DPrintf("  Operation: new = old %s\n", PrintOp(m))
		DPrintf("  Test: divisible by %s\n", m.Divisor.String())
		DPrintf("    If true: throw to monkey %d\n", m.DestTrue)
		DPrintf("    If false: throw to monkey %d\n", m.DestFalse)
		DPrintln("")
	}

	for i := 0; i < 10000; i++ {
		for n, m := range monkeys {
			DPrintf("Monkey %d:\n", n)
			for _, w := range m.Queue {
				m.Inspections++
				DPrintf("  Monkey inspects an item with worry level of %s.\n", w.String())
				if m.MultBy.Cmp(&BIGZERO) > 0 {
					w = *w.Mul(&w, &m.MultBy)
				} else if m.AddBy.Cmp(&BIGZERO) > 0 {
					w = *w.Add(&w, &m.AddBy)
				} else if m.MultBy.Cmp(BIGM1) == 0 {
					w = *w.Mul(&w, &w)
				} else if m.AddBy.Cmp(BIGM1) == 0 {
					w = *w.Add(&w, &w)
				}
				DPrintf("    Worry level is now %s.\n", w.String())
				next := -1
				t := big.Int{}
				if t.Mod(&w, &m.Divisor).Int64() == 0 {
					DPrintf("    Current worry level is divisible by %s\n", m.Divisor.String())
					next = m.DestTrue
					w.Div(&w, &m.Divisor)
				} else {
					DPrintf("    Current worry level is not divisible by %s\n", m.Divisor.String())
					next = m.DestFalse
				}
				//fmt.Println(next)
				nextM := monkeys[next]
				nextM.Queue = append(nextM.Queue, w)
				DPrintf("    Item with worry level %s is thrown to monkey %d.\n", w.String(), next)
			}
			m.Queue = []big.Int{}
		}

		for n, m := range monkeys {
			DPrintf("Monkey %d: %s\n", n, PrintItems(m.Queue))
		}
		DPrintln("")
		if (i > 0 && i < 1000 && i%20 == 0) || (i > 1000 && i%1000 == 0) {
			fmt.Printf("== After round %d ==\n", i)
			for n, m := range monkeys {
				fmt.Printf("Monkey %d inspected items %d times.\n", n, m.Inspections)
			}
			fmt.Println("")
			break
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

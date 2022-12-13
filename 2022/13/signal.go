// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 13, Puzzle 1.
// Distress Signal - ordering.

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

func Min(a, b int) int {
	if a == -1 {
		return b
	}
	if b == -1 {
		return a
	}
	if a < b {
		return a
	}
	return b
}

type Signal []any

func (s Signal) String() string {
	v := []string{}
	for _, e := range s {
		if i, ok := e.(int); ok {
			v = append(v, fmt.Sprintf("%d", i))
		} else if l, ok := e.(Signal); ok {
			v = append(v, fmt.Sprintf(" %s ", l.String()))
		} else {
			log.Fatalf("Unknown type in signal: %v", e)
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(v, ","))
}

func ParseSignal(s string) (Signal, int) {
	if s[0] != '[' {
		return Signal{Int(s)}, len(s)
	}
	rv := Signal{}
	i := 1
	//fmt.Println(len(s), s)
	for i < len(s) {
		//fmt.Printf("%d: %c", i, s[i])
		if s[i] == ',' {
			//fmt.Println()
			i++
		} else if s[i] == ']' {
			//fmt.Println(" end")
			i++
			break
		} else if s[i] == '[' {
			//fmt.Println(" going deeper")
			e, l := ParseSignal(s[i:])
			rv = append(rv, e)
			//fmt.Println(" and back, used ", i, " to ", i+l)
			i += l
		} else {
			c := strings.Index(s[i:], ",")
			e := strings.Index(s[i:], "]")
			if c == -1 && e == -1 {
				log.Fatalf("Found unterminated list at %d parsing: %s", i, s)
			}
			n := Min(c, e)
			rv = append(rv, Signal{Int(s[i : i+n])})
			//fmt.Printf(" used from %d to %d\n", i, i+n)
			i += n + 1
			if e == n {
				break
			}
		}
	}
	return rv, i
}

func IsOrdered(a, b Signal) bool {
	fmt.Println("A", a)
	fmt.Println("B", b)
	fmt.Println()
	return false
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	sum := 0
	pair := 1

	signals := []Signal{}
	for s.Scan() {
		if s.Text() == "" {
			if len(signals) != 2 {
				log.Fatal("blank line but don't have 2 signals!")
			}
			if IsOrdered(signals[0], signals[1]) {
				sum += pair
			}
			pair++
			signals = []Signal{}
		} else {
			e, used := ParseSignal(s.Text())
			if used < len(s.Text()) {
				log.Fatalf("Parsing of %s only consumed %d/%d characters!", s.Text(), used, len(s.Text()))
			}
			signals = append(signals, e)
		}
	}
	if len(signals) == 2 {
		if IsOrdered(signals[0], signals[1]) {
			sum += pair
		}
	}
	fmt.Println(sum)
}

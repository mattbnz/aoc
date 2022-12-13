// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 13, Puzzle 1.
// Distress Signal - sorting a list.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
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

func (s Signal) Compare(o Signal) int {
	i := 0
	for i < len(s) || i < len(o) {
		if len(s) <= i {
			return -1 // left ran out first!
		}
		if len(o) <= i {
			return 1 // right ran out first!
		}
		sL, slOk := s[i].(Signal)
		oL, olOk := o[i].(Signal)
		if slOk && olOk {
			if v := sL.Compare(oL); v != 0 {
				return v
			}
			i++
			continue
		}
		sI, siOk := s[i].(int)
		oI, oiOk := o[i].(int)
		if siOk && oiOk {
			if sI > oI {
				return 1
			} else if sI < oI {
				return -1
			}
			i++
			continue
		}
		if slOk && oiOk {
			if v := sL.Compare(Signal{oI}); v != 0 {
				return v
			}
		} else if siOk && olOk {
			if v := (Signal{sI}).Compare(oL); v != 0 {
				return v
			}
		} else {
			log.Fatal("mismatch at ", i)
		}
		i++
	}

	// both ended - so equal
	return 0
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
			rv = append(rv, Int(s[i:i+n]))
			//fmt.Printf(" used from %d to %d\n", i, i+n)
			i += n + 1
			if e == n {
				break
			}
		}
	}
	return rv, i
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	signals := []Signal{{Signal{2}}, {Signal{6}}}
	key1 := signals[0].String()
	key2 := signals[1].String()
	for s.Scan() {
		if s.Text() == "" {
			continue
		}
		e, used := ParseSignal(s.Text())
		if used < len(s.Text()) {
			log.Fatalf("Parsing of %s only consumed %d/%d characters!", s.Text(), used, len(s.Text()))
		}
		signals = append(signals, e)
	}

	sort.SliceStable(signals, func(i, j int) bool {
		return signals[i].Compare(signals[j]) == -1
	})

	sum := 1
	for i, s := range signals {
		k := s.String()
		fmt.Printf(k)
		if k == key1 || k == key2 {
			sum *= (i + 1)
			fmt.Print("  ***** ", i+1)
		}
		fmt.Println()
	}
	fmt.Println(sum)
}

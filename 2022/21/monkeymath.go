// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 19, Puzzle 1.
// Grove Positioning System - decryption mixing.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
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

type MonkeyMap map[string]*Monkey

type Monkey struct {
	Name   string
	Number int // for announcing monkeys
	// for computing monkeys
	Left  string
	Op    string
	Right string
}

func (m *Monkey) String() string {
	if m.Number != 0 {
		return fmt.Sprintf("%s: %d", m.Name, m.Number)
	} else {
		return fmt.Sprintf("%s: %s %s %s", m.Name, m.Left, m.Op, m.Right)
	}
}

func (m *Monkey) Result(monkeys MonkeyMap) int {
	if m.Number != 0 {
		return m.Number
	}

	leftV := monkeys[m.Left].Result(monkeys)
	rightV := monkeys[m.Right].Result(monkeys)
	if m.Op == "+" {
		return leftV + rightV
	} else if m.Op == "-" {
		return leftV - rightV
	} else if m.Op == "*" {
		return leftV * rightV
	} else if m.Op == "/" {
		return leftV / rightV
	}
	log.Fatal("Bad Monkey Result: ", m)
	return -1
}

var MATH_MONKEY_RE = regexp.MustCompile(`([a-z]+): ([a-z]+) ([+*/-]) ([a-z]+)`)
var ANNOUNCING_MONKEY_RE = regexp.MustCompile(`([a-z]+): ([\d]+)$`)

func NewMonkey(s string) *Monkey {
	m := MATH_MONKEY_RE.FindStringSubmatch(s)
	if len(m) == 5 {
		return &Monkey{
			Name:  m[1],
			Left:  m[2],
			Op:    m[3],
			Right: m[4],
		}

	}
	m = ANNOUNCING_MONKEY_RE.FindStringSubmatch(s)
	if len(m) == 3 {
		return &Monkey{
			Name:   m[1],
			Number: Int(m[2]),
		}
	}
	log.Fatal("Bad Monkey: ", s)
	return nil
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	monkeys := MonkeyMap{}
	for s.Scan() {
		m := NewMonkey(s.Text())
		monkeys[m.Name] = m
	}

	root := monkeys["root"]
	if root == nil {
		log.Fatal("no root monkey")
	}
	fmt.Println(root.Result(monkeys))
}

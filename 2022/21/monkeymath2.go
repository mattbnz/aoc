// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 21, Puzzle 2.
// // Monkey Math. Algebra Solver...

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

var REVERSE = map[string]string{
	"+": "-",
	"-": "+",
	"*": "/",
	"/": "*",
}

func Calc(a int, op string, b int) int {
	if op == "+" {
		return a + b
	} else if op == "-" {
		return a - b
	} else if op == "*" {
		return a * b
	} else if op == "/" {
		return a / b
	}
	log.Fatal("Can't calc bad op: ", a, op, b)
	return -1
}

type MonkeyMap map[string]*Monkey

type Result struct {
	Value int

	Left  *Result
	Op    string
	Right *Result

	X bool
}

func (r Result) String() string {
	if r.X {
		return "X"
	}
	if r.Value > 0 {
		return fmt.Sprintf("%d", r.Value)
	}
	return fmt.Sprintf("(%s %s %s)", r.Left, r.Op, r.Right)
}

func (r Result) SolveFor(v int) int {
	if r.X {
		fmt.Println("X must be ", v)
		return v
	}
	solve := r.Left
	other := r.Right.Value
	if r.Right.Value == -1 {
		solve = r.Right
		other = r.Left.Value
	}
	newV := Calc(v, REVERSE[r.Op], other)
	fmt.Printf("%d %s %d = %d\n", v, REVERSE[r.Op], other, newV)
	return solve.SolveFor(newV)
}

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

func (m *Monkey) Calc(monkeys MonkeyMap) int {
	if m.Number != 0 {
		return m.Number
	}

	leftV := monkeys[m.Left].Calc(monkeys)
	rightV := monkeys[m.Right].Calc(monkeys)
	return Calc(leftV, m.Op, rightV)
}

func (m *Monkey) Result(monkeys MonkeyMap) *Result {
	if m.Name == "humn" {
		return &Result{X: true}
	}
	if m.Number != 0 {
		return &Result{Value: m.Number}
	}

	leftV := monkeys[m.Left].Result(monkeys)
	rightV := monkeys[m.Right].Result(monkeys)
	if leftV.Value > 0 && rightV.Value > 0 {
		return &Result{Value: Calc(leftV.Value, m.Op, rightV.Value)}
	}
	return &Result{Value: -1, Left: leftV, Op: m.Op, Right: rightV}
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
	left := monkeys[root.Left]
	right := monkeys[root.Right]

	orig := root.Calc(monkeys)
	fmt.Println("Original Result: ", orig)
	origLeft := left.Calc(monkeys)
	origRight := right.Calc(monkeys)
	fmt.Println(origLeft, root.Op, origRight, " = ", orig)
	fmt.Println()

	leftV := left.Result(monkeys)
	fmt.Println(root.Left, " equation is: ", leftV)
	rightV := right.Result(monkeys)
	fmt.Println(root.Right, "expected value is: ", rightV)

	x := 0
	if leftV.Value != -1 {
		log.Fatal("expected X to be in left")
	}
	fmt.Println("Test SolveFor against original: ", origLeft)
	x = leftV.SolveFor(origLeft)
	fmt.Println(x)
	humn := monkeys["humn"]
	if x != humn.Number {
		log.Fatal("Solved value doesn't match original hummn: ", humn.Number)
	}
	fmt.Println("Solving Left to equal", rightV.Value)
	x = leftV.SolveFor(rightV.Value)
	fmt.Println(x)

	fmt.Println("Using it gives:")
	humn.Number = x
	newLeft := left.Calc(monkeys)

	fmt.Println("Left: ", newLeft)
	fmt.Println("Right: ", right.Calc(monkeys), " vs ", origRight)
	if newLeft != origRight {
		fmt.Println(x, " cannot be right answer!")
	} else {
		fmt.Println(x, " looks like the right answer!")
	}
}

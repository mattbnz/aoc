// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 21, Puzzle 2.
// // Monkey Math. Algebra Solver...

package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"regexp"
	"strconv"
)

func Int(s string) *big.Int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("%s is not an int: %v", s, err)
	}
	return big.NewInt(int64(v))
}

var REVERSE = map[string]string{
	"+": "-",
	"-": "+",
	"*": "/",
	"/": "*",
}

func Calc(a *big.Int, op string, b *big.Int) *big.Int {
	var r big.Int
	if op == "+" {
		return r.Add(a, b)
	} else if op == "-" {
		return r.Sub(a, b)
	} else if op == "*" {
		return r.Mul(a, b)
	} else if op == "/" {
		return r.Div(a, b)
	}
	log.Fatal("Can't calc bad op: ", a, op, b)
	return big.NewInt(-1)
}

var INVALID = big.NewInt(-1)

type MonkeyMap map[string]*Monkey

type Result struct {
	Value *big.Int
	From  string

	Left  *Result
	Op    string
	Right *Result

	X bool
}

func (r Result) String() string {
	if r.X {
		return "X"
	}
	if r.Value.Cmp(&big.Int{}) == 1 {
		return fmt.Sprintf("%d", r.Value)
	}
	return fmt.Sprintf("(%s %s %s)", r.Left, r.Op, r.Right)
}

func (r Result) StringAsNames() string {
	if r.X {
		return "X"
	}
	if r.Value.Cmp(&big.Int{}) == 1 {
		if r.From != "" {
			return fmt.Sprintf("%s", r.From)
		}
		return fmt.Sprintf("%d", r.Value)
	}
	return fmt.Sprintf("(%s %s %s)", r.Left.StringAsNames(), r.Op, r.Right.StringAsNames())
}

func (r Result) SolveFor(v *big.Int) *big.Int {
	//fmt.Print("SolveFor ", r, " == ", v)
	solve := r.Left
	other := r.Right.Value
	reverseX := true
	if r.Right.Value == nil || r.Right.Value.Cmp(INVALID) == 0 {
		solve = r.Right
		other = r.Left.Value
		reverseX = false
	}
	//fmt.Println("solve", solve)
	//fmt.Println("other", other)
	if solve.X {
		if reverseX {
			newV := Calc(v, REVERSE[r.Op], other)
			fmt.Printf(" ===> %d %s %d = %d == X!!\n", v, REVERSE[r.Op], other, newV)
			return newV
		} else {
			newV := Calc(other, r.Op, v)
			fmt.Printf(" ===> %d %s %d = %d == X!!\n", other, r.Op, v, newV)
			return newV
		}
	}
	newV := Calc(v, REVERSE[r.Op], other)
	fmt.Printf(" <===> == %d %s %d = %d\n", v, REVERSE[r.Op], other, newV)
	return solve.SolveFor(newV)
}

type Monkey struct {
	Name   string
	Number *big.Int // for announcing monkeys
	// for computing monkeys
	Left  string
	Op    string
	Right string
}

func (m *Monkey) String() string {
	if m.Number != nil {
		return fmt.Sprintf("%s: %d", m.Name, m.Number)
	} else {
		return fmt.Sprintf("%s: %s %s %s", m.Name, m.Left, m.Op, m.Right)
	}
}

func (m *Monkey) Calc(monkeys MonkeyMap) *big.Int {
	if m.Number != nil {
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
	if m.Number != nil {
		return &Result{Value: m.Number, From: m.Name}
	}

	leftV := monkeys[m.Left].Result(monkeys)
	rightV := monkeys[m.Right].Result(monkeys)
	if leftV.Value != nil && leftV.Value != INVALID && rightV.Value != nil && rightV.Value != INVALID {
		return &Result{Value: Calc(leftV.Value, m.Op, rightV.Value)}
	}
	fmt.Println(m, ": L", leftV)
	fmt.Println(m, ": R", rightV)
	return &Result{Value: INVALID, Left: leftV, Op: m.Op, Right: rightV, From: m.Name}
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

	if len(os.Args) > 1 {
		m := monkeys[os.Args[1]]
		if m == nil {
			log.Fatal("Bad Monkey: ", m)
		}
		r := m.Result(monkeys)
		fmt.Println(m, ": ", r.StringAsNames())
		fmt.Println(m, ": ", r)
		return
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
	fmt.Println()
	fmt.Println("Equation: ", root.Left, ":", leftV.StringAsNames())
	fmt.Println("Equation: ", root.Left, ":", leftV)

	x := big.NewInt(0)
	if leftV.Value.Cmp(INVALID) != 0 {
		log.Fatal("expected X to be in left")
	}

	fmt.Println("Solving Left for original value as test: ", origLeft) //rightV.Value)
	x = leftV.SolveFor(origLeft)                                       // rightV.Value)
	fmt.Println(x)
	humn := monkeys["humn"]
	if x.Cmp(humn.Number) != 0 {
		log.Fatal("Found value doesn't match original hummn: ", humn.Number)
	}

	fmt.Println("Solving Left to equal", origRight) //rightV.Value)
	x = leftV.SolveFor(origRight)                   // rightV.Value)
	fmt.Println(x)

	fmt.Println("Using it gives:")
	humn.Number = x
	newLeft := left.Calc(monkeys)

	fmt.Println("Left: ", newLeft)
	fmt.Println("Right: ", right.Calc(monkeys), " vs ", origRight)
	if newLeft.Cmp(origRight) != 0 {
		fmt.Println(x, " cannot be right answer!")
	} else {
		fmt.Println(x, " looks like the right answer!")
	}
	var v big.Int
	fmt.Println(v.Add(x, newLeft))
}

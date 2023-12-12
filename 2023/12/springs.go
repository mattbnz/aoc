package day12

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Spring int

func (s *Spring) New(t rune) (Spring, error) {
	if t == '?' {
		return 0, nil
	} else if t == '.' {
		return 1, nil
	} else if t == '#' {
		return 2, nil
	}
	return -1, fmt.Errorf("unknown spring type")
}

func (s *Spring) String() string {
	if *s == 0 {
		return "?"
	} else if *s == 1 {
		return "."
	} else if *s == 2 {
		return "#"
	}
	return "!"
}

type SpringRow struct {
	Springs []Spring
	BadRuns []int
}

type SpringRows []SpringRow

func NewSpringRows(filename string) (rows SpringRows, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	var sFactory Spring
	s := bufio.NewScanner(f)

	for s.Scan() {
		springs, seqs, ok := strings.Cut(s.Text(), " ")
		if !ok {
			return nil, fmt.Errorf("bad format for input: %s", s.Text())
		}
		sr := SpringRow{BadRuns: numberList(seqs)}
		for _, c := range springs {
			spring, err := sFactory.New(c)
			if err != nil {
				return nil, fmt.Errorf("bad spring (%c) for input (%s): %v", c, s.Text(), err)
			}
			sr.Springs = append(sr.Springs, spring)
		}
		rows = append(rows, sr)
	}
	return
}

func (r *SpringRows) Arrangements(row int) int {
	return 0
}

func (r *SpringRows) SumArrangements() (rv int) {
	for n, _ := range *r {
		rv += r.Arrangements(n)
	}
	return
}

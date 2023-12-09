// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 8.
// Haunted Wasteland.

package day9

import (
	"bufio"
	"os"
)

type Reading []int

type Scan struct {
	Readings []Reading
}

func NewScan(filename string) (scan Scan, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		scan.Readings = append(scan.Readings, numberList(s.Text()))
	}
	return
}

func (s *Scan) ExtrapolateAndSum() int {
	sum := 0
	for _, seq := range s.Readings {
		sum += s.Extrapolate(seq)
	}
	return sum
}

func (s *Scan) Extrapolate(r Reading) int {
	diffs := Reading{}
	last := r[0]
	for n := 1; n < len(r); n++ {
		diffs = append(diffs, r[n]-last)
		last = r[n]
	}
	if Sum(diffs) == 0 {
		return r[len(r)-1]
	}
	nextDiff := s.Extrapolate(diffs)
	return r[len(r)-1] + nextDiff
}

func (s *Scan) ExtrapolateHistoryAndSum() int {
	sum := 0
	for _, seq := range s.Readings {
		sum += s.ExtrapolateHistory(seq)
	}
	return sum
}

func (s *Scan) ExtrapolateHistory(r Reading) int {
	diffs := Reading{}
	last := r[0]
	for n := 1; n < len(r); n++ {
		diffs = append(diffs, r[n]-last)
		last = r[n]
	}
	if Sum(diffs) == 0 {
		return r[0]
	}
	nextDiff := s.ExtrapolateHistory(diffs)
	return r[0] - nextDiff
}

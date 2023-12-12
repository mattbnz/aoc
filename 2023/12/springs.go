package day12

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
)

type Spring int

const (
	S_UNKNOWN Spring = iota
	S_OK
	S_BAD
)

func (s *Spring) New(t rune) (Spring, error) {
	if t == '?' {
		return S_UNKNOWN, nil
	} else if t == '.' {
		return S_OK, nil
	} else if t == '#' {
		return S_BAD, nil
	}
	return -1, fmt.Errorf("unknown spring type")
}

func (s Spring) String() string {
	if s == S_UNKNOWN {
		return "?"
	} else if s == S_OK {
		return "."
	} else if s == S_BAD {
		return "#"
	}
	return "!"
}

type SpringList []Spring

func (sl SpringList) String() (rv string) {
	for _, s := range sl {
		rv += s.String()
	}
	return
}

// NewWith returns a new SpringList with the spring at n, in state
func (sl *SpringList) NewWith(n int, state Spring) (rv SpringList) {
	// Found first unknown spring
	rv = append(rv, (*sl)[:n]...)
	rv = append(rv, state)
	rv = append(rv, (*sl)[n+1:]...)
	return
}

// Arrangements recursively computes the potential arrangements that can match spec by replacing unknown springs with either good or bad ones.
func (sl *SpringList) Arrangements(spec Ints) int {
	return sl.findArrangements(0, 0, -1, spec)
}

// findArrangements iteratively checks along the list of springs, recursing only when an unknown spring is found following an already valid
// pattern.
func (sl *SpringList) findArrangements(from, sI, badTil int, spec Ints) int {
	for n := from; n < len(*sl); n++ {
		s := (*sl)[n]
		if s == S_UNKNOWN {
			a := sl.NewWith(n, S_OK)
			b := sl.NewWith(n, S_BAD)
			return a.findArrangements(n, sI, badTil, spec) + b.findArrangements(n, sI, badTil, spec)
		}
		if s == S_OK {
			if badTil != -1 && n <= badTil {
				glog.V(2).Infof("good spring at %d when expecting bad til %d in %s, %s", n, badTil, sl, spec)
				return 0 // found a good spring when expecting a run of bads
			}
			if n > badTil {
				badTil = -1 // reset expectation
				continue
			}
			continue
		}

		if badTil != -1 && n == badTil+1 {
			// bad following expected end of sequence, when we need an OK!
			glog.V(2).Infof("bad spring at %d needs to be OK in %s, %s", n, sl, spec)
			return 0 // too many bad springs
		}

		if badTil == -1 {
			if sI >= len(spec) {
				glog.V(2).Infof("bad spring at %d makes %d bad runs, when expecting %d in %s, %s", n, sI, len(spec), sl, spec)
				return 0 // too many bad springs
			}
			// first bad spring in a run
			badTil = n + spec[sI] - 1 // s counts as the first item in the run
			sI++
		} else {
			if n > badTil {
				glog.V(2).Infof("bad spring at %d exceeds expected run til %d, when expecting %d in %s, %s", n, badTil, sl, spec)
				return 0 // sequence too long
			}
		}
	}
	if badTil != -1 && badTil >= len(*sl) {
		// bad sequence ran off the end
		glog.V(2).Infof("ran out of bad springs, when expecting til %d in %s of length %d with %s", badTil, sl, len(*sl), spec)
		return 0
	}
	if sI != len(spec) {
		glog.V(2).Infof("only matched %d/%d bad spring sequences in %s with %s", sI, len(spec), sl, spec)
		return 0
	}
	return 1
}

type SpringRow struct {
	Springs SpringList
	BadRuns Ints
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
		sr := SpringRow{BadRuns: NewIntsFromCSV(seqs)}
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
	if row < 0 || row >= len(*r) {
		glog.Fatalf("invalid row index: %d", row)
	}
	sr := (*r)[row]
	return sr.Springs.Arrangements(sr.BadRuns)
}

func (r *SpringRows) SumArrangements() (rv int) {
	for n, _ := range *r {
		rv += r.Arrangements(n)
	}
	return
}

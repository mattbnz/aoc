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
func (sl *SpringList) Arrangements(spec Ints, unfoldFactor int) int {
	return sl.findArrangements(0, 0, -1, spec, unfoldFactor)
}

// Returns true if pattern can match sl[from:from+len(pattern)] (e.g. each item in SL either directly matches pattern, or is unknown).
func (sl *SpringList) canMatch(pattern SpringList, from int) bool {
	if from+len(pattern) > len(*sl) {
		glog.V(1).Infof("cannot match %s from %d in row of only %d", pattern, from, len(*sl))
		return false
	}
	for n := 0; n < len(pattern); n++ {
		if (*sl)[from+n] == S_UNKNOWN {
			continue
		}
		if (*sl)[from+n] != pattern[n] {
			glog.V(1).Infof("pattern mismatch at %d (pattern idx %d, expected %s got %s): %s vs %s", from+n, n, pattern[n], (*sl)[from+n], (*sl)[from:], pattern)
			return false
		}
	}
	return true
}

type Match struct {
	Springs SpringList
	Length  int
}

// Matches returns the positions a "link" spring after a match of the pattern in spec could be, for all matches in sl.
func (sl SpringList) potentialMatches(spec Ints, from int) (rv []Match) {
	glog.V(1).Infof(" canMatchTil %s from %d (of %d)", sl, from, len(sl))
	sI := 0
	badTil := -1
	for n := from; n < len(sl); n++ {
		s := sl[n]
		if s == S_UNKNOWN {
			rv = append(rv, sl.NewWith(n, S_OK).potentialMatches(spec, from)...)
			rv = append(rv, sl.NewWith(n, S_BAD).potentialMatches(spec, from)...)
			return
		}
		if s == S_OK && badTil != -1 && n <= badTil {
			glog.V(2).Infof("good spring at %d when expecting bad til %d in %s, %s", n, badTil, sl, spec)
			return // found a good spring when expecting a run of bads
		}
		if s != S_BAD {
			if n > badTil {
				if sI >= len(spec) {
					return []Match{{Springs: sl, Length: n}} // found end of sequence that matched spec.
				}
				badTil = -1 // reset expectation
			}
			continue
		}

		if badTil != -1 && n == badTil+1 {
			// bad following expected end of sequence, when we need an OK!
			glog.V(2).Infof("bad spring at %d needs to be OK in %s, %s", n, sl, spec)
			return // too many bad springs
		}

		if badTil == -1 {
			if sI >= len(spec) {
				glog.Fatalf("bad spring at %d makes %d bad runs, when expecting %d in %s, %s", n, sI, len(spec), sl, spec)
			}
			// first bad spring in a run
			badTil = n + spec[sI] - 1 // s counts as the first item in the run
			sI++
		} else {
			if n > badTil {
				glog.V(2).Infof("bad spring at %d exceeds expected run til %d, when expecting %d in %s, %s", n, badTil, sl, spec)
				return // sequence too long
			}
		}
	}
	if badTil != -1 && badTil >= len(sl) {
		// bad sequence ran off the end
		glog.V(2).Infof("ran out of bad springs, when expecting til %d in %s of length %d with %s", badTil, sl, len(sl), spec)
		return
	}
	if sI != len(spec) {
		glog.V(2).Infof("only matched %d/%d bad spring sequences in %s with %s", sI, len(spec), sl, spec)
		return
	}
	return []Match{{Springs: sl, Length: len(sl)}}
}

// Returns how many springs backwards the last bad spring was. Zero if no previous bad spring was found.
func (sl *SpringList) lastBad(from int) int {
	for n := from - 1; n >= 0; n-- {
		if (*sl)[n] == S_BAD {
			return from - n
		}
	}
	return 0
}

// findArrangements iteratively checks along the list of springs, recursing only when an unknown spring is found following an already valid
// pattern.
func (sl *SpringList) findArrangements(from, match, badTil int, spec Ints, matchLimit int) (rv int) {
	glog.V(1).Infof(" %s%s from %d (of %d) at match %d", "", sl, from, len(*sl), match)

	matches := sl.potentialMatches(spec, from)
	match++

MATCHES:
	for _, m := range matches {
		glog.V(1).Infof("Found %d match from %d to %d: %s", match, from, m.Length, m.Springs[from:m.Length])
		for linkPos := m.Length; linkPos < len(*sl); linkPos++ {
			spring := m.Springs[linkPos]
			// don't need to check for BAD before unknown, OK here, because potentialMatches will only return options with a following OK/UNKNOWN or end.
			if spring == S_OK {
				continue
			}
			if spring == S_UNKNOWN {
				if match != matchLimit {
					// explore this path as a potential new start
					glog.V(1).Infof("Exploring %d as BAD", linkPos)
					next := m.Springs.NewWith(linkPos, S_BAD)
					rv += next.findArrangements(linkPos, match, -1, spec, matchLimit)
				}
			} else if spring == S_BAD && match != matchLimit {
				// must be the start of the next pattern
				rv += m.Springs.findArrangements(linkPos, match, -1, spec, matchLimit)
				break
			} else {
				glog.V(1).Infof("found cycle %d/%d prematurely ending at %d/%d", match, matchLimit, m.Length, len(*sl))
				continue MATCHES
			}
		}
		// hit the end without finding another BAD spring, so this is a good match.
		if match == matchLimit {
			rv++
		}
	}

	return
}

// findArrangements iteratively checks along the list of springs, recursing only when an unknown spring is found following an already valid
// pattern.
/*
func (sl *SpringList) oldfindArrangements(from, sI, badTil int, spec Ints, unfoldFactor int) int {
	glog.Infof(" %s%s from %d (of %d)", "", sl, from, len(*sl))
	for n := from; n < len(*sl); n++ {
		s := (*sl)[n]
		if s == S_UNKNOWN {
			a := sl.NewWith(n, S_OK)
			b := sl.NewWith(n, S_BAD)
			return a.findArrangements(n, sI, badTil, spec, unfoldFactor) + b.findArrangements(n, sI, badTil, spec, unfoldFactor)
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
			if sI >= (len(spec) * unfoldFactor) {
				glog.V(2).Infof("bad spring at %d makes %d bad runs, when expecting %d in %s, %s", n, sI, len(spec), sl, spec)
				return 0 // too many bad springs
			}
			eSi := sI
			if unfoldFactor > 1 && sI >= len(spec) {
				eSi = sI % len(spec)
				if eSi == 0 {
					linkSize := sl.lastBad(n) - 1 // 1 this spring, 1 the last bad spring
					if linkSize == 0 {
						glog.Fatalf("End of bad spring cycle in %s at %d with no previous bad spring!", sl, n)
					}
					cycLen := n - linkSize
					cycN := (len(*sl) / cycLen) - 1 // -1 because first cycle is already up to n; only want remaining.
					glog.Infof("Evaluating potential %d spring cycle (x%d) from %d-%d + %d links: %s", cycLen, cycN, 0, cycLen, linkSize, (*sl)[0:n])
					cycI := 0
					linkArrangements := 1
					for {
						cycS := n + (cycI * (cycLen + linkSize))
						if cycS >= len(*sl) {
							break
						}
						if !sl.canMatchTil(spec, cycS) {
							return 0
						}
						for linkN := 0; linkN < linkSize; linkN++ {
							linkPos := cycS + cycLen + linkN
							if linkPos >= len(*sl) {
								break
							}
							// non final cycles must be followed by an OK spring to "link" to the start of the next cycle of bad springs.
							link := (*sl)[linkPos]
							if link == S_BAD {
								glog.Infof("Expected link at %d (%s) was BAD; from cycI=%d, cycS=%d", linkPos, link, cycI, cycS)
								return 0
							}
							if link == S_UNKNOWN {
								glog.Infof("Recursing in cycle at %d", linkPos)
								next := sl.NewWith(linkPos, S_BAD)
								linkArrangements += next.findArrangements(linkPos, 0, -1, spec, unfoldFactor-(sI/len(spec)))
							}
						}
						cycI++
					}
					glog.Infof("Cycle is viable!")
					return 1
				}
			}

			// first bad spring in a run
			badTil = n + spec[eSi] - 1 // s counts as the first item in the run
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
*/

type SpringRow struct {
	Springs SpringList
	BadRuns Ints
}

// Returns a new SpringRow based on this row, but unfolded!
func (sr SpringRow) Unfold() (nr SpringRow) {
	nr.BadRuns = append(nr.BadRuns, sr.BadRuns...)
	for n := 0; n < 5; n++ {
		nr.Springs = append(nr.Springs, sr.Springs...)
		//nr.BadRuns = append(nr.BadRuns, sr.BadRuns...)
		if n != 4 {
			nr.Springs = append(nr.Springs, S_UNKNOWN)
		}
	}
	return
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

func (r *SpringRows) Unfold() {
	for n, sr := range *r {
		(*r)[n] = sr.Unfold()
	}
}

func (r *SpringRows) Arrangements(row int, unfoldFactor int) int {
	if row < 0 || row >= len(*r) {
		glog.Fatalf("invalid row index: %d", row)
	}
	sr := (*r)[row]
	return sr.Springs.Arrangements(sr.BadRuns, unfoldFactor)
}

func (r *SpringRows) SumArrangements(unfoldFactor int) (rv int) {
	for n, sr := range *r {
		glog.V(1).Infof("Starting row %d: %s (len=%d), %s", n, sr, len(sr.Springs), sr.BadRuns)
		rv += r.Arrangements(n, unfoldFactor)
	}
	return
}

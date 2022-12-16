// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 15, Puzzle 2.
// Beacon Exclusion Zone - where is it?.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

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
func Max(a, b int) int {
	if a == -1 {
		return b
	}
	if b == -1 {
		return a
	}
	if a > b {
		return a
	}
	return b
}

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

type Pos struct {
	row, col int
}

func (p Pos) String() string {
	return fmt.Sprintf("x=%d,y=%d", p.col, p.row)
}

func (p Pos) Dist(o Pos) int {
	return Abs(p.col-o.col) + Abs(p.row-o.row)
}

type Sensor struct {
	at      Pos
	closest Pos

	scope int
}

func (s Sensor) String() string {
	return fmt.Sprintf("Sensor at x=%d, y=%d\t: scope % 8d, from beacon at x=%d, y=%d", s.at.col, s.at.row, s.scope, s.closest.col, s.closest.row)
}

func (s Sensor) BBox(limit int) string {
	return fmt.Sprintf("Sensor at x=% 8d, y=% 8d : x=% 8d, y=% 8d - x=% 8d, y=% 8d",
		s.at.col, s.at.row, Max(s.at.col-s.scope, 0), Max(s.at.row-s.scope, 0), Min(s.at.col+s.scope, limit), Min(s.at.row+s.scope, limit))
}

func (s Sensor) BoundFor(row int, limit int) (int, int) {
	offset := Abs(s.at.row - row)
	if offset > s.scope {
		return -1, -1
	}
	return Max(s.at.col-(s.scope-offset), 0), Min(Max(s.at.col+(s.scope-offset), 0), limit)
}

func NewSensor(at, closest Pos) Sensor {
	return Sensor{at: at, closest: closest, scope: at.Dist(closest)}
}

var INPUT_RE = regexp.MustCompile(`Sensor at x=([-\d]+), y=([-\d]+): closest beacon is at x=([-\d]+), y=([-\d]+)`)

func timer(start time.Time, what string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", what, elapsed)
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	sensors := []Sensor{}

	for s.Scan() {
		m := INPUT_RE.FindStringSubmatch(s.Text())
		if len(m) != 5 {
			fmt.Println(len(m), m)
			log.Fatalf("Couldn't parse: %s", s.Text())
		}
		s := NewSensor(Pos{col: Int(m[1]), row: Int(m[2])}, Pos{col: Int(m[3]), row: Int(m[4])})
		sensors = append(sensors, s)
		fmt.Println(s)
	}

	// Take the limit from argv so we can easily test sample vs input
	limit := Int(os.Args[1])
	for _, s := range sensors {
		fmt.Println(s.BBox(limit))
	}

	// Take an optional second parameter to only scan 1 row (easier for debugging)
	onlyRow := 0
	if len(os.Args) > 2 {
		onlyRow = Int(os.Args[2])
	}

	missing := [][]int{}
	{
		defer timer(time.Now(), "loop")
		// Iterate over all rows.
ROWS:
		for i := onlyRow; i <= limit; i++ {
			// Build a list of min,max coverage from each sensor over this row.
			covers := [][]int{}
			for _, s := range sensors {
				min, max := s.BoundFor(i, limit)
				if min == -1 || max == -1 {
					continue // doesn't cover this row at all
				}
				covers = append(covers, []int{min, max})
			}
			// Sort that list of coverage, and check it covers from 0-limit
			sort.Slice(covers, func(a, b int) bool { return Min(covers[a][0], covers[a][1]) < Min(covers[b][0], covers[b][1]) })
			if covers[0][0] != 0 {
				missing = append(missing, []int{i, 0}) // missing the start
			}
			expect := covers[0][1] + 1
			for _, c := range covers[1:] {
				if c[0] > expect {
					missing = append(missing, []int{i, expect}) // missing a point in the middle
					break ROWS
				}
				expect = Max(Min(c[1]+1, limit), expect)
				if expect == limit {
					break
				}
			}
			if expect != limit {
				missing = append(missing, []int{i, limit}) // missing the end
			}

			if onlyRow != 0 {
				break
			}
		}
	}

	if len(missing) == 1 {
		fmt.Println(missing[0])
		fmt.Println("Tuning frequency: ", (missing[0][1]*4000000)+missing[0][0])
	} else if len(missing) > 1 {
		fmt.Println("Multiple missing bits!", missing)
	}
}

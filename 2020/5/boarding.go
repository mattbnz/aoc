// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 5, Puzzle 1
// Converting boarding code to seat id.

package main

import (
    "bufio"
    "fmt"
	 "os"
)

func main() {
	 maxseat := 0
    // Loop over scanned codes
    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
		  start := 0
		  end := 127
		  row := -1
		  for _, char := range s.Text()[:7] {
				diff := end-start
				if (diff > 1) {
					 if char == 'F' {
						  end = start + (diff/2)
					 } else if char == 'B' {
						  start = end - (diff/2)
					 }
				} else {
					 if char == 'F' {
						  row = start
					 } else {
						  row = end
					 }
				}
		  }
		  start = 0
		  end = 7
		  col := -1
		  for _, char := range s.Text()[7:] {
				diff := end-start
				if (diff > 1) {
					 if char == 'L' {
						  end = start + (diff/2)
					 } else if char == 'R' {
						  start = end - (diff/2)
					 }
				} else {
					 if char == 'L' {
						  col = start
					 } else {
						  col = end
					 }
				}
		  }
		  seatid := (row * 8) + col
		  if seatid > maxseat {
				maxseat = seatid
		  }
    }

    fmt.Println(maxseat);
}

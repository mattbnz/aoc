// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 5, Puzzle 2
// Find missing seat ID. 

package main

import (
    "bufio"
    "fmt"
	 "os"
	 "sort"
)

func main() {
	 seats := []int{}
    // Loop over scanned codes, add to list.
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
		  seats = append(seats, seatid)
    }

	 // Sort and then find missing seat. 
	 sort.Ints(seats)
	 last := seats[0]
	 for _, seatid := range seats[1:] {
		  if seatid - 1 != last {
				fmt.Println(seatid -1)
		  }
		  last = seatid
	 }
}

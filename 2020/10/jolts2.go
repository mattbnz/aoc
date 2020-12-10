// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 10, Puzzle 2
// Jolt combinations.

package main

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
    "sort"
    "strconv"
)

func Find(list []int, jolt int) []int {
    rv := []int{}
    for idx, j := range list {
        diff := 0
        if jolt == 0 {
            diff = int(math.Abs(float64(jolt - j)))
        } else {
            diff = j - jolt
        }
        if diff >= 0 && diff <= 3 {
            //fmt.Printf("Found %d at %d for %d\n", j, idx, jolt)
            rv = append(rv, idx)
        }
    }
    return rv
}

func Next(list []int, path []int, start int) int {
    count := 0
    c := Find(list, start)
    if len(c) == 0 {
        //fmt.Printf("Found path %v\n", path)
        return 1
    }
    for _, idx := range c {
        newlist := make([]int, len(list))
        copy(newlist, list)
        newlist = append(newlist[:idx], newlist[idx+1:]...)
        newpath := []int{}
        newpath = append(newpath, path...)
        newpath = append(newpath, list[idx])
        count += Next(newlist, newpath, list[idx])
    }
    //fmt.Printf("Next list=%d, start=%d; count=%d\n", len(list), start, count)
    return count
}

func main() {
    s := bufio.NewScanner(os.Stdin)
    adaptors := []int{}
    for s.Scan() {
        n, err := strconv.Atoi(s.Text())
        if err != nil {
            log.Fatal("Cannot parse: ", s.Text())
        }
        adaptors = append(adaptors, n)
    }
    // Work backwards from the largest adaptor, and calculate how many paths
    // forward from each one there are.
    paths := make([]int, len(adaptors))
    paths[len(adaptors)-1] = 1  // Only one path from the biggest. 
    sort.Ints(adaptors)
    fmt.Printf("%v\n", adaptors)
    for i := len(adaptors) -2; i>=0; i-- {
        sum := 0
        for diff := 1; diff <= 3; diff++ {
            try := i+diff
            if try < len(adaptors) {
                if adaptors[try] - adaptors[i] <= 3 {
                    sum += paths[try]
                }
            }
        }
        paths[i] = sum
    }
    fmt.Printf("%v\n", paths)

    // Then work out which of the first 3 positions are valid, and sum them.
    sum := 0
    for i := 0; i<3; i++ {
        if adaptors[i] <= 3 {
            sum += paths[i]
        }
    }
    fmt.Println(sum)
}

// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 7, Puzzle 1
// Bags in bags... 

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "regexp"
    "strconv"
    "strings"
)

type BagMap map[string][]ContainedBag
type ContainedBag struct {
    color string
    num int
}
var BagRe = regexp.MustCompile(`(\d)+\s(.*)`)
var StripRe = regexp.MustCompile(` +bags?.?$`)

func StripName(input string) string {
    return StripRe.ReplaceAllString(input, "")
}

func FindBagsFor(want string, rules BagMap) []string {
    rv := []string{}
    for bag, contains := range rules {
        for _, inner := range contains {
            if inner.color == want {
                rv = append(rv, bag)
                break
            }
        }
    }
    return rv
}

func main() {
    // Read in the rules.
    allbags := make(BagMap)
    s := bufio.NewScanner(os.Stdin)
    for s.Scan() {
        parts := strings.Split(s.Text(), " contain ")
        contents := []string{}
        if parts[1] != "no other bags." {
            contents = strings.Split(parts[1], ", ")
        }
        bags := []ContainedBag{}
        for _, bag := range contents {
            match := BagRe.FindStringSubmatch(strings.TrimRight(bag, "."))
            q, err := strconv.Atoi(match[1])
            if err != nil {
                log.Fatal("Cannot parse: ", s.Text())
            }
            bags = append(bags, ContainedBag{color: StripName(match[2]), num: q})
        }
        allbags[StripName(parts[0])] = bags
    }
    
    // Loop and find out what can be in what... 
    containers := make(map[string]int)
    bag := ""
    next := []string{"shiny gold"}
    for len(next) > 0 {
        bag, next = next[len(next)-1], next[:len(next)-1]
        t := FindBagsFor(bag, allbags)
        next = append(next, t...)
        for _, o := range t {
            containers[o] = 1
        }
    }
    fmt.Println(len(containers))
}

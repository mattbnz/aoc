// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 7, Puzzle 2
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

func CountBagsIn(want string, rules BagMap) int {
    sum := 1
    for _, inner := range rules[want] {
        sum += inner.num * CountBagsIn(inner.color, rules)
    }
    return sum
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

    // Loop and find out what's in a shiny gold bag
    fmt.Println(CountBagsIn("shiny gold", allbags)-1)
}

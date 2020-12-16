// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 16, Puzzle 2
// Ticket Validity

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"
)

func Remove(in []int, pop int) []int {
    rv := []int{}
    for _, v := range in {
        if v != pop {
            rv = append(rv, v)
        }
    }
    return rv
}

func ValidForField(rules [][]int, v int) bool {
    for _, rule := range rules {
        if v >= rule[0] && v <= rule[1] {
            return true
        }
    }
    return false
}

func main() {
    valid := make(map[int]bool)
    fields := make(map[string][][]int)

    s := bufio.NewScanner(os.Stdin)
    // Read valid values for each field.
    for s.Scan() {
        if strings.TrimSpace(s.Text()) == "" {
            break
        }
        split := strings.Split(strings.TrimSpace(s.Text()), ":")
        field := strings.TrimSpace(split[0])
        fields[field] = [][]int{}

        pairs := strings.Split(split[1], " or ")
        for _, pair := range pairs {
            minmax := strings.Split(pair, "-")
            min, err := strconv.Atoi(strings.TrimSpace(minmax[0]))
            if err != nil {
                log.Fatal("Bad min for ", minmax[0], err)
            }
            max, err := strconv.Atoi(strings.TrimSpace(minmax[1]))
            if err != nil {
                log.Fatal("Bad max for ", minmax[1], err)
            }
            for i:=min; i<=max; i++ {
                valid[i] = true
            }
            fields[field] = append(fields[field], []int{min, max})
        }
    }
    //fmt.Println(fields)

    // Read my ticket for later (skip heading line first)
    s.Scan()
    s.Scan()
    myticket := s.Text()
    //fmt.Println(myticket)

    // Skip 2 lines (blank, nearby tickets header)
    s.Scan()
    s.Scan()

    // Accumulate valid nearby tickets
    nearby := [][]int{}
    for s.Scan() {
        values := strings.Split(strings.TrimSpace(s.Text()), ",")
        positions := []int{}
        good := true
        for _, vstr := range values {
            v, err := strconv.Atoi(vstr)
            if err != nil {
                log.Fatal("Bad ticket value for ", s.Scan())
            }
            _, found := valid[v]
            if !found {
                good = false
                break
            }
            positions = append(positions, v)
        }
        if good {
            nearby = append(nearby, positions)
        }
    }
    //fmt.Println(nearby)

    // Start by assuming any position is OK for any field
    allpos := []int{}
    for i:=0; i<len(nearby[0]); i++ {
        allpos = append(allpos, i)
    }
    fieldpos := make(map[string][]int)
    for field, _ := range fields {
        fieldpos[field] = make([]int, len(allpos))
        copy(fieldpos[field], allpos)
    }

    // Loop through nearby tickets and eliminate positions that can't work for
    // each field.
    for field, _ := range fieldpos {
        //fmt.Printf("Checking %s: %v\n", field, fieldpos[field])
        for _, ticket := range nearby {
            //fmt.Printf("- Against %v\n", ticket)
            for pos, v := range ticket {
                if ValidForField(fields[field], v) {
                    continue
                }
                //fmt.Printf("%s can't be in pos %d from %v\n", field, pos, ticket)
                fieldpos[field] = Remove(fieldpos[field], pos)
                if len(fieldpos[field]) == 1 {
                    break
                }
            }
            if len(fieldpos[field]) == 1 {
                break
            }
        }
    }
    for i:=0; i<len(fieldpos); i++ {
        for field, _ := range fieldpos {
            // If we know which position this field is, eliminate that position 
            // from the valid positions of all other fields
            if len(fieldpos[field]) == 1 {
                for field2, _ := range fieldpos {
                    if field2 == field {
                        continue
                    }
                    fieldpos[field2] = Remove(fieldpos[field2], fieldpos[field][0])
                }
            }
        }
    }
    //fmt.Println(fieldpos)

    // Convert myticket to ints
    values := strings.Split(strings.TrimSpace(myticket), ",")
    mine := []int{}
    for _, vstr := range values {
        v, err := strconv.Atoi(vstr)
        if err != nil {
            log.Fatal("Bad myticket value for ", s.Scan())
        }
        mine = append(mine, v)
    }
    //fmt.Println(mine)
    product := 1
    for field, pos := range fieldpos {
        if strings.HasPrefix(field, "departure") {
            product *= mine[pos[0]]
        }
        //fmt.Printf("%s: %d %d\n", field, pos[0], product)
    }
    fmt.Println(product)
}

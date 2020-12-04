// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 4, Puzzle 1
// Check passport field presence for validity.

package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

var PP_FIELDS = []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid", "cid"}
var CID = FindFieldIdx("cid")

// Return the index of the field in the passport, or -1 if not found.
func FindFieldIdx(name string) (int) {
    for idx, field := range PP_FIELDS {
        if strings.ToLower(name) == field {
            return idx
        }
    }
    return -1
}

// Create a new passport structure (slice) all zeroed.
func NewPassport() ([]int) {
    passport := make([]int, len(PP_FIELDS))
    for idx, _ := range passport {
        passport[idx] = 0
    }
    return passport
}

func main() {
    // Loop over passports and parse them all. 
    s := bufio.NewScanner(os.Stdin)
    passports := [][]int{}
    passport := NewPassport()
    for s.Scan() {
        if strings.TrimSpace(s.Text()) == "" {
            passports = append(passports, passport)
            passport = NewPassport()
        }
        fields := strings.Fields(s.Text())

        for _, field := range fields {
            k_v := strings.Split(field, ":")
            idx := FindFieldIdx(k_v[0])
            if idx != -1 {
                passport[idx] = 1
            }
        }
    }
    passports = append(passports, passport)

    // Now check validity of the passports we parsed.
    valid := 0
    for _, passport := range passports {
        all_present := true
        for idx, _ := range PP_FIELDS {
            if idx == CID {
                continue
            }
            if passport[idx] == 0 {
                all_present = false
                break
            }
        }
        if all_present {
            valid++
        }
    }

    fmt.Println(valid);
}

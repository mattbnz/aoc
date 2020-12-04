// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 4, Puzzle 1
// Check passport field presence for validity.

package main

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"
)

type Validator func(string) (bool)

var PP_FIELDS = []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid", "cid"}
var PP_VALIDATORS = []Validator{ValidByr, ValidIyr, ValidEyr, ValidHgt, ValidHcl,
                               ValidEcl, ValidPid, nil}
var ECL_VALUES = []string{"amb", "blu", "brn", "gry", "grn", "hzl", "oth"}
var CID = FindFieldIdx("cid", PP_FIELDS)

// Validator Helpers
func GetYear(value string) (int, error) {
    if len(value) != 4 {
        return -1, errors.New("year must be 4 characters")
    }
    return strconv.Atoi(value)
}
func ValidYear(value string, min, max int) (bool) {
    year, valid := GetYear(value)
    if valid != nil {
        return false
    }
    if year >= min && year <= max {
        return true
    }
    return false
}
// Validators
func ValidByr(value string) (bool) {
    return ValidYear(value, 1920, 2002)
}
func ValidIyr(value string) (bool) {
    return ValidYear(value, 2010, 2020)
}
func ValidEyr(value string) (bool) {
    return ValidYear(value, 2020, 2030)
}
func ValidHgt(value string) (bool) {
    suffix := value[len(value)-2:]
    intVal, err := strconv.Atoi(value[:len(value)-2])
    if err != nil {
        return false
    }
    if suffix == "cm" {
        if intVal >= 150 && intVal <= 193 {
            return true
        }
    } else if suffix == "in" {
        if intVal >= 59 && intVal <= 76 {
            return true
        }
    }
    return false;
}
func ValidHcl(value string) (bool) {
    match, err := regexp.MatchString(`^#[0-9a-f]{6}$`, value)
    if err == nil && match {
        return true
    }
    return false
}
func ValidEcl(value string) (bool) {
    idx := FindFieldIdx(strings.ToLower(value), ECL_VALUES)
    if idx != -1 {
        return true
    }
    return false
}
func ValidPid(value string) (bool) {
    match, err := regexp.MatchString(`^[0-9]{9}$`, value)
    if err == nil && match {
        return true
    }
    return false
}

// Return the index of the field in the list, or -1 if not found.
func FindFieldIdx(name string, list []string) (int) {
    for idx, field := range list {
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
            idx := FindFieldIdx(k_v[0], PP_FIELDS)
            if idx != -1 {
                validator := PP_VALIDATORS[idx]
                if validator != nil && validator(k_v[1]) {
                    passport[idx] = 1
                }
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

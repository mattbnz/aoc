package day5

import (
	"log"
	"strconv"
	"strings"
)

func numberList(list string) (rv []int) {
	for _, nStr := range strings.Split(strings.TrimSpace(list), " ") {
		nStr = strings.TrimSpace(nStr)
		if nStr == "" {
			continue
		}
		i, err := strconv.Atoi(nStr)
		if err != nil {
			log.Fatalf("Bad number '%s' (from %s): %v", nStr, list, err)
		}
		rv = append(rv, i)
	}
	return
}

func Max(a int, b ...int) (rv int) {
	rv = a
	for _, t := range b {
		if t > rv {
			rv = t
		}
	}
	return
}

func Min(a int, b ...int) (rv int) {
	rv = a
	for _, t := range b {
		if t < rv {
			rv = t
		}
	}
	return
}

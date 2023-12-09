package day9

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

func Sum(l []int) (rv int) {
	for _, n := range l {
		rv += n
	}
	return
}

package day13

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Ints []int

func (c Ints) String() string {
	s := []string{}
	for _, i := range c {
		s = append(s, fmt.Sprintf("%d", i))
	}
	return strings.Join(s, ",")
}

func NewIntsFromCSV(list string) (rv Ints) {
	for _, nStr := range strings.Split(strings.TrimSpace(list), ",") {
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

func Abs(x int) int {
	if x < 0 {
		return x * -1
	}
	return x
}

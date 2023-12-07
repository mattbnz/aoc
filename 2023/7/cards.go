// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 7.
// Camel Cards

package day7

import (
	"cmp"
	"fmt"
	"sort"
)

var cardValues = map[string]Card{
	"A": 14,
	"K": 13,
	"Q": 12,
	"J": 11,
	"T": 10,
	"9": 9,
	"8": 8,
	"7": 7,
	"6": 6,
	"5": 5,
	"4": 4,
	"3": 3,
	"2": 2,
}

type Card int

func (c Card) String() string {
	for s, v := range cardValues {
		if v == c {
			return s
		}
	}
	return "!"
}

func NewCard(s string) (c Card, err error) {
	c, ok := cardValues[s]
	if !ok {
		err = fmt.Errorf("unknown card symbol: %s", s)
		return
	}
	return
}

type HandType int

const (
	H_Unknown HandType = iota
	H_High
	H_Pair
	H_2Pair
	H_3Kind
	H_FullHouse
	H_4Kind
	H_5Kind

	H_Max // Must be last!
)

func (h HandType) String() string {
	switch h {
	case H_High:
		return "High"
	case H_Pair:
		return "Pair"
	case H_2Pair:
		return "2Pair"
	case H_3Kind:
		return "3Kind"
	case H_FullHouse:
		return "FullHouse"
	case H_4Kind:
		return "4Kind"
	case H_5Kind:
		return "5Kind"
	}
	return "Unknown"
}

type Hand struct {
	Cards []Card

	value HandType
}

func NewHand(s string) (h Hand, err error) {
	for _, r := range s {
		var c Card
		c, err = NewCard(string(r))
		if err != nil {
			return
		}
		h.Cards = append(h.Cards, c)
	}
	if len(h.Cards) != 5 {
		err = fmt.Errorf("bad hand size")
		return
	}
	h.recognize()
	return
}

func (h Hand) String() string {
	return fmt.Sprintf("%s (%s%s%s%s%s)", h.value, h.Cards[0], h.Cards[1], h.Cards[2], h.Cards[3], h.Cards[4])
}

func (h *Hand) recognize() {
	kMap := map[Card]int{}
	for _, c := range h.Cards {
		kMap[c]++
	}
	syms := []Card{}
	for s := range kMap {
		syms = append(syms, s)
	}
	sort.Slice(syms, func(i, j int) bool {
		// Not Less; to sort to largest first.
		return kMap[syms[i]] > kMap[syms[j]]
	})
	pairs := 0
	for _, sym := range syms {
		if kMap[sym] == 5 {
			h.value = H_5Kind
			return
		}
		if kMap[sym] == 4 {
			h.value = H_4Kind
			return
		}
		if kMap[sym] == 3 {
			if len(kMap) == 2 && kMap[syms[1]] == 2 {
				h.value = H_FullHouse
			} else {
				h.value = H_3Kind
			}
			return
		}
		if kMap[sym] == 2 {
			pairs++
		}
	}
	if pairs == 2 {
		h.value = H_2Pair
	} else if pairs == 1 {
		h.value = H_Pair
	} else {
		h.value = H_High
	}
}

// HandSortFunc returns a negative number when a < b, a positive number when a > b and zero when a == b.
func HandSortFunc(a, b Hand) (rv int) {
	rv = cmp.Compare(a.value, b.value)
	if rv != 0 {
		return
	}
	for n := 0; n < 5; n++ {
		rv = cmp.Compare(a.Cards[n], b.Cards[n])
		if rv != 0 {
			return
		}
	}
	return 0
}

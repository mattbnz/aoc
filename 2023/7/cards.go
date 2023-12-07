// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 7.
// Camel Cards

package day7

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

const JokerPlain = "J"
const JokerMarker = "-"
const Joker = JokerPlain + JokerMarker

var cardValues = map[string]Card{
	"A":        14,
	"K":        13,
	"Q":        12,
	JokerPlain: 11,
	"T":        10,
	"9":        9,
	"8":        8,
	"7":        7,
	"6":        6,
	"5":        5,
	"4":        4,
	"3":        3,
	"2":        2,
	Joker:      1,
}

type Card int

func (c Card) String() string {
	for s, v := range cardValues {
		if v == c {
			return strings.TrimSuffix(s, JokerMarker)
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
	Bet   int

	WithJoker bool

	value HandType
}

func NewHand(s string, withJoker bool) (h Hand, err error) {
	cards, bet, ok := strings.Cut(s, " ")
	if !ok {
		err = fmt.Errorf("invalid format: %s", s)
		return
	}
	h.Bet, err = strconv.Atoi(bet)
	if err != nil {
		err = fmt.Errorf("bad bet value (%s): %w", bet, err)
		return
	}
	h.WithJoker = withJoker
	for _, r := range cards {
		cS := string(r)
		if withJoker && cS == "J" {
			cS += JokerMarker
		}
		var c Card
		c, err = NewCard(cS)
		if err != nil {
			return
		}
		h.Cards = append(h.Cards, c)
	}
	if len(h.Cards) != 5 {
		err = fmt.Errorf("bad hand size")
		return
	}
	err = h.recognize()
	return
}

func (h Hand) String() string {
	return fmt.Sprintf("%s (%s)", h.value, h.CardString())
}

func (h Hand) CardString() string {
	return fmt.Sprintf("%s%s%s%s%s", h.Cards[0], h.Cards[1], h.Cards[2], h.Cards[3], h.Cards[4])
}

func (h *Hand) bestJokerHand(jh Hand, level int) (ht HandType, err error) {
	best := h.value
	cs := jh.CardString()
	glog.V(1).Infof("%sbestJokerHand for %s", strings.Repeat("", 0), jh)
	for n, c := range jh.Cards {
		if c != cardValues[Joker] {
			continue
		}
		for r := cardValues["2"]; r <= cardValues["A"]; r++ {
			if r == cardValues[JokerPlain] {
				continue
			}
			var nh Hand
			nh, err = NewHand(cs[:n]+r.String()+cs[n+1:]+" 0", true)
			if err != nil {
				return
			}
			if nh.value > best {
				best = nh.value
				glog.V(1).Infof("%s - %s is new best hand from %s", strings.Repeat("", 0), nh, best)
			}
		}
	}
	return best, nil
}

func (h *Hand) recognize() error {
	kMap := map[Card]int{}
	for _, c := range h.Cards {
		kMap[c]++
	}
	h.value = h.recognizeHand(kMap)
	if h.WithJoker && kMap[1] != 0 {
		if kMap[1] == 5 {
			return nil // Can't improve on 5 of a kind!
		}
		v, err := h.bestJokerHand(*h, 0)
		if err != nil {
			return err
		}
		h.value = v
	}
	return nil
}

func (h *Hand) recognizeHand(kMap map[Card]int) HandType {
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
			return H_5Kind
		}
		if kMap[sym] == 4 {
			return H_4Kind
		}
		if kMap[sym] == 3 {
			if len(kMap) == 2 && kMap[syms[1]] == 2 {
				return H_FullHouse
			} else {
				return H_3Kind
			}
		}
		if kMap[sym] == 2 {
			pairs++
		}
	}
	if pairs == 2 {
		return H_2Pair
	} else if pairs == 1 {
		return H_Pair
	} else {
		return H_High
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

type Hands []Hand

func NewHands(filename string, withJoker bool) (hl Hands, er error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	lineno := 0
	for s.Scan() {
		h, err := NewHand(s.Text(), withJoker)
		if err != nil {
			err = fmt.Errorf("bad hand on line %d: %w", lineno, err)
			return
		}
		lineno++

		hl = append(hl, h)
	}
	glog.Infof("Read %d hands", len(hl))
	return
}

func (h Hands) Winnings() (rv int) {
	slices.SortFunc(h, HandSortFunc)
	for rank, hand := range h {
		winnings := hand.Bet * (rank + 1)
		glog.V(1).Infof("Rank % 4d: %s bets % 4d and wins % 6d", rank, hand, hand.Bet, winnings)
		rv += winnings
	}
	return
}

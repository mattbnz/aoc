// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 4, Puzzle 2
// Scratchcards

package main

import (
	"bufio"
	"log"
	"os"
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

func PlayCards(cards []*Card, card int) {
	cards[card].NumberWon++
	if cards[card].Score == 0 {
		return
	}
	for i := 1; i <= cards[card].Score; i++ {
		PlayCards(cards, card+i)
	}
}

type Card struct {
	NumberWon int // how many we've won
	Score     int // the number of matching numbers on the card.
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	cards := []*Card{}
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		wStr, nStr, found := strings.Cut(s.Text(), "|")
		if !found {
			log.Fatalf("Bad game line: %s", s.Text())
		}
		_, wStr2, found := strings.Cut(wStr, ":")
		if !found {
			log.Fatalf("Bad game line: %s", s.Text())
		}
		winners := numberList(wStr2)
		numbers := numberList(nStr)
		matches := 0
	N:
		for _, n := range numbers {
			for _, wN := range winners {
				if n == wN {
					matches++
					continue N
				}
			}
		}
		cards = append(cards, &Card{NumberWon: 0, Score: matches})
	}
	for n, _ := range cards {
		PlayCards(cards, n)
		for n, c := range cards {
			log.Printf("Card %d: %d won", n, c.NumberWon)
		}
	}
	sum := 0
	for n, c := range cards {
		log.Printf("Card %d: %d won", n, c.NumberWon)
		sum += c.NumberWon
	}

	log.Printf("%d", sum)
}

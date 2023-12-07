package day7

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func MustHand(s string) Hand {
	h, err := NewHand(s + " 0")
	if err != nil {
		panic(err)
	}
	return h
}

var canonicalHands = map[HandType]Hand{
	H_High:      MustHand("23456"),
	H_Pair:      MustHand("22456"),
	H_2Pair:     MustHand("22466"),
	H_3Kind:     MustHand("22256"),
	H_FullHouse: MustHand("22266"),
	H_4Kind:     MustHand("22226"),
	H_5Kind:     MustHand("22222"),
}
var canonicalHandOrder = []HandType{H_High, H_Pair, H_2Pair, H_3Kind, H_FullHouse, H_4Kind, H_5Kind}

func Test_Hand(t *testing.T) {
	for hType, h := range canonicalHands {
		assert.Equal(t, hType.String(), h.value.String())
		assert.Equal(t, 0, HandSortFunc(h, h), "hand not equal to itself!")
	}
}

func Test_HandCmp(t *testing.T) {
	last := canonicalHands[canonicalHandOrder[0]]
	for _, ht := range canonicalHandOrder[1:] {
		h := canonicalHands[ht]
		assert.Equal(t, -1, HandSortFunc(last, h), "%s was not less than %s", h, last)
		last = h
	}
}

// in increasing order
var testHands = []Hand{
	MustHand("23456"),
	MustHand("32456"), // 3 beats 2 in first place
	MustHand("65432"), // 6 beats 3

	MustHand("22345"), // pair beets singles
	MustHand("33456"), // first 3 beats 2
	MustHand("34AA5"), // second 4 beats 3

	MustHand("25566"), // 2 pair beats pair
	MustHand("55466"), // 5 beats 2

	MustHand("23555"), // 3Kind beats 2 pair
	MustHand("55523"), // first 5 beats 2

	MustHand("22255"), // Full house beats 3Kind
	MustHand("33222"), // First 3 beats 2

	MustHand("22223"), // 4Kind beats Full House
	MustHand("22224"), // 4 beats 3
	MustHand("A2222"), // A beats 2

	MustHand("22222"), // 5Kind beats 4Kind
	MustHand("33333"), // 3 beats 2
}

func Test_HandCardCmp(t *testing.T) {
	last := testHands[0]
	for _, h := range testHands[1:] {
		assert.Equal(t, -1, HandSortFunc(last, h), "%s was not less than %s", h, last)
		last = h
	}
}

func Test_Sample(t *testing.T) {
	hands, err := NewHands("sample")
	require.NoError(t, err)
	assert.Equal(t, 6440, hands.Winnings())
}

func Test_Part1(t *testing.T) {
	hands, err := NewHands("input")
	require.NoError(t, err)
	log.Printf("Winnings are: %d", hands.Winnings())
}

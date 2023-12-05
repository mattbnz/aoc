package day5

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Override(t *testing.T) {
	o := Override{SourceBase: 98, DestBase: 50, Count: 2}
	assert.Equal(t, -1, o.Contains(97))
	assert.Equal(t, -1, o.Contains(100))
	assert.Equal(t, 0, o.Contains(98))
	assert.Equal(t, 1, o.Contains(99))
}

func Test_Map(t *testing.T) {
	//
	o := &Override{SourceBase: 98, DestBase: 50, Count: 2}
	m := Mapping{
		Overrides: []*Override{o},
	}
	assert.Equal(t, 97, m.Map(97))
	assert.Equal(t, 50, m.Map(98))
	assert.Equal(t, 51, m.Map(99))
	assert.Equal(t, 100, m.Map(100))
}

func Test_Sample(t *testing.T) {
	almanac, err := NewAlmanac("sample")
	require.NoError(t, err)
	for seed, soil := range map[int]int{
		0:  0,
		1:  1,
		48: 48,
		49: 49,
		50: 52,
		51: 53,
		96: 98,
		97: 99,
		98: 50,
		99: 51,
	} {
		assert.Equal(t, soil, almanac.Lookup("seed", seed, "soil"))
	}
	for seed, location := range map[int]int{
		79: 82,
		14: 43,
		55: 86,
		13: 35,
	} {
		assert.Equal(t, location, almanac.Lookup("seed", seed, "location"))
	}
	assert.Equal(t, 35, almanac.BestLocation())
}

func Test_Input(t *testing.T) {
	almanac, err := NewAlmanac("input")
	require.NoError(t, err)
	log.Printf("Best Location is: %d", almanac.BestLocation())
}

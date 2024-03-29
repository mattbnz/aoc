package day5

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Override(t *testing.T) {
	//
	o := &Override{SourceBase: 98, DestBase: 50, Count: 2}

	assert.Equal(t, -1, o.ToDest(97))
	assert.Equal(t, 50, o.ToDest(98))
	assert.Equal(t, 51, o.ToDest(99))
	assert.Equal(t, -1, o.ToDest(100))

	assert.Equal(t, -1, o.ToSource(49))
	assert.Equal(t, 98, o.ToSource(50))
	assert.Equal(t, 99, o.ToSource(51))
	assert.Equal(t, -1, o.ToSource(52))
}

func Test_Map(t *testing.T) {
	//
	o := &Override{SourceBase: 98, DestBase: 50, Count: 2}
	m := Mapping{
		Overrides: []*Override{o},
	}
	v, _ := m.ToDest(97, 100)
	assert.Equal(t, 97, v)
	v, _ = m.ToDest(98, 100)
	assert.Equal(t, 50, v)
	v, _ = m.ToDest(99, 100)
	assert.Equal(t, 51, v)
	v, _ = m.ToDest(100, 100)
	assert.Equal(t, 100, v)

	v, _ = m.ToSource(97, 100)
	assert.Equal(t, 97, v)
	v, _ = m.ToSource(50, 100)
	assert.Equal(t, 98, v)
	v, _ = m.ToSource(51, 100)
	assert.Equal(t, 99, v)
	v, _ = m.ToSource(100, 100)
	assert.Equal(t, 100, v)
}

func Test_MapBounds(t *testing.T) {
	almanac, err := NewAlmanac("sample")
	require.NoError(t, err)
	m := almanac.getMap("seed")
	v, b := m.ToDest(0, 100)
	assert.Equal(t, 0, v)
	assert.Equal(t, Override{SourceBase: 0, DestBase: 0, Count: 50}, b)
	v, b = m.ToDest(49, 100)
	assert.Equal(t, 49, v)
	assert.Equal(t, Override{SourceBase: 0, DestBase: 0, Count: 50}, b)
	v, b = m.ToDest(50, 100)
	assert.Equal(t, 52, v)
	assert.Equal(t, Override{SourceBase: 50, DestBase: 52, Count: 48}, b)
	v, b = m.ToDest(97, 100)
	assert.Equal(t, 99, v)
	assert.Equal(t, Override{SourceBase: 50, DestBase: 52, Count: 48}, b)
	v, b = m.ToDest(98, 100)
	assert.Equal(t, 50, v)
	assert.Equal(t, Override{SourceBase: 98, DestBase: 50, Count: 2}, b)
	v, b = m.ToDest(99, 100)
	assert.Equal(t, 51, v)
	assert.Equal(t, Override{SourceBase: 98, DestBase: 50, Count: 2}, b)
	v, b = m.ToDest(100, 100)
	assert.Equal(t, 100, v)
	assert.Equal(t, Override{SourceBase: 100, DestBase: 100, Count: 0}, b)
}

func Test_BoundedLookup(t *testing.T) {
	almanac, err := NewAlmanac("sample")
	require.NoError(t, err)

	// https://docs.google.com/spreadsheets/d/1NZkrEvOFNTz-vwKHRL2aP34UEeK0YWwgIurKgChdnIA/edit#gid=1157275920
	d, b := almanac.BoundedLookup("seed", 79, "soil")
	assert.Equal(t, 81, d)
	assert.Equal(t, Override{SourceBase: 50, DestBase: 52, Count: 48}, b)

	d, b = almanac.BoundedLookup("seed", 79, "fertilizer")
	assert.Equal(t, 81, d)
	assert.Equal(t, Override{SourceBase: 52, DestBase: 54, Count: 46}, b)

	d, b = almanac.BoundedLookup("seed", 79, "water")
	assert.Equal(t, 81, d)
	assert.Equal(t, Override{SourceBase: 59, DestBase: 61, Count: 39}, b)

	d, b = almanac.BoundedLookup("seed", 79, "light")
	assert.Equal(t, 74, d)
	assert.Equal(t, Override{SourceBase: 59, DestBase: 61, Count: 34}, b)

	d, b = almanac.BoundedLookup("seed", 79, "temperature")
	assert.Equal(t, 78, d)
	assert.Equal(t, Override{SourceBase: 69, DestBase: 71, Count: 13}, b)

	d, b = almanac.BoundedLookup("seed", 79, "humidity")
	assert.Equal(t, 78, d)
	assert.Equal(t, Override{SourceBase: 71, DestBase: 73, Count: 11}, b)

	d, b = almanac.BoundedLookup("seed", 79, "location")
	assert.Equal(t, 82, d)
	assert.Equal(t, Override{SourceBase: 71, DestBase: 73, Count: 11}, b)
}

func Test_LookupSource(t *testing.T) {
	almanac, err := NewAlmanac("sample")
	require.NoError(t, err)

	// https://docs.google.com/spreadsheets/d/1NZkrEvOFNTz-vwKHRL2aP34UEeK0YWwgIurKgChdnIA/edit#gid=1157275920
	d, b := almanac.LookupSource("humidity", 0, "location")
	assert.Equal(t, 0, d)
	assert.Equal(t, Override{SourceBase: 0, DestBase: 0, Count: 56}, b)

	d, b = almanac.LookupSource("temperature", 0, "location")
	assert.Equal(t, 69, d)
	assert.Equal(t, Override{SourceBase: 0, DestBase: 0, Count: 1}, b)

	d, b = almanac.LookupSource("light", 0, "location")
	assert.Equal(t, 65, d)
	assert.Equal(t, Override{SourceBase: 0, DestBase: 0, Count: 1}, b)

	d, b = almanac.LookupSource("water", 0, "location")
	assert.Equal(t, 72, d)
	assert.Equal(t, Override{SourceBase: 0, DestBase: 0, Count: 1}, b)

	d, b = almanac.LookupSource("fertilizer", 0, "location")
	assert.Equal(t, 72, d)
	assert.Equal(t, Override{SourceBase: 0, DestBase: 0, Count: 1}, b)

	d, b = almanac.LookupSource("soil", 0, "location")
	assert.Equal(t, 72, d)
	assert.Equal(t, Override{SourceBase: 0, DestBase: 0, Count: 1}, b)

	d, b = almanac.LookupSource("seed", 0, "location")
	assert.Equal(t, 70, d)
	assert.Equal(t, Override{SourceBase: 0, DestBase: 0, Count: 1}, b)
}

func Test_Sample(t *testing.T) {
	almanac, err := NewAlmanac("sample")
	require.NoError(t, err)
	assert.Equal(t, 100, almanac.Max)
	assert.Equal(t, 35, almanac.BestLocation())
	assert.Equal(t, 46, almanac.BestLocation2())
	assert.Equal(t, 46, almanac.BestLocation2b())
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
		l, _ := almanac.BoundedLookup("seed", seed, "soil")
		assert.Equal(t, soil, l, "BoundedLookup wrong for seed %d", seed)
		s, _ := almanac.LookupSource("seed", soil, "soil")
		assert.Equal(t, seed, s, "LookupSource wrong for soil %d", soil)
	}
	for seed, location := range map[int]int{
		79: 82,
		14: 43,
		55: 86,
		13: 35,
		82: 46,
	} {
		assert.Equal(t, location, almanac.Lookup("seed", seed, "location"))
		l, _ := almanac.BoundedLookup("seed", seed, "location")
		assert.Equal(t, location, l, "BoundedLookup wrong for seed %d", seed)
		s, _ := almanac.LookupSource("seed", location, "location")
		assert.Equal(t, seed, s, "LookupSource wrong for location %d", location)
	}
}

func Test_Part1(t *testing.T) {
	almanac, err := NewAlmanac("input")
	require.NoError(t, err)
	log.Printf("Best Location is: %d", almanac.BestLocation())
}

func Test_Part2(t *testing.T) {
	almanac, err := NewAlmanac("input")
	require.NoError(t, err)

	best := almanac.BestLocation2()
	assert.Less(t, best, 53266420, "guess 1")
	log.Printf("Best Location for all seeds is: %d (WRONG)", best)

	best = almanac.BestLocation2b()
	assert.Less(t, best, 53266420, "guess 1")
	assert.Greater(t, best, 0, "guess 2")

	log.Printf("Best Location for all seeds is: %d", best)
}

func Test_Test(t *testing.T) {
	almanac, err := NewAlmanac("input")
	require.NoError(t, err)
	best := almanac.BestLocation2b()
	assert.Less(t, best, 53266420, "guess 1")
	assert.Greater(t, best, 0, "guess 2")
}

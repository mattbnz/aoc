package day14

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Sample(t *testing.T) {
	dish, err := NewDish("sample")
	require.NoError(t, err)
	expected, err := NewDish("sample-tilted")
	require.NoError(t, err)
	dish.Tilt(NORTH)
	assert.True(t, expected.Equal(dish.Grid))

	assert.Equal(t, 136, dish.Load())
}

func Test_Part1(t *testing.T) {
	dish, err := NewDish("input")
	require.NoError(t, err)
	dish.Tilt(NORTH)

	load := dish.Load()
	assert.Equal(t, 108840, load)
	t.Logf("Total Load: %d", load)
	t.Logf("Took %d checks % 12d/% 12d from cache (%.2f%%)", dish.checks, dish.hit, dish.miss, float64((dish.hit)/(dish.hit+dish.miss))*100.0)
}

func Test_Cycle1(t *testing.T) {
	dish, err := NewDish("sample")
	require.NoError(t, err)
	expected, err := NewDish("sample-cycle1")
	require.NoError(t, err)

	dish.Cycle()
	if !assert.True(t, expected.Equal(dish.Grid)) {
		t.Log("Got")
		dish.PrintNumbered()
		t.Log("Expected")
		expected.PrintNumbered()
	}
	t.Logf("Took %d checks % 12d/% 12d from cache (%.2f%%)", dish.checks, dish.hit, dish.miss, float64((dish.hit)/(dish.hit+dish.miss))*100.0)
}
func Test_Cycle2(t *testing.T) {
	dish, err := NewDish("sample")
	require.NoError(t, err)
	expected, err := NewDish("sample-cycle2")
	require.NoError(t, err)
	dish.Cycle()
	dish.Cycle()
	assert.True(t, expected.Equal(dish.Grid))
	t.Logf("Took %d checks % 12d/% 12d from cache (%.2f%%)", dish.checks, dish.hit, dish.miss, float64((dish.hit)/(dish.hit+dish.miss))*100.0)
}
func Test_Cycle3(t *testing.T) {
	dish, err := NewDish("sample")
	require.NoError(t, err)
	expected, err := NewDish("sample-cycle3")
	require.NoError(t, err)
	dish.Cycle()
	dish.Cycle()
	dish.Cycle()
	assert.True(t, expected.Equal(dish.Grid))
	t.Logf("Took %d checks % 12d/% 12d from cache (%.2f%%)", dish.checks, dish.hit, dish.miss, float64((dish.hit)/(dish.hit+dish.miss))*100.0)
}

func noTest_Part2(t *testing.T) {
	dish, err := NewDish("input")
	require.NoError(t, err)

	for i := 0; i < 1000000000; i++ {
		dish.Cycle()
		if i%1000 == 0 {
			t.Logf("%s: cycle=% 12d; % 12d/% 12d tilts (%.2f%%) from cache", time.Now(), i, dish.hit, dish.miss, float64((dish.hit)/(dish.hit+dish.miss))*100.0)
		}
	}

	load := dish.Load()
	t.Logf("Total Load: %d", load)
}

func noTest_Debug(t *testing.T) {
	dish, err := NewDish("sample")
	require.NoError(t, err)
	dish.PrintNumbered()
	dish.Tilt(SOUTH)
	dish.PrintNumbered()
}

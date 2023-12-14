package day14

import (
	"testing"

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
}
func Test_Cycle2(t *testing.T) {
	dish, err := NewDish("sample")
	require.NoError(t, err)
	expected, err := NewDish("sample-cycle2")
	require.NoError(t, err)
	dish.Cycle()
	dish.Cycle()
	assert.True(t, expected.Equal(dish.Grid))
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
}

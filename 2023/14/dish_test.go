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
	dish.Tilt()
	assert.True(t, expected.Equal(dish.Grid))

	assert.Equal(t, 136, dish.Load())
}

func Test_Part1(t *testing.T) {
	dish, err := NewDish("input")
	require.NoError(t, err)
	dish.Tilt()

	t.Logf("Total Load: %d", dish.Load())
}

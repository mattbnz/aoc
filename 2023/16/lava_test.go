package day16

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Sample(t *testing.T) {
	grid, err := NewMirrorGrid("sample")
	require.NoError(t, err)
	grid.PrintNumbered()
	energized := grid.Beam(1, 1, Pos{1, 1}, EAST)
	assert.Equal(t, 46, energized)
}

func Test_Part1(t *testing.T) {
	grid, err := NewMirrorGrid("input")
	require.NoError(t, err)
	energized := grid.Beam(1, 1, Pos{1, 1}, EAST)
	assert.Equal(t, 6921, energized)
	t.Logf("Enegerized Tiles: %d", energized)
}

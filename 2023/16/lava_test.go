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
	grid.StartBeam(Pos{1, 1}, EAST)
	grid.Wait()
	assert.Equal(t, 46, grid.Energized())
}

func Test_Part1(t *testing.T) {
	grid, err := NewMirrorGrid("input")
	require.NoError(t, err)
	grid.StartBeam(Pos{1, 1}, EAST)
	grid.Wait()
	t.Logf("Enegerized Tiles: %d", grid.Energized())
}

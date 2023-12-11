package day11

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Expansion(t *testing.T) {
	space, err := NewSpace("sample")
	require.NoError(t, err)

	r, err := os.OpenFile("expanded-sample", os.O_RDONLY, 0)
	require.NoError(t, err)
	defer r.Close()
	expected := NewGrid[BaseCell](r)

	space.Expand()
	if !assert.True(t, expected.Equal(space.g)) {
		t.Log("Got")
		space.g.PrintNumbered()
		t.Log("Expected")
		expected.PrintNumbered()
	}
}

func Test_Sample(t *testing.T) {
	space, err := NewSpace("sample")
	require.NoError(t, err)
	assert.Equal(t, 374, space.PathSum())

	assert.Equal(t, 9, len(space.Galaxies))
	space.g.PrintNumbered()

	assert.Equal(t, 9, space.PathLength(space.Galaxies[5], space.Galaxies[9]))
	assert.Equal(t, 15, space.PathLength(space.Galaxies[1], space.Galaxies[7]))
	assert.Equal(t, 17, space.PathLength(space.Galaxies[3], space.Galaxies[6]))
	assert.Equal(t, 5, space.PathLength(space.Galaxies[8], space.Galaxies[9]))

}

func Test_Part1(t *testing.T) {
	space, err := NewSpace("input")
	require.NoError(t, err)

	t.Logf("Path sum is %d", space.PathSum())
}

/*
func Test_Part2(t *testing.T) {
	maze, err := NewMaze("input")
	require.NoError(t, err)
	c, err := maze.CountEnclosed()
	require.NoError(t, err)
	t.Logf("Enclosed Tiles: %d", c)
}
*/

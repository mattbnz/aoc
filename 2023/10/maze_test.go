package day10

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Sample(t *testing.T) {
	maze, err := NewMaze("sample")
	require.NoError(t, err)
	maze.g.Print()
	assert.Equal(t, 4, maze.LongestPath())
	t.Log()

	maze, err = NewMaze("sample2")
	require.NoError(t, err)
	maze.g.Print()
	assert.Equal(t, 8, maze.LongestPath())
}

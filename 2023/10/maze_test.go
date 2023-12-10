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
	p, err := maze.LongestPath()
	require.NoError(t, err)
	assert.Equal(t, 4, p)
	t.Log()

	maze, err = NewMaze("sample2")
	require.NoError(t, err)
	maze.g.Print()
	p, err = maze.LongestPath()
	require.NoError(t, err)
	assert.Equal(t, 8, p)
}

func Test_Part1(t *testing.T) {
	maze, err := NewMaze("input")
	require.NoError(t, err)
	p, err := maze.LongestPath()
	require.NoError(t, err)
	t.Logf("Longest path is %d", p)
}

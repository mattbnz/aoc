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
}

func Test_Sample2(t *testing.T) {
	maze, err := NewMaze("sample2")
	require.NoError(t, err)
	maze.g.Print()
	p, err := maze.LongestPath()
	require.NoError(t, err)
	assert.Equal(t, 8, p)
}

func Test_Sample3(t *testing.T) {
	maze, err := NewMaze("sample3")
	maze.g.PrintNumbered()
	require.NoError(t, err)
	c, err := maze.CountEnclosed()
	require.NoError(t, err)
	assert.Equal(t, 4, c)
}

func Test_Sample4(t *testing.T) {
	maze, err := NewMaze("sample4")
	require.NoError(t, err)
	c, err := maze.CountEnclosed()
	require.NoError(t, err)
	assert.Equal(t, 4, c)
}

func Test_Sample5(t *testing.T) {
	maze, err := NewMaze("sample5")
	require.NoError(t, err)
	c, err := maze.CountEnclosed()
	require.NoError(t, err)
	assert.Equal(t, 8, c)
}

func Test_Sample6(t *testing.T) {
	maze, err := NewMaze("sample6")
	require.NoError(t, err)
	c, err := maze.CountEnclosed()
	require.NoError(t, err)
	assert.Equal(t, 10, c)
}

func Test_Part1(t *testing.T) {
	maze, err := NewMaze("input")
	require.NoError(t, err)
	p, err := maze.LongestPath()
	require.NoError(t, err)
	t.Logf("Longest path is %d", p)
}

func Test_Part2(t *testing.T) {
	maze, err := NewMaze("input")
	require.NoError(t, err)
	c, err := maze.CountEnclosed()
	require.NoError(t, err)
	t.Logf("Enclosed Tiles: %d", c)
}

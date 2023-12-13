package day13

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Sample(t *testing.T) {
	valley, err := NewMirrorValley("sample")
	require.NoError(t, err)
	assert.Len(t, valley.Mirrors, 2)

	p, _, ok := valley.ReflectsVertical(valley.Mirrors[0], -1)
	assert.True(t, ok)
	assert.Equal(t, 6, p)
	assert.Equal(t, 5, valley.Score(valley.Mirrors[0]))

	p, _, ok = valley.ReflectsHorizontal(valley.Mirrors[1], -1)
	assert.True(t, ok)
	assert.Equal(t, 5, p)
	assert.Equal(t, 400, valley.Score(valley.Mirrors[1]))

	assert.Equal(t, 405, valley.ScoreSum())
}

func Test_Part1(t *testing.T) {
	valley, err := NewMirrorValley("input")
	require.NoError(t, err)
	t.Logf("Sum is %d", valley.ScoreSum())
}

func Test_Sample_Part2(t *testing.T) {
	valley, err := NewMirrorValley("sample")
	require.NoError(t, err)
	assert.Len(t, valley.Mirrors, 2)

	assert.Equal(t, 300, valley.ScoreSmudged(valley.Mirrors[0]))
	assert.Equal(t, 100, valley.ScoreSmudged(valley.Mirrors[1]))
	assert.Equal(t, 400, valley.ScoreSmudgedSum())
}

func Test_Part2(t *testing.T) {
	valley, err := NewMirrorValley("input")
	require.NoError(t, err)
	score := valley.ScoreSmudgedSum()
	assert.Greater(t, 21356, score) // first guess.
	t.Logf("Sum is %d", score)
}

func Test_Debug(t *testing.T) {
	valley, err := NewMirrorValley("sample")
	require.NoError(t, err)
	assert.Len(t, valley.Mirrors, 2)

	assert.Equal(t, 100, valley.ScoreSmudged(valley.Mirrors[1]))

}

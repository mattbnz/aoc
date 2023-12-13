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

	p, ok := valley.ReflectsVertical(valley.Mirrors[0])
	assert.True(t, ok)
	assert.Equal(t, 6, p)
	assert.Equal(t, 5, valley.Score(valley.Mirrors[0]))

	p, ok = valley.ReflectsHorizontal(valley.Mirrors[1])
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

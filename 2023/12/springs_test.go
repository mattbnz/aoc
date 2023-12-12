package day12

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Sample(t *testing.T) {
	rows, err := NewSpringRows("sample")
	require.NoError(t, err)
	assert.Len(t, rows, 6)

	expect := []int{1, 4, 1, 1, 4, 10}
	for n, e := range expect {
		assert.Equal(t, e, rows.Arrangements(n))
	}
	assert.Equal(t, 21, rows.SumArrangements())
}

func Test_Part1(t *testing.T) {
	rows, err := NewSpringRows("input")
	require.NoError(t, err)

	t.Logf("Arrangement sum is %d", rows.SumArrangements())
}

func Test_Part2(t *testing.T) {
	rows, err := NewSpringRows("input")
	require.NoError(t, err)

	t.Logf("Arrangement sum is %d", rows.SumArrangements())
}

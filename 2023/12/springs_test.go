package day12

import (
	"testing"

	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewWith(t *testing.T) {
	sl := SpringList{S_OK}
	assert.Equal(t, SpringList{S_OK}, sl.NewWith(0, S_OK))
	assert.Equal(t, SpringList{S_BAD}, sl.NewWith(0, S_BAD))
	sl = SpringList{S_OK, S_BAD}
	assert.Equal(t, SpringList{S_BAD, S_BAD}, sl.NewWith(0, S_BAD))
	assert.Equal(t, SpringList{S_OK, S_OK}, sl.NewWith(1, S_OK))
	sl = SpringList{S_OK, S_BAD, S_OK}
	assert.Equal(t, SpringList{S_OK, S_OK, S_OK}, sl.NewWith(1, S_OK))
}

func Test_Sample(t *testing.T) {
	rows, err := NewSpringRows("sample")
	require.NoError(t, err)
	assert.Len(t, rows, 6)
	assert.Equal(t, 21, rows.SumArrangements())

	expect := []int{1, 4, 1, 1, 4, 10}
	for n, e := range expect {
		glog.V(1).Infof("Test Case %d", n)
		assert.Equal(t, e, rows.Arrangements(n), "sample %d", n)
	}
}

func Test_Part1(t *testing.T) {
	rows, err := NewSpringRows("input")
	require.NoError(t, err)
	a := rows.SumArrangements()
	assert.Equal(t, 7407, a)

	t.Logf("Arrangement sum is %d", a)
}

func Test_Part2(t *testing.T) {
	rows, err := NewSpringRows("input")
	require.NoError(t, err)

	t.Logf("Arrangement sum is %d", rows.SumArrangements())
}

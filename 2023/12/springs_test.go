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
	assert.Equal(t, 21, rows.SumArrangements(1))

	expect := []int{1, 4, 1, 1, 4, 10}
	for n, e := range expect {
		glog.V(1).Infof("Test Case %d", n)
		assert.Equal(t, e, rows.Arrangements(n, 1), "sample %d", n)
	}
}

func Test_Part1(t *testing.T) {
	rows, err := NewSpringRows("input")
	require.NoError(t, err)
	a := rows.SumArrangements(1)
	assert.Equal(t, 7407, a)

	t.Logf("Arrangement sum is %d", a)
}

func Test_CanMatch(t *testing.T) {
	sl := SpringList{S_OK, S_BAD}
	assert.True(t, sl.canMatch(SpringList{S_OK, S_BAD}, 0))
	assert.True(t, sl.canMatch(SpringList{S_OK}, 0))
	assert.True(t, sl.canMatch(SpringList{S_BAD}, 1))

	assert.False(t, sl.canMatch(SpringList{S_OK, S_BAD}, 1))
	assert.False(t, sl.canMatch(SpringList{S_OK}, 1))
	assert.False(t, sl.canMatch(SpringList{S_BAD}, 0))

	sl = SpringList{S_OK, S_UNKNOWN, S_BAD}
	assert.True(t, sl.canMatch(SpringList{S_OK, S_OK}, 0))
	assert.True(t, sl.canMatch(SpringList{S_BAD, S_BAD}, 1))
	assert.True(t, sl.canMatch(SpringList{S_OK, S_BAD}, 1))
	assert.True(t, sl.canMatch(SpringList{S_OK, S_BAD}, 0))

	assert.False(t, sl.canMatch(SpringList{S_BAD, S_OK}, 0))
}

func Test_LastBad(t *testing.T) {
	sl := SpringList{S_OK, S_BAD, S_BAD, S_OK, S_BAD, S_OK, S_OK, S_BAD}
	assert.Equal(t, 0, sl.lastBad(0))
	assert.Equal(t, 0, sl.lastBad(1))
	assert.Equal(t, 1, sl.lastBad(2))
	assert.Equal(t, 1, sl.lastBad(3))
	assert.Equal(t, 2, sl.lastBad(4))
	assert.Equal(t, 1, sl.lastBad(5))
	assert.Equal(t, 2, sl.lastBad(6))
	assert.Equal(t, 3, sl.lastBad(7))
}

func Test_Unfold(t *testing.T) {
	sr := SpringRow{
		Springs: SpringList{S_OK, S_BAD},
		BadRuns: Ints{1},
	}
	nr := sr.Unfold()
	assert.Equal(t, SpringList{S_OK, S_BAD, S_UNKNOWN, S_OK, S_BAD, S_UNKNOWN, S_OK, S_BAD, S_UNKNOWN, S_OK, S_BAD, S_UNKNOWN, S_OK, S_BAD}, nr.Springs)
	assert.Equal(t, Ints{1, 1, 1, 1, 1}, nr.BadRuns)
}

func Test_Sample_Part2(t *testing.T) {
	rows, err := NewSpringRows("sample")
	require.NoError(t, err)
	assert.Len(t, rows, 6)
	rows.Unfold()

	assert.Equal(t, 525152, rows.SumArrangements(5))

	expect := []int{1, 16384, 1, 16, 2500, 506250}
	for n, e := range expect {
		glog.V(1).Infof("Test Case %d", n)
		assert.Equal(t, e, rows.Arrangements(n, 5), "sample %d", n)
	}
}

func Test_Part2(t *testing.T) {
	rows, err := NewSpringRows("input")
	require.NoError(t, err)
	rows.Unfold()

	t.Logf("Arrangement sum is %d", rows.SumArrangements(5))
}

func Test_Debug(t *testing.T) {
	rows, err := NewSpringRows("sample")
	require.NoError(t, err)
	rows.Unfold()
	glog.Infof("%d Springs", len(rows[3].Springs))
	assert.Equal(t, 16, rows.Arrangements(3, 5))
}

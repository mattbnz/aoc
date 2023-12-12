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

func Test_Matches(t *testing.T) {
	for n, tCase := range []SpringRow{
		{SpringList{}, Ints{}},
		{SpringList{S_OK}, Ints{}},
		{SpringList{S_OK, S_OK}, Ints{}},
		{SpringList{S_BAD}, Ints{1}},
		{SpringList{S_BAD, S_BAD}, Ints{2}},
		{SpringList{S_BAD, S_OK, S_BAD, S_OK, S_BAD, S_BAD, S_BAD}, Ints{1, 1, 3}},
		{SpringList{S_OK, S_OK, S_BAD, S_OK, S_BAD, S_OK, S_BAD, S_BAD, S_BAD}, Ints{1, 1, 3}}, // OK prefix
		{SpringList{S_BAD, S_OK, S_BAD, S_OK, S_BAD, S_BAD, S_BAD, S_OK, S_OK}, Ints{1, 1, 3}}, // OK suffix
		{SpringList{S_BAD, S_OK, S_OK, S_BAD, S_OK, S_OK, S_BAD, S_BAD, S_BAD}, Ints{1, 1, 3}}, // OKs middle
	} {
		assert.True(t, tCase.Springs.Matches(tCase.BadRuns), "true case %d: %s", n, tCase)
	}
	for n, tCase := range []SpringRow{
		{SpringList{}, Ints{1}},
		{SpringList{S_OK}, Ints{1}},
		{SpringList{S_OK, S_OK}, Ints{1}},
		{SpringList{S_BAD}, Ints{2}},
		{SpringList{S_BAD}, Ints{}},
		{SpringList{}, Ints{1, 1, 3}},
		{SpringList{S_OK}, Ints{1, 1, 3}},
		{SpringList{S_BAD}, Ints{1, 1, 3}},
		{SpringList{S_OK, S_OK, S_BAD, S_OK, S_BAD, S_BAD, S_BAD}, Ints{1, 1, 3}},                                // matches 1,3
		{SpringList{S_OK, S_BAD, S_BAD, S_BAD, S_OK, S_OK, S_OK, S_OK, S_BAD, S_OK, S_OK, S_BAD}, Ints{3, 2, 1}}, // matches 1,3
	} {
		assert.False(t, tCase.Springs.Matches(tCase.BadRuns), "false case %d: %s", n, tCase)
	}
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

	t.Logf("Arrangement sum is %d", rows.SumArrangements())
}

func Test_Part2(t *testing.T) {
	rows, err := NewSpringRows("input")
	require.NoError(t, err)

	t.Logf("Arrangement sum is %d", rows.SumArrangements())
}

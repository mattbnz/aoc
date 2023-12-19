package day19

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testParts = []Part{
	{"x": 787, "m": 2655, "a": 1222, "s": 2876},
	{"x": 1679, "m": 44, "a": 2067, "s": 496},
	{"x": 2036, "m": 264, "a": 79, "s": 2244},
	{"x": 2461, "m": 1339, "a": 466, "s": 291},
	{"x": 2127, "m": 1623, "a": 2188, "s": 1013},
}

func Test_Sample(t *testing.T) {
	heap, err := NewHeap("sample")
	require.NoError(t, err)

	in := heap.Workflows["in"]
	for n, p := range testParts {
		assert.Equal(t, p, *in.Queue[n], "part %d does not match expected", n)
	}
	assert.Equal(t, &Rule{attr: "s", op: "<", val: 1351, dest: "px", RemainingSpace: map[string]int64{}}, in.Rules[0])
	assert.Equal(t, &Rule{dest: "qqz", RemainingSpace: map[string]int64{}}, in.Rules[1])
	assert.True(t, in.Rules[1].IsDefault())

	crn := heap.Workflows["crn"]
	assert.Equal(t, &Rule{attr: "x", op: ">", val: 2662, dest: "A", RemainingSpace: map[string]int64{}}, crn.Rules[0])
	assert.Equal(t, &Rule{dest: "R", RemainingSpace: map[string]int64{}}, crn.Rules[1])
	assert.True(t, crn.Rules[1].IsDefault())

	heap.SortParts()
	for n, w := range heap.Workflows {
		e, found := map[string]int{"A": 3, "R": 2}[n]
		if !found {
			e = 0
		}
		assert.Len(t, w.Queue, e, "Queue %s had unexpected length", n)
	}

	accepted := heap.Sum("A")
	assert.Equal(t, 19114, accepted)
}

func Test_Part1(t *testing.T) {
	heap, err := NewHeap("input")
	require.NoError(t, err)
	heap.SortParts()

	accepted := heap.Sum("A")
	t.Logf("Accepted Sum: %d", accepted)
}

func Test_Sample_Part2(t *testing.T) {
	heap, err := NewHeap("sample")
	require.NoError(t, err)

	heap.BuildGraph()
	in := heap.Workflows["in"]
	assert.EqualValues(t, 1350, in.Rules[0].Size)
	assert.EqualValues(t, 4000, in.Rules[1].RemainingSpace["x"])
	assert.EqualValues(t, 4000, in.Rules[1].RemainingSpace["m"])
	assert.EqualValues(t, 4000, in.Rules[1].RemainingSpace["a"])
	assert.EqualValues(t, 4000-1350, in.Rules[1].RemainingSpace["s"])
	px := heap.Workflows["px"]
	assert.EqualValues(t, 2005, px.Rules[0].Size)
	assert.EqualValues(t, 4000-2090, px.Rules[1].Size)
	assert.EqualValues(t, 4000, px.Rules[2].RemainingSpace["x"])
	assert.EqualValues(t, 2090, px.Rules[2].RemainingSpace["m"])
	assert.EqualValues(t, 4000-2005, px.Rules[2].RemainingSpace["a"])
	assert.EqualValues(t, 4000, px.Rules[2].RemainingSpace["s"])
	assert.Equal(t, big.NewInt(167409079868000), heap.Combinations())
}

func Test_CalcAttrSize(t *testing.T) {
	rSizes, aSizes := calcAttrSize([]Rule{
		{attr: "x", op: ">", val: 3000},
		{attr: "m", op: "<", val: 2},
		{attr: "a", op: ">", val: 3990},
		{dest: "default"},
	}, false)
	assert.Equal(t, []int64{1000, 1, 10, 0}, rSizes)
	assert.Equal(t, map[string]int64{"x": 1000, "m": 1, "a": 10, "s": 0}, aSizes)

	rSizes, aSizes = calcAttrSize([]Rule{
		{attr: "x", op: "<", val: 10},
		{attr: "x", op: ">", val: 3990},
	}, false)
	assert.Equal(t, []int64{9, 10}, rSizes)
	assert.Equal(t, map[string]int64{"x": 19, "m": 0, "a": 0, "s": 0}, aSizes)

	rSizes, aSizes = calcAttrSize([]Rule{
		{attr: "x", op: "<", val: 10},
		{attr: "x", op: ">", val: 3990},
		{dest: "default", RemainingSpace: map[string]int64{"a": 1000}}, // ignored
	}, false)
	assert.Equal(t, []int64{9, 10, 0}, rSizes)
	assert.Equal(t, map[string]int64{"x": 19, "m": 0, "a": 0, "s": 0}, aSizes)

	rSizes, aSizes = calcAttrSize([]Rule{
		{attr: "x", op: "<", val: 10},
		{attr: "x", op: ">", val: 3990},
		{dest: "default", RemainingSpace: map[string]int64{"a": 1000}}, // used
	}, true)
	assert.Equal(t, []int64{9, 10, 0}, rSizes)
	assert.Equal(t, map[string]int64{"x": 19, "m": 0, "a": 1000, "s": 0}, aSizes)

	rSizes, aSizes = calcAttrSize([]Rule{
		{attr: "a", op: ">", val: 10},
		{dest: "default", RemainingSpace: map[string]int64{"a": 1000}}, // used
	}, true)
	assert.Equal(t, []int64{3990, 0}, rSizes)
	assert.Equal(t, map[string]int64{"x": 0, "m": 0, "a": 1000, "s": 0}, aSizes)

	rSizes, aSizes = calcAttrSize([]Rule{
		{attr: "a", op: ">", val: 3990},
		{dest: "default", RemainingSpace: map[string]int64{"a": 1000}}, // used
	}, true)
	assert.Equal(t, []int64{10, 0}, rSizes)
	assert.Equal(t, map[string]int64{"x": 0, "m": 0, "a": 10, "s": 0}, aSizes)
}

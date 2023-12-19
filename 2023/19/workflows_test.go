package day19

import (
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
	assert.Equal(t, Rule{attr: "s", op: "<", val: 1351, dest: "px"}, in.Rules[0])
	assert.Equal(t, Rule{dest: "qqz"}, in.Rules[1])
	assert.True(t, in.Rules[1].IsDefault())

	crn := heap.Workflows["crn"]
	assert.Equal(t, Rule{attr: "x", op: ">", val: 2662, dest: "A"}, crn.Rules[0])
	assert.Equal(t, Rule{dest: "R"}, crn.Rules[1])
	assert.True(t, crn.Rules[1].IsDefault())

	heap.SortParts()

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

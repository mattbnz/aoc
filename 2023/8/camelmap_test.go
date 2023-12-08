package day8

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Parsing(t *testing.T) {
	m, err := NewMap("sample")
	require.NoError(t, err)

	assert.Equal(t, "RL", m.Instructions)
	assert.Equal(t, 7, len(m.Nodes))
	assert.Equal(t, "BBB", m.Nodes["AAA"].Elements['L'])
	assert.Equal(t, "CCC", m.Nodes["AAA"].Elements['R'])
}

func Test_Sample(t *testing.T) {
	m, err := NewMap("sample")
	require.NoError(t, err)
	assert.Equal(t, 2, m.StepsFrom("AAA", "ZZZ"))

	m, err = NewMap("sample2")
	require.NoError(t, err)
	assert.Equal(t, 6, m.StepsFrom("AAA", "ZZZ"))
}

func Test_Part1(t *testing.T) {
	m, err := NewMap("input")
	require.NoError(t, err)
	log.Printf("Steps: %d", m.StepsFrom("AAA", "ZZZ"))
}

func Test_Part2_Sample(t *testing.T) {
	m, err := NewMap("sample3")
	require.NoError(t, err)
	assert.Equal(t, 6, m.SimultaneousSteps())
}

func Test_Part2(t *testing.T) {
	m, err := NewMap("input")
	require.NoError(t, err)
	log.Printf("Simultaneous Steps: %d", m.SimultaneousSteps())
}

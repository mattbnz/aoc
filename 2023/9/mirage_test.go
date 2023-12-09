package day9

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Extrapolate(t *testing.T) {
	s := Scan{}
	assert.Equal(t, 0, s.Extrapolate(Reading{0}))
	assert.Equal(t, 0, s.Extrapolate(Reading{0, 0, 0}))
	assert.Equal(t, 1, s.Extrapolate(Reading{1, 1, 1}))
}

func Test_Sample(t *testing.T) {
	scan, err := NewScan("sample")
	require.NoError(t, err)
	assert.Equal(t, 3, len(scan.Readings))
	assert.Equal(t, 114, scan.ExtrapolateAndSum())
}

func Test_Part1(t *testing.T) {
	scan, err := NewScan("input")
	require.NoError(t, err)
	log.Printf("Sum is: %d", scan.ExtrapolateAndSum())
}

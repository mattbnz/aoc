package day18

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Sample(t *testing.T) {
	lagoon, err := NewLagoon("sample")
	require.NoError(t, err)
	lagoon.VisitAll()
	lagoon.PrintNumbered()

	length := lagoon.TrenchLength()
	assert.Equal(t, 38, length)
	volume := lagoon.Volume()
	assert.Equal(t, 62, volume)
}

func testLagoon(s string) (l Lagoon) {
	l.FlexGrid = NewGridFromScanner[BaseCell](bufio.NewScanner(bytes.NewBufferString(s)))
	return
}

func Test_Part1(t *testing.T) {
	lagoon, err := NewLagoon("input")
	require.NoError(t, err)

	volume := lagoon.Volume()
	assert.Less(t, volume, 100724)   // First Guess
	assert.Less(t, volume, 93240)    // Second Guess
	assert.Greater(t, volume, 90818) // Third Guess
	t.Logf("Lagon Volume: %d", volume)
}

func Test_Debug(t *testing.T) {
	lagoon, err := NewLagoon("input")
	require.NoError(t, err)
	v := 0
	for row := lagoon.minrow; row <= lagoon.maxrow; row++ {
		v += lagoon.RowVolume(row)
	}
	volume := lagoon.Volume()
	assert.Equal(t, v, volume)
	volume = v
	assert.Less(t, volume, 100724)   // First Guess
	assert.Less(t, volume, 93240)    // Second Guess
	assert.Greater(t, volume, 90818) // Third Guess
	t.Logf("Lagon Volume: %d", volume)
}

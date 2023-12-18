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

func assertRow(t *testing.T, s string, n int) {
	tl := testLagoon(s)
	tl.PrintNumbered()
	assert.Equal(t, n, tl.RowVolume(1))

}

func Test_InputDebug(t *testing.T) {
	assertRow(t, ".####.####.", 9)
	assertRow(t, ".###########.", 11)

	lagoon, err := NewLagoon("input")
	require.NoError(t, err)

	assert.Equal(t, 38, lagoon.RowVolume(lagoon.minrow))

}

func Test_Part1(t *testing.T) {
	lagoon, err := NewLagoon("input")
	require.NoError(t, err)

	volume := lagoon.Volume()
	assert.Less(t, volume, 100724) // First Guess
	assert.Less(t, volume, 93240)  // First Guess
	t.Logf("Lagon Volume: %d", volume)
}

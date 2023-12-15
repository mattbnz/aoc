package day15

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Sample(t *testing.T) {
	manual := Manual{}
	assert.Equal(t, 52, manual.HashInstruction("HASH"))
	assert.Equal(t, 30, manual.HashInstruction("rn=1"))
	assert.Equal(t, 253, manual.HashInstruction("cm-"))
	assert.Equal(t, 97, manual.HashInstruction("qp=3"))
	assert.Equal(t, 47, manual.HashInstruction("cm=2"))
	assert.Equal(t, 14, manual.HashInstruction("qp- "))
	assert.Equal(t, 180, manual.HashInstruction("pc=4"))
	assert.Equal(t, 9, manual.HashInstruction("ot=9"))
	assert.Equal(t, 197, manual.HashInstruction("ab=5"))
	assert.Equal(t, 48, manual.HashInstruction("pc-"))
	assert.Equal(t, 214, manual.HashInstruction("pc=6"))
	assert.Equal(t, 231, manual.HashInstruction("ot=7"))
	assert.Equal(t, 19, manual.HashInstruction("tjcr"))

	manual, err := NewManual("sample")
	require.NoError(t, err)
	assert.Equal(t, 1320, manual.Hash())
}

func Test_Part1(t *testing.T) {
	manual, err := NewManual("input")
	require.NoError(t, err)

	hash := manual.Hash()
	assert.Greater(t, hash, 93549)
	t.Logf("Instruction Hash: %d", hash)
}

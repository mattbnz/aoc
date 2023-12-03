package day3

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseRow(t *testing.T) {
	grid := Grid{}
	grid.ParseRow("467..114..")
	assert.Len(t, grid.cells, 1)
}

func run(t *testing.T, filename string, skipLines int, limitLines int) int {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	grid := Grid{}
	n := 0
	for s.Scan() {
		if skipLines == -1 || (n > skipLines) {
			grid.ParseRow(s.Text())
		}
		n++
		if limitLines != -1 && n >= limitLines {
			break
		}
	}
	grid.Print()
	fmt.Println()
	sum := 0
	for _, n := range grid.FindPartNumbers() {
		sum += n
	}
	return sum
}

func Test_Sample(t *testing.T) {
	assert.Equal(t, 4361, run(t, "sample", -1, -1))
}

func Test_Input(t *testing.T) {
	log.Printf("Sum is: %d", run(t, "input", -1, -1))
}

func Test__Head(t *testing.T) {
	log.Printf("Sum is: %d", run(t, "input", -1, 4))
}

func Test_Middle(t *testing.T) {
	log.Printf("Sum is: %d", run(t, "input", 136, 140))
}

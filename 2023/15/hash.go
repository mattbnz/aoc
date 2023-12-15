package day15

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

type Lens struct {
	Label  string
	Length int
}

type Box struct {
	Lenses []Lens
}

func (b Box) Index(label string) int {
	for n, l := range b.Lenses {
		if l.Label == label {
			return n
		}
	}
	return -1
}

type Manual struct {
	Steps []string

	Boxes []*Box
}

func NewManual(filename string) (m Manual, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	m.Boxes = make([]*Box, 256)
	for b := range m.Boxes {
		m.Boxes[b] = &Box{}
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() == "" {
			return
		}
		ins := strings.Split(s.Text(), ",")
		m.Steps = append(m.Steps, ins...)
	}
	return
}

func (m Manual) HashInstruction(s string) (sum int) {
	for _, r := range s {
		sum += int(r)
		sum *= 17
		sum %= 256
	}
	return
}

func (m Manual) Hash() (sum int) {
	for _, s := range m.Steps {
		h := m.HashInstruction(s)
		glog.V(1).Infof("% 10s => %d", s, h)
		sum += h
	}
	return
}

func (m Manual) Focus() (sum int) {
	for _, s := range m.Steps {
		if strings.HasSuffix(s, "-") {
			// Replace
			label := strings.TrimSuffix(s, "-")
			boxIdx := m.HashInstruction(label)
			box := m.Boxes[boxIdx]
			lensIdx := box.Index(label)
			if lensIdx != -1 {
				box.Lenses = append(box.Lenses[:lensIdx], box.Lenses[lensIdx+1:]...)
			}
		} else {
			label, lengthS, ok := strings.Cut(s, "=")
			if !ok {
				glog.Fatalf("Step %s: invalid format", s)
			}
			length, err := strconv.Atoi(lengthS)
			if err != nil {
				glog.Fatalf("Step %s: invalid format: %v", s, err)
			}
			boxIdx := m.HashInstruction(label)
			box := m.Boxes[boxIdx]
			lensIdx := box.Index(label)
			if lensIdx == -1 {
				box.Lenses = append(box.Lenses, Lens{Label: label, Length: length})
			} else {
				box.Lenses[lensIdx].Length = length
			}
		}
	}

	for boxNum, box := range m.Boxes {
		for lensNum, l := range box.Lenses {
			power := (boxNum + 1) * (lensNum + 1) * l.Length
			glog.V(1).Infof("Lens %s at %d in box %d with length %d has power %d", l.Label, lensNum+1, boxNum+1, l.Length, power)
			sum += power
		}
	}
	return
}

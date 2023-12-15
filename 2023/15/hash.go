package day15

import (
	"bufio"
	"os"
	"strings"

	"github.com/golang/glog"
)

type Manual struct {
	Steps []string
}

func NewManual(filename string) (dish Manual, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() == "" {
			return
		}
		ins := strings.Split(s.Text(), ",")
		for _, i := range ins {
			dish.Steps = append(dish.Steps, i)
		}
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

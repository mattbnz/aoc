// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 8.
// Haunted Wasteland.

package day8

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/golang/glog"
)

type Node struct {
	Name string
	m    *Map

	Elements map[rune]string
}

type Map struct {
	Instructions string

	Nodes map[string]*Node
}

var reInstruct = regexp.MustCompile(`^[RL]+$`)
var mapLine = regexp.MustCompile(`^([A-Z0-9]{3}) = \(([A-Z0-9]{3}), ([A-Z0-9]{3})\)$`)

func NewMap(filename string) (m Map, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	if !s.Scan() {
		return m, fmt.Errorf("missing instructions line")
	}
	if !reInstruct.MatchString(s.Text()) {
		return m, fmt.Errorf("invalid instruction line")
	}
	m.Instructions = s.Text()

	if !s.Scan() {
		return m, fmt.Errorf("missing blank line following instructions")
	}

	m.Nodes = make(map[string]*Node)
	for s.Scan() {
		l := s.Text()
		p := mapLine.FindStringSubmatch(l)
		if len(p) != 4 {
			return m, fmt.Errorf("bad map line: %s", l)
		}
		m.Nodes[p[1]] = &Node{
			Name: p[1],
			m:    &m,
			Elements: map[rune]string{
				'L': p[2],
				'R': p[3],
			},
		}
	}
	return
}

func (m Map) StepsFrom(start, end string) int {
	ip := 0
	node := m.Nodes["AAA"]
	n := 0
	for node.Name != end {
		glog.V(1).Infof("Step % 3d: %c from %s", n, m.Instructions[ip], node.Name)
		node = m.Nodes[node.Elements[rune(m.Instructions[ip])]]
		ip = (ip + 1) % len(m.Instructions)
		n++
	}
	return n
}

func print(l []*Node) string {
	nl := []string{}
	for _, n := range l {
		nl = append(nl, n.Name)
	}
	return strings.Join(nl, ", ")
}

func (m Map) SimultaneousSteps() int {
	nodes := []*Node{}
	for name, node := range m.Nodes {
		if strings.HasSuffix(name, "A") {
			nodes = append(nodes, node)
		}
	}
	glog.V(1).Infof("Found %d starting nodes", len(nodes))
	glog.V(1).Infof("Found %d instructions: %s", len(m.Instructions), m.Instructions)
	ip := 0
	step := 0
	for {
		glog.V(1).Infof("Step % 3d: ip=%d going %c from %s", step, ip, m.Instructions[ip], print(nodes))
		zEnds := 0
		for n := 0; n < len(nodes); n++ {
			nodes[n] = m.Nodes[nodes[n].Elements[rune(m.Instructions[ip])]]
			if strings.HasSuffix(nodes[n].Name, "Z") {
				zEnds++
			}
		}
		if zEnds > 0 {
			glog.Infof("Step % 3d: went %c arrived at %d Zs", step, m.Instructions[ip], zEnds)
		}
		ip = (ip + 1) % len(m.Instructions)
		step++
		if zEnds == len(nodes) {
			break
		}
	}
	return step
}

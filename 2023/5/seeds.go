// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 5.
// If You Give A Seed A Fertilizer

package day5

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Mapping struct {
	Source string
	Dest   string

	Overrides []*Override
}

func (m Mapping) String() string {
	return fmt.Sprintf("%s-to-%s map: %d overrides", m.Source, m.Dest, len(m.Overrides))
}

func (m Mapping) Map(id int) int {
	for _, o := range m.Overrides {
		if n := o.Contains(id); n != -1 {
			return o.DestBase + n
		}
	}
	return id
}

type Override struct {
	SourceBase int
	DestBase   int
	Count      int
}

func (o Override) Contains(id int) int {
	if id >= o.SourceBase && id < o.SourceBase+o.Count {
		return id - o.SourceBase
	}
	return -1
}

type Almanac struct {
	Seeds []int

	Maps []*Mapping
}

func NewAlmanac(filename string) (a Almanac, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	if !s.Scan() {
		return a, fmt.Errorf("missing seeds line")
	}
	_, seeds, ok := strings.Cut(s.Text(), ": ")
	a.Seeds = numberList(seeds)
	if !ok || len(a.Seeds) < 1 {
		return a, fmt.Errorf("implausible seed count: %#v", a.Seeds)
	}
	if !s.Scan() {
		return a, fmt.Errorf("missing blank line following seeds")
	}

	var currentMap *Mapping
	for s.Scan() {
		l := s.Text()
		if l == "" {
			a.Maps = append(a.Maps, currentMap)
			currentMap = nil
			continue
		}

		if strings.HasSuffix(l, ":") {
			if currentMap != nil {
				return a, fmt.Errorf("found map header (%s) but already have map (%s)", l, currentMap)
			}
			map_name, _, found := strings.Cut(l, " ")
			if !found {
				return a, fmt.Errorf("bad map header on line: %s", l)
			}
			parts := strings.Split(map_name, "-")
			if len(parts) != 3 {
				return a, fmt.Errorf("bad map name on line; %s", l)
			}
			currentMap = &Mapping{Source: parts[0], Dest: parts[2]}
			continue
		}

		if currentMap == nil {
			return a, fmt.Errorf("found map line (%s) but don't have a map source", l)
		}
		nums := numberList(l)
		if len(nums) != 3 {
			return a, fmt.Errorf("bad map line: %s", l)
		}
		currentMap.Overrides = append(currentMap.Overrides, &Override{DestBase: nums[0], SourceBase: nums[1], Count: nums[2]})
	}
	if currentMap != nil {
		a.Maps = append(a.Maps, currentMap)
	}
	return
}

func (a *Almanac) getMap(source string) *Mapping {
	for _, m := range a.Maps {
		if m.Source == source {
			return m
		}
	}
	return nil
}

func (a *Almanac) Lookup(source string, id int, dest string) int {
	m := a.getMap(source)
	if m == nil {
		log.Fatalf("unknown category: %s", source)
	}
	d := m.Map(id)
	if m.Dest == dest {
		return d
	}
	return a.Lookup(m.Dest, d, dest)
}

func (a *Almanac) BestLocation() int {
	locs := []int{}
	for _, s := range a.Seeds {
		locs = append(locs, a.Lookup("seed", s, "location"))
	}
	sort.Ints(locs)
	return locs[0]
}

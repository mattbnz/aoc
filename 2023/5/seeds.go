// Copyright (C) 2023 Matt Brown

// Advent of Code 2023 - Day 5.
// If You Give A Seed A Fertilizer

package day5

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/golang/glog"
)

type cacheKey struct {
	source     string
	dest       string
	sourceBase int
	count      int
}

type Mapping struct {
	Source string
	Dest   string

	Overrides []*Override

	Cache     map[cacheKey]Override
	hit, miss int

	sorted bool
}

func (m Mapping) String() string {
	return fmt.Sprintf("%s-to-%s map: %d overrides", m.Source, m.Dest, len(m.Overrides))
}

func (m Mapping) Map(id int, max int) (int, Override) {
	if !m.sorted {
		sort.Slice(m.Overrides, func(i, j int) bool {
			return m.Overrides[i].SourceBase < m.Overrides[j].SourceBase
		})
	}
	lastEnd := 0
	for _, o := range m.Overrides {
		if o.SourceBase > id {
			// Wasn't in previous override, not in this one, fake up a gap override
			return id, Override{SourceBase: lastEnd, DestBase: lastEnd, Count: o.SourceBase - lastEnd}
		}
		if d := o.ToDest(id); d != -1 {
			return d, *o
		}
		lastEnd = o.SourceBase + o.Count
	}
	// Bigger than any override
	return id, Override{SourceBase: lastEnd, DestBase: lastEnd, Count: max - lastEnd}
}

func (m *Mapping) FromCache(dest string, id int) (rv int, bounds Override, found bool) {
	for k, o := range m.Cache {
		if k.dest != dest {
			continue
		}
		if id < k.sourceBase || id >= k.sourceBase+k.count {
			continue
		}
		m.hit++
		return o.DestBase + (id - k.sourceBase), o, true
	}
	m.miss++
	return -1, Override{}, false
}

func (m *Mapping) PutCache(dest string, o Override, sourceID, destID int) {
	if m.Cache == nil {
		m.Cache = map[cacheKey]Override{}
	}
	key := cacheKey{source: m.Source, dest: dest, sourceBase: o.SourceBase, count: o.Count}
	m.Cache[key] = o
	glog.V(1).Infof("Cache Add: %#v => %#v", key, o)
}

type Override struct {
	SourceBase int
	DestBase   int
	Count      int
}

func (o Override) ToSource(dest int) int {
	if dest < o.DestBase || dest >= o.DestBase+o.Count {
		return -1
	}
	return o.SourceBase + (dest - o.DestBase)
}

func (o Override) ToDest(source int) int {
	if source < o.SourceBase || source >= o.SourceBase+o.Count {
		return -1
	}
	return o.DestBase + (source - o.SourceBase)
}

type Almanac struct {
	Seeds []int

	Maps []*Mapping

	Max int // biggest thing we see during parsing.
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
	for n := 0; n < len(a.Seeds); n += 2 {
		a.Max = Max(a.Max, a.Seeds[n]+a.Seeds[n+1])
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
		a.Max = Max(a.Max, nums[0]+nums[2], nums[1]+nums[2])
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
		glog.Fatalf("unknown category: %s", source)
	}
	d, _ := m.Map(id, a.Max)
	if m.Dest == dest {
		return d
	}
	return a.Lookup(m.Dest, d, dest)
}

func (a *Almanac) BoundedLookup(source string, id int, dest string) (int, Override) {
	callPrefix := fmt.Sprintf("BoundedLookup (%s,%d,%s) =>", source, id, dest)
	m := a.getMap(source)
	if m == nil {
		glog.Fatalf("unknown category: %s", source)
	}
	if source == "seed" && (m.hit+m.miss) > 0 && (m.hit+m.miss)%100000 == 0 {
		glog.V(2).Infof("cache has %d/%d hits %.2f%%", m.hit, m.hit+m.miss, float64(m.hit)/float64(m.hit+m.miss)*100.0)
	}
	if d, b, found := m.FromCache(dest, id); found {
		glog.V(1).Infof("%s (%s,%d) (%#v) (from cache)", callPrefix, dest, d, b)
		return d, b
	}
	d, b := m.Map(id, a.Max)
	if m.Dest == dest {
		glog.V(1).Infof("%s (%s,%d) (%#v)", callPrefix, dest, d, b)
		return d, b
	}
	d2, b2 := a.BoundedLookup(m.Dest, d, dest)
	final := b
	final.SourceBase = Max(b.SourceBase, b.ToSource(b2.SourceBase))
	trimmedBase := final.SourceBase - b.SourceBase
	final.DestBase += trimmedBase
	final.Count = Min(b.Count-trimmedBase, (b2.SourceBase+b2.Count)-final.DestBase)
	m.PutCache(dest, final, id, d2)
	glog.V(1).Infof("%s (%s,%d) (%#v); (%#v) => (%s,%d) (%#v)",
		callPrefix, dest, d2, final, b, m.Dest, d, b2)
	return d2, final
}

func (a *Almanac) BestLocation() int {
	locs := []int{}
	for _, s := range a.Seeds {
		locs = append(locs, a.Lookup("seed", s, "location"))
	}
	sort.Ints(locs)
	return locs[0]
}

func (a *Almanac) BestLocation2() int {
	locs := []int{}
	for n := 0; n < len(a.Seeds); n += 2 {
		seed, count := a.Seeds[n], a.Seeds[n+1]
		for c := 0; c < count; c++ {
			l, b := a.BoundedLookup("seed", seed+c, "location")
			glog.Infof("%d: Seed %d (%d+%d) goes to %d (%#v)", n, seed+c, seed, c, l, b)
			locs = append(locs, l)
			c += b.Count
		}
	}
	sort.Ints(locs)
	return locs[0]
}

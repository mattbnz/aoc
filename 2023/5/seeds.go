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
type cacheResult struct {
	override Override
	diff     int
}

type Mapping struct {
	Source string
	Dest   string

	Overrides []*Override

	Cache     map[cacheKey]cacheResult
	hit, miss int

	sorted     bool
	sortedDest bool
	oLByDest   []*Override
}

func (m Mapping) String() string {
	return fmt.Sprintf("%s-to-%s map: %d overrides", m.Source, m.Dest, len(m.Overrides))
}

// Returns the destination id for the given source id
func (m *Mapping) ToDest(id int, max int) (int, Override) {
	if !m.sorted {
		sort.Slice(m.Overrides, func(i, j int) bool {
			return m.Overrides[i].SourceBase < m.Overrides[j].SourceBase
		})
		m.sorted = true
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

// Returns the source id for the given destination id
func (m *Mapping) ToSource(id int, max int) (int, Override) {
	if !m.sortedDest {
		m.oLByDest = append(m.oLByDest, m.Overrides...)
		sort.Slice(m.oLByDest, func(i, j int) bool {
			return m.oLByDest[i].DestBase < m.oLByDest[j].DestBase
		})
		m.sortedDest = true
	}
	lastEnd := 0
	for _, o := range m.oLByDest {
		if o.DestBase > id {
			// Wasn't in previous override, not in this one, fake up a gap override
			return id, Override{SourceBase: lastEnd, DestBase: lastEnd, Count: o.DestBase - lastEnd}
		}
		if d := o.ToSource(id); d != -1 {
			return d, *o
		}
		lastEnd = o.DestBase + o.Count
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
		return id + o.diff, o.override, true
	}
	m.miss++
	return -1, Override{}, false
}

func (m *Mapping) PutCache(dest string, o Override, sourceID, destID int) {
	if m.Cache == nil {
		m.Cache = map[cacheKey]cacheResult{}
	}
	key := cacheKey{source: m.Source, dest: dest, sourceBase: o.SourceBase, count: o.Count}
	m.Cache[key] = cacheResult{override: o, diff: destID - sourceID}
	glog.V(1).Infof("Cache Add: %#v => %#v", key, m.Cache[key])
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

func (a *Almanac) getDestMap(dest string) *Mapping {
	for _, m := range a.Maps {
		if m.Dest == dest {
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
	d, _ := m.ToDest(id, a.Max)
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
	d, b := m.ToDest(id, a.Max)
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
			glog.Infof("% 3d: Seed % 12d (% 12d +% 12d)\tgoes to % 12d", n, seed+c, seed, c, l)
			locs = append(locs, l)
			c += (b.SourceBase + b.Count)
			glog.V(2).Infof("      - jumping search to seed %d from %#v", c, b)
			c -= seed + 1 // for base seed and loop increment itself.
		}
	}
	sort.Ints(locs)
	return locs[0]
}

// Returns the first (lowest) seed within bounds that exists in the starting seed list, along with
// it's offset from the base of bounds.
func (a *Almanac) LowestStart(bounds Override) (startSeed int, startOffset int) {
	return -1, -1 // TODO
}

func (a *Almanac) LookupSource(source string, destID int, dest string) (int, Override) {
	callPrefix := fmt.Sprintf("ReverseLookup (%s,%d,%s) =>", source, destID, dest)
	m := a.getDestMap(dest)
	if m == nil {
		glog.Fatalf("unknown category: %s", dest)
	}
	s, b := m.ToSource(destID, a.Max)
	if m.Source == source {
		glog.V(1).Infof("%s (%s,%d) (%#v)", callPrefix, source, s, b)
		return s, b
	}
	s2, b2 := a.LookupSource(source, s, m.Source)
	final := b
	final.DestBase = Max(b.DestBase, b.ToDest(b2.DestBase))
	trimmedBase := final.DestBase - b.DestBase
	final.SourceBase += trimmedBase
	final.Count = Min(b.Count-trimmedBase, (b2.SourceBase+b2.Count)-final.SourceBase)

	/*
		final.SourceBase = Max(b.SourceBase, b.ToSource(b2.SourceBase))
		trimmedBase := final.SourceBase - b.SourceBase
		final.DestBase += trimmedBase
		final.Count = Min(b.Count-trimmedBase, (b2.SourceBase+b2.Count)-final.DestBase)
	*/
	//m.PutCache(dest, final, id, s2)
	glog.V(1).Infof("%s (%s,%d) (%#v); (%#v) => (%s,%d) (%#v)",
		callPrefix, source, s2, final, b, m.Source, s, b2)
	return s2, final
}

func (a *Almanac) BestLocation2b() int {
	m := a.getMap("humidity")
	if m == nil {
		glog.Fatalf("can't find humidity map")
	}

	for loc := 0; loc < a.Max; loc++ {
		seed, b := a.LookupSource("seed", loc, "location")
		glog.Infof("Location % 12d maps to seed % 12d with bounds %#v", loc, seed, b)
		if startSeed, startOffset := a.LowestStart(b); startSeed != -1 {
			glog.Infof(" + Bounds contain starting seed % 12d at offset % 12d!", startSeed, startOffset)
			return loc + startOffset
		}

	}
	glog.Fatalf("Did not find any location that mapped to a starting seed!")
	return -1
}

// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 19, Puzzle 1.
// Not Enough Minerals. Robots to the rescue.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Int(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("%s is not an int: %v", s, err)
	}
	return v
}

type Resource int

const (
	ORE Resource = iota
	CLAY
	OBSIDIAN
	GEODE
)

func (r Resource) String() string {
	return []string{"ore", "clay", "obsidian", "geode"}[r]
}

func (r Resource) Action() string {
	if r == GEODE {
		return "crack"
	}
	return "collect"
}

type ResourceData map[Resource]int
type Cost ResourceData

func (c Cost) String() string {
	cost := []string{}
	for nr, nq := range c {
		cost = append(cost, fmt.Sprintf("%d %s", nq, nr))
	}
	return strings.Join(cost, " and ")
}

type Blueprint struct {
	number int

	robotCost map[Resource]Cost

	inventory ResourceData

	robots ResourceData
}

func (bp *Blueprint) PrettyPrint() {
	fmt.Printf("Blueprint %d:\n", bp.number)
	for r, cost := range bp.robotCost {
		fmt.Printf("  Each %s robot costs %s.\n", r, cost)
	}
	fmt.Println()
}

func (bp *Blueprint) Collect() {
	for r, count := range bp.robots {
		bp.inventory[r] += count
		fmt.Printf("%d %s-%sing robot %ss %d %s; you now have %d %s.\n", count, r, r.Action(), r.Action(), count, r, bp.inventory[r], r)
	}
}

func (bp *Blueprint) canBuild(r Resource) bool {
	for nr, nq := range bp.robotCost[r] {
		if bp.inventory[nr] < nq {
			return false
		}
	}
	return true
}

func (bp *Blueprint) consume(r Resource) {
	for nr, nq := range bp.robotCost[r] {
		bp.inventory[nr] -= nq
	}
	fmt.Printf("Spend %s to start building an %s-%sing robot.\n", bp.robotCost[r], r, r.Action())
}

func (bp *Blueprint) Build() ResourceData {
	queue := ResourceData{}
	for _, r := range []Resource{GEODE, OBSIDIAN, CLAY, ORE} {
		for bp.canBuild(r) {
			queue[r] = 1
			bp.consume(r)
		}
	}
	return queue
}

func (bp *Blueprint) CollectRobots(queue ResourceData) {
	for r, count := range queue {
		bp.robots[r] += count
		fmt.Printf("The new %s-%sing robot is ready; you now have %d of them.\n", r, r.Action(), bp.robots[r])
	}
}

var BP_RE = regexp.MustCompile(`Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. ` +
	`Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.`)

func NewBlueprint(l string) Blueprint {
	m := BP_RE.FindStringSubmatch(l)
	if len(m) != 8 {
		log.Fatal("Not a blueprint: ", l)
	}
	costs := map[Resource]Cost{
		ORE:      {ORE: Int(m[2])},
		CLAY:     {ORE: Int(m[3])},
		OBSIDIAN: {ORE: Int(m[4]), CLAY: Int(m[5])},
		GEODE:    {ORE: Int(m[6]), OBSIDIAN: Int(m[7])},
	}
	return Blueprint{
		number:    Int(m[1]),
		robotCost: costs,
		inventory: ResourceData{},
		robots:    ResourceData{ORE: 1},
	}
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	bps := []Blueprint{}
	n := 0
	for s.Scan() {
		b := NewBlueprint(s.Text())
		bps = append(bps, b)
		b.PrettyPrint()
		n++
	}
	fmt.Printf("Loaded %d blueprints\n\n", len(bps))

	for _, bp := range bps {
		fmt.Printf("** Blueprint %d **\n\n", bp.number)
		for min := 1; min <= 24; min++ {
			fmt.Printf("== Minute %d\n", min)
			// Choose what (if anything) to build
			queue := bp.Build()

			// Collect resources
			bp.Collect()

			// Collect production
			bp.CollectRobots(queue)

			fmt.Println()
		}
		fmt.Println()
		break
	}
}

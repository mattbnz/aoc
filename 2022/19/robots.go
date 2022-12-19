// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 19, Puzzle 1.
// Not Enough Minerals. Robots to the rescue.

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Plural(base string, n int) string {
	if n > 1 {
		base += "s"
	}
	return base
}
func InvPlural(base string, n int) string {
	if n == 1 {
		base += "s"
	}
	return base
}

type Resource int

const (
	NOTHING Resource = iota
	ORE
	CLAY
	OBSIDIAN
	GEODE
	EVERYTHING
)

func (r Resource) String() string {
	return []string{"NOTHING", "ore", "clay", "obsidian", "geode", "EVERYTHING"}[r]
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

type Brain interface {
	BuildChoice(bp *Blueprint) Resource
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
	for r := ORE; r < EVERYTHING; r++ {
		count := bp.robots[r]
		if count == 0 {
			continue
		}
		bp.inventory[r] += count
		fmt.Printf("%d %s-%sing %s %s %d %s; you now have %d %s.\n", count, r, r.Action(),
			Plural("robot", count), InvPlural(r.Action(), count), count, r, bp.inventory[r], r)
	}
}

func (bp *Blueprint) CanBuild(r Resource) bool {
	for nr, nq := range bp.robotCost[r] {
		if bp.inventory[nr] < nq {
			return false
		}
	}
	return true
}

func Consume(inv ResourceData, cost Cost) ResourceData {
	newInv := ResourceData{}
	for nr, q := range inv {
		newInv[nr] = q - cost[nr]
	}
	return newInv
}

func (bp *Blueprint) consume(r Resource) {
	bp.inventory = Consume(bp.inventory, bp.robotCost[r])
	fmt.Printf("Spend %s to start building a %s-%sing robot.\n", bp.robotCost[r], r, r.Action())
}

func (bp *Blueprint) Build(brain Brain) Resource {
	b := brain.BuildChoice(bp)
	if b != NOTHING {
		bp.consume(b)
	}
	return b
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

type BabyBrain struct {
	wants ResourceData
}

func (b BabyBrain) Want(r Resource, bp *Blueprint) bool {
	if bp.robots[r] < b.wants[r] {
		//fmt.Println("want ", r)
		return true
	}
	//fmt.Println("don't want ", r, bp.robots[r], b.wants[r])
	return false
}

// how long until we can build r
func (b BabyBrain) BuildLatency(r Resource, bp *Blueprint) int {
	need := bp.robotCost[r]
	l := 0
	for nr, nq := range need {
		t := 0
		if bp.inventory[nr] < nq {
			if bp.robots[nr] > 0 {
				t = int(math.Ceil(float64(nq) / float64(bp.robots[nr])))
			} else {
				t = b.BuildLatency(nr, bp) + 1 + nq
			}
		}
		l = Max(l, t)
	}
	return l
}

func (b BabyBrain) BuildChoice(bp *Blueprint) Resource {
	delay, build := b.RobotDelay(GEODE, bp.inventory, bp.robots, bp)
	if delay > 0 {
		return NOTHING
	}
	return build
}

// how long until we have enough of this resource to build the requested item?
func (b BabyBrain) InventoryDelay(need int, have int, robots int) int {
	qo := need - have
	if qo <= 0 {
		return 0
	}
	if robots > 0 {
		return int(math.Ceil(float64(qo) / float64(robots)))
	}
	return math.MaxInt
}

func (b BabyBrain) RobotDelay(r Resource, inv ResourceData, robots ResourceData, bp *Blueprint) (int, Resource) {
	if bp.CanBuild(r) {
		return 0, r
	}
	if r == ORE {
		cost := int(math.Ceil(float64(bp.robotCost[ORE][ORE]) / float64(robots[ORE])))
		//fmt.Println(r, " robot will take ", cost, " minutes to build")
		return cost, ORE
	}

	need := bp.robotCost[r]

	bd := ResourceData{}
	maxDelay := 0
	for nr := r - 1; nr > NOTHING; nr-- {
		// Work out how long until we have enough of this resource
		bd[nr] = b.InventoryDelay(need[nr], inv[nr], bp.robots[nr])
		if bd[nr] == 0 {
			continue // have enough, easy!
		}
		//fmt.Println("sufficient ", nr, " is ", bd[nr], " minutes away")
		// If we'll never have enough of this resource, then we need to build it first
		if bd[nr] == math.MaxInt {
			delay, build := b.RobotDelay(nr, inv, robots, bp)
			//fmt.Println("To build ", r, " we need ", nr, " so we need to build ", build, " in ", delay, " minutes first!")
			return delay, build
		}
		maxDelay = Max(maxDelay, bd[nr])
	}

	// Given the delay above, would it be quicker to build another robot to make this item quicker?
	for nr := r - 1; nr > NOTHING; nr-- {
		if bd[nr] <= 2 {
			continue // have enough, easy!
		}
		buildDelay, build := b.RobotDelay(nr, inv, robots, bp)
		newInv := Consume(inv, bp.robotCost[build])
		newDelay := b.InventoryDelay(need[nr], newInv[nr], bp.robots[nr]+1)
		if buildDelay+newDelay < bd[nr] && buildDelay+newDelay < maxDelay {
			//fmt.Println("Building ", build, " will speed up ", nr, " to build ", r, " by ", bd[nr]-buildDelay-newDelay)
			return 0, build
		}
	}

	// If we get to here, there's nothing useful for us to build, so just wait...
	return maxDelay, NOTHING
}

func NewBabyBrain(bp *Blueprint) BabyBrain {
	bb := BabyBrain{}
	bb.wants = ResourceData{
		ORE:      1,
		CLAY:     3,
		OBSIDIAN: 1,
	}
	return bb
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

	sum := 0
	for _, bp := range bps {
		fmt.Printf("** Blueprint %d **\n\n", bp.number)
		brain := NewBabyBrain(&bp)
		for min := 1; min <= 24; min++ {
			fmt.Printf("== Minute %d ==\n", min)
			// Choose what (if anything) to build
			queue := bp.Build(brain)

			// Collect resources
			bp.Collect()

			// Collect production
			if queue != NOTHING {
				bp.robots[queue] += 1
				fmt.Printf("The new %s-%sing robot is ready; you now have %d of them.\n", queue, queue.Action(), bp.robots[queue])
			}
			fmt.Println()
		}
		fmt.Println()
		quality := bp.number * bp.inventory[GEODE]
		fmt.Printf("Blueprint %d: Quality Level = %d\n", bp.number, quality)
		fmt.Println()
		sum += quality
	}
	fmt.Println()
	fmt.Println(sum)
}

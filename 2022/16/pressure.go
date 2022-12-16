// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 16, Puzzle 1.
// Proboscidea Volcanium - pressure release.

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

type Valve struct {
	name     string
	flowRate int

	paths []*Valve
}

func (v Valve) String() string {
	names := []string{}
	for _, o := range v.paths {
		names = append(names, o.name)
	}
	return fmt.Sprintf("Valve %s; flow rate=%d; leads to %s", v.name, v.flowRate, strings.Join(names, ", "))
}
func GetOrCreateValve(name string) *Valve {
	v := Valves[name]
	if v != nil {
		return v
	}
	v = &Valve{name: name, flowRate: -1}
	Valves[v.name] = v
	return v
}

func GetOrCreateValves(nameList string) []*Valve {
	rv := []*Valve{}
	for _, name := range strings.Split(nameList, ", ") {
		rv = append(rv, GetOrCreateValve(name))
	}
	return rv
}

func ShowDot(save bool) {
	file, err := ioutil.TempFile("", "dot")
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString("digraph {\n")
	for _, v := range Valves {
		file.WriteString(fmt.Sprintf("%s [label=\"%s flow=%d\"]\n", v.name, v.name, v.flowRate))
		for _, ov := range v.paths {
			file.WriteString(fmt.Sprintf("%s -> %s\n", v.name, ov.name))
		}
		file.WriteString("\n")
	}
	file.WriteString("}\n")

	cmd := exec.Command("dot", "-Tpng", "-o", file.Name()+".png", file.Name())
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Couldn't render %s: %v", file.Name(), err)
		return
	}
	//os.Remove(file.Name())
	if save {
		log.Print("Image is at ", file.Name()+".png")
	} else {
		cmd = exec.Command("eog", file.Name()+".png")
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Couldn't display %s: %v", file.Name()+".png", err)
			return
		}
		os.Remove(file.Name() + ".png")
	}
}

var INPUT_RE = regexp.MustCompile(`Valve (.*?) has flow rate=(\d+); tunnel.*? lead.*? to valve.*? (.*)$`)

var Valves map[string]*Valve

func main() {
	s := bufio.NewScanner(os.Stdin)

	Valves = map[string]*Valve{}

	for s.Scan() {
		m := INPUT_RE.FindStringSubmatch(s.Text())
		if len(m) != 4 {
			log.Fatalf("Couldn't parse: %s", s.Text())
		}
		v := GetOrCreateValve(m[1])
		if v.flowRate != -1 {
			log.Fatal("multiple flow definitions of ", v)
		}
		v.flowRate = Int(m[2])
		v.paths = GetOrCreateValves(m[3])
		fmt.Println(v)
	}
	//ShowDot(false)
}

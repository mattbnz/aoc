// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 7, Puzzle 2.
// Find the directory to delete.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Dir struct {
	Name   string
	Parent *Dir
	Dirs   map[string]*Dir
	Files  map[string]*File
}

func (d *Dir) String() string {
	return d.Name // TODO: walk paarents
}

func (d *Dir) Size() int {
	size := 0
	for _, sd := range d.Dirs {
		size += sd.Size()
	}
	for _, f := range d.Files {
		size += f.Size
	}
	return size
}

// depth-first walk with callback
func (d *Dir) Walk(f func(dd *Dir)) {
	for _, sd := range d.Dirs {
		sd.Walk(f)
	}
	f(d)
}

type File struct {
	Size int
	Name string
}

func Pad(lvl int) string {
	return strings.Repeat(" ", lvl)
}

func PrintDirs(d *Dir, lvl int) {
	fmt.Printf("%s- %s (dir, size=%d)\n", Pad(lvl), d.Name, d.Size())
	for _, sd := range d.Dirs {
		PrintDirs(sd, lvl+1)
	}
	for _, f := range d.Files {
		fmt.Printf("%s- %s (file, size=%d)\n", Pad(lvl+1), f.Name, f.Size)
	}
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	root := &Dir{Name: "/", Parent: nil, Dirs: map[string]*Dir{}, Files: map[string]*File{}}
	curDir := root

	for s.Scan() {
		if s.Text() == "$ cd /" {
			// Special case cd to root for ease.
			curDir = root
			continue
		}
		if curDir == nil {
			log.Fatalf("Error: cannot '%s' without cd", s.Text())
		}
		if strings.HasPrefix(s.Text(), "$ cd") {
			// Change dir
			name := s.Text()[5:]
			if name == ".." {
				if curDir.Parent != nil {
					curDir = curDir.Parent
				} else {
					log.Fatalf("Error: cannot '%s' from %v", s.Text(), curDir)
				}
			} else {
				dir, exists := curDir.Dirs[name]
				if !exists {
					log.Fatalf("Error: cannot '%s' from %s, path does not exist", s.Text(), curDir)
				}
				curDir = dir
			}
			continue
		}
		if strings.HasPrefix(s.Text(), "$") {
			continue // skip other commands like ls, not useful
		}
		if strings.HasPrefix(s.Text(), "dir") {
			name := s.Text()[4:]
			log.Printf("Adding dir %s to %s", name, curDir)
			curDir.Dirs[name] = &Dir{Name: name, Parent: curDir, Dirs: map[string]*Dir{}, Files: map[string]*File{}}
		} else {
			sizeS, name, ok := strings.Cut(s.Text(), " ")
			if !ok {
				log.Fatalf("Error: cannot parse file: %s", s.Text())
			}
			size, err := strconv.Atoi(sizeS)
			if err != nil {
				log.Fatalf("Error: bad file size for %s: %v", s.Text(), err)
			}
			log.Printf("Adding file %s to %s", name, curDir)
			curDir.Files[name] = &File{Size: size, Name: name}
		}
	}

	PrintDirs(root, 0)

	const FSSIZE = 70000000
	const HEADROOM = 30000000
	free := FSSIZE - root.Size()
	need := HEADROOM - free

	candidates := []*Dir{}
	root.Walk(func(d *Dir) {
		if d.Size() >= need {
			candidates = append(candidates, d)
		}
	})
	fmt.Println(candidates)

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Size() < candidates[j].Size()
	})
	fmt.Println(candidates)
	fmt.Println(candidates[0].Name, candidates[0].Size())
}

// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 17, Puzzle 1
// Conway Cube Initialization.

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "reflect"
    "sort"
)

var STATES = map[int]string{-1:"_", 0:".", 1:"#"}

type COL map[int]int
type ROW map[int]COL
type Z map[int]ROW

// Return Cube State (1, 0) or -1 if the requested cube doesn't exist.
func CubeState(z Z, zk, rk, ck int) int {
    r, zE := z[zk]
    if !zE {
        return -1
    }
    c, rE := r[rk]
    if !rE {
        return -1
    }
    state, cE := c[ck]
    if !cE {
        return -1
    }
    return state
}

// Return count of active neighbors for a cube; -1 if cube doesn't exist
func CubeActiveNeighbors(z Z, zB, rB, cB int) int {
    sum := 0
    for zd:=-1; zd<=1; zd++ {
        zk := zB + zd
        for rd:=-1; rd<=1; rd++ {
            rk := rB + rd
            for cd:=-1; cd<=1; cd++ {
                ck := cB + cd
                if cB==ck && rB==rk && zB==zk {
                    //fmt.Printf(" - skip at (%d, %d, %d)\n", zk, rk, ck)
                    continue
                }
                s := CubeState(z, zk, rk, ck)
                if s == 1 {
                    //fmt.Printf(" - neighbour at (%d, %d, %d)\n", zk, rk, ck)
                    sum++
                //} else {
                //    fmt.Printf(" - empty at (%d, %d, %d)\n", zk, rk, ck)
                }
            }
        }
    }
    return sum
}

func ActivateCubes(z Z) (int, Z) {
    zkeys := SortedKeys(z)
    rkeys := SortedKeys(z[zkeys[0]])
    ckeys := SortedKeys(z[zkeys[0]][rkeys[0]])

    sum := 0
    newz := Z{}
    for zk := zkeys[0]-1; zk<=zkeys[len(zkeys)-1]+1; zk++ {
        newrows := ROW{}
        for rk := rkeys[0]-1; rk<=rkeys[len(rkeys)-1]+1; rk++ {
            newcols := COL{}
            for ck := ckeys[0]-1; ck<=ckeys[len(ckeys)-1]+1; ck++ {
                state := CubeState(z, zk, rk, ck)
                actives := CubeActiveNeighbors(z, zk, rk, ck)
                newstate := 0
                if state == 1 {
                    if actives==2 || actives==3 {
                        newstate = 1
                    }
                } else if actives == 3 {
                    newstate = 1
                }
                newcols[ck] = newstate
                if newcols[ck] == 1 {
                    sum++
                }
                //fmt.Printf("(%d,%d,%d) is %s and has %d neighbors => %s (%d)\n",
                //            zk, rk, ck, STATES[state], actives,
                //            STATES[newcols[ck]], sum)
            }
            newrows[rk] = newcols
        }
        newz[zk] = newrows
    }
    return sum, newz
}

func SortedKeys(in interface{}) []int {
    m := reflect.ValueOf(in)
    keys := []int{}
    for _, v := range m.MapKeys() {
        keys = append(keys, v.Interface().(int))
    }
    sort.Ints(keys)
    return keys
}

func PrintCube(z Z) {
    for _, zk := range SortedKeys(z) {
        fmt.Printf("z=%d\n", zk)
        for _, rk := range SortedKeys(z[zk]) {
            for _, ck := range SortedKeys(z[zk][rk]) {
                s, _ := STATES[z[zk][rk][ck]]
                fmt.Printf("%s", s)
            }
            fmt.Printf("\n")
        }
        fmt.Printf("\n")
    }
}

func main() {
    s := bufio.NewScanner(os.Stdin)
    z := Z{}  // z,x,y; 0=./Inactive, 1=#/Active
    rows := ROW{}
    r := 0
    for s.Scan() {
        cols := COL{}
        for c, r := range s.Text() {
            if r == '.' {
                cols[c] = 0
            } else if r == '#' {
                cols[c] = 1
            } else {
                log.Fatal("Bad cube state: ", s.Text())
            }
        }
        rows[r] = cols
        r++
    }
    z[0] = rows

    PrintCube(z)

    // Initialize with 6 cycles
    active := 0
    for i:=0; i<6; i++ {
        active, z = ActivateCubes(z)
        //fmt.Printf("%d active after iteration %d\n", active, i+1)
        //PrintCube(z)
    }
    fmt.Println(active)
}

// Copyright (C) 2020 Matt Brown

// Advent of Code 2020 - Day 17, Puzzle 2
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
type W map[int]Z

// Return Cube State (1, 0) or -1 if the requested cube doesn't exist.
func CubeState(w W, wk, zk, rk, ck int) int {
    z, wE := w[wk]
    if !wE {
        return -1
    }
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
func CubeActiveNeighbors(w W, wB, zB, rB, cB int) int {
    sum := 0
    for wd:=-1; wd<=1; wd++ {
        wk := wB + wd
        for zd:=-1; zd<=1; zd++ {
            zk := zB + zd
            for rd:=-1; rd<=1; rd++ {
                rk := rB + rd
                for cd:=-1; cd<=1; cd++ {
                    ck := cB + cd
                    if cB==ck && rB==rk && zB==zk && wB==wk {
                        //fmt.Printf(" - skip at (%d, %d, %d)\n", zk, rk, ck)
                        continue
                    }
                    s := CubeState(w, wk, zk, rk, ck)
                    if s == 1 {
                        //fmt.Printf(" - neighbour at (%d, %d, %d)\n", zk, rk, ck)
                        sum++
                    //} else {
                    //    fmt.Printf(" - empty at (%d, %d, %d)\n", zk, rk, ck)
                    }
                }
            }
        }
    }
    return sum
}

func ActivateCubes(w W) (int, W) {
    wkeys := SortedKeys(w)
    zkeys := SortedKeys(w[wkeys[0]])
    rkeys := SortedKeys(w[wkeys[0]][zkeys[0]])
    ckeys := SortedKeys(w[wkeys[0]][zkeys[0]][rkeys[0]])

    sum := 0
    neww := W{}
    for wk := wkeys[0]-1; wk<=wkeys[len(wkeys)-1]+1; wk++ {
        newz := Z{}
        for zk := zkeys[0]-1; zk<=zkeys[len(zkeys)-1]+1; zk++ {
            newrows := ROW{}
            for rk := rkeys[0]-1; rk<=rkeys[len(rkeys)-1]+1; rk++ {
                newcols := COL{}
                for ck := ckeys[0]-1; ck<=ckeys[len(ckeys)-1]+1; ck++ {
                    state := CubeState(w, wk, zk, rk, ck)
                    actives := CubeActiveNeighbors(w, wk, zk, rk, ck)
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
        neww[wk] = newz
    }
    return sum, neww
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

func PrintCube(w W) {
    for _, wk := range SortedKeys(w) {
        for _, zk := range SortedKeys(w[wk]) {
            fmt.Printf("z=%d, w=%d\n", zk, wk)
            for _, rk := range SortedKeys(w[wk][zk]) {
                for _, ck := range SortedKeys(w[wk][zk][rk]) {
                    s, _ := STATES[w[wk][zk][rk][ck]]
                    fmt.Printf("%s", s)
                }
                fmt.Printf("\n")
            }
            fmt.Printf("\n")
        }
    }
}

func main() {
    s := bufio.NewScanner(os.Stdin)
    w := W{}  // w,z,x,y; 0=./Inactive, 1=#/Active
    z := Z{}
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
    w[0] = z

    PrintCube(w)

    // Initialize with 6 cycles
    active := 0
    for i:=0; i<6; i++ {
        active, w = ActivateCubes(w)
        fmt.Printf("%d active after iteration %d\n", active, i+1)
        PrintCube(w)
    }
    fmt.Println(active)
}

package day19

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

type Part map[string]int

func (p Part) String() string {
	jb, _ := json.Marshal(p)
	return string(jb)
}

func NewPart(s string) (rv Part) {
	rv = make(Part)
	for _, attrVal := range strings.Split(strings.Trim(s, "{}"), ",") {
		attr, valS, found := strings.Cut(attrVal, "=")
		if !found {
			glog.Fatalf("Bad Part: %s", s)
		}
		val, err := strconv.Atoi(valS)
		if err != nil {
			glog.Fatalf("Bad Part Attribute: %s", attrVal)
		}
		rv[attr] = val
	}
	return
}

type Rule struct {
	attr string
	op   string
	val  int
	dest string

	DestW          *Workflow
	Size           int64
	RemainingSpace map[string]int64
}

func (r Rule) IsDefault() bool {
	return r.dest != "" && r.attr == "" && r.op == "" && r.val == 0
}

func (r Rule) String() string {
	if r.IsDefault() {
		return r.dest
	}
	return fmt.Sprintf("%s%s%d:%s", r.attr, r.op, r.val, r.dest)
}

var ruleRegexp = regexp.MustCompile(`^([a-z]+)([><])(\d+):([a-zAR]+)$`)

func NewRule(s string) (rv *Rule) {
	rv = &Rule{}
	rv.RemainingSpace = make(map[string]int64)
	if !strings.Contains(s, ":") {
		rv.dest = s
		return
	}

	m := ruleRegexp.FindStringSubmatch(s)
	if m == nil {
		glog.Fatalf("Bad Rule: %s", s)
	}
	rv.attr = m[1]
	rv.op = m[2]
	val, err := strconv.Atoi(m[3])
	if err != nil {
		glog.Fatalf("Bad Rule (%v): %s", err, s)
	}
	rv.val = val
	rv.dest = m[4]

	return
}

type Workflow struct {
	Name string

	Rules   []*Rule
	Default string

	Queue []*Part

	Result []*big.Int
}

func (w *Workflow) Append(p *Part) {
	w.Queue = append(w.Queue, p)
}

func NewWorkflow(s string) (rv Workflow) {
	b := strings.Index(s, "{")
	if b == -1 {
		glog.Fatalf("bad Workflow: %s", s)
	}
	rv.Name = s[:b]
	for _, r := range strings.Split(strings.Trim(s[b:], "{}"), ",") {
		rv.Rules = append(rv.Rules, NewRule(r))
	}
	for n, r := range rv.Rules {
		last := n == len(rv.Rules)-1
		if r.IsDefault() && !last {
			glog.Fatalf("Got default rule (%s) at non-final position %d of %s", r, n, rv.Name)
		}
		if last && !r.IsDefault() {
			glog.Fatalf("Final rule (%s) of %s is not a default", r, n, rv.Name)
		}
	}
	return
}

type Heap struct {
	Workflows map[string]*Workflow
}

func NewHeap(filename string) (h Heap, err error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	h.Workflows = make(map[string]*Workflow)
	parseWorkflows := true
	var in *Workflow
	for s.Scan() {
		if s.Text() == "" {
			parseWorkflows = false
			if i, ok := h.Workflows["in"]; ok {
				in = i
				continue
			}
			glog.Fatal("in workflow not found!")
		}

		if parseWorkflows {
			w := NewWorkflow(s.Text())
			h.Workflows[w.Name] = &w
		} else {
			p := NewPart(s.Text())
			in.Append(&p)
		}
	}
	glog.V(1).Infof("Loaded %d workflows", len(h.Workflows))
	glog.V(1).Infof("Loaded %d parts", len(in.Queue))
	h.Workflows["A"] = &Workflow{Name: "A"}
	h.Workflows["R"] = &Workflow{Name: "R"}
	return
}

func (h *Heap) SortParts() {
	i := 0
	for {
		processed := 0
		for _, w := range h.Workflows {
			if len(w.Queue) == 0 || w.Name == "A" || w.Name == "R" {
				continue // don't process empty, or final queues.
			}
			//q := append([]*Part{}, w.Queue...)
			q := w.Queue
			w.Queue = []*Part{}
			for _, p := range q {
				dest := ""
				for _, r := range w.Rules {
					if r.IsDefault() {
						dest = r.dest
						break
					}
					match := false
					switch r.op {
					case ">":
						match = (*p)[r.attr] > r.val
					case "<":
						match = (*p)[r.attr] < r.val
					default:
						glog.Fatalf("bad rule (%s) in %s, cannot match %s", r, w, p)
					}
					if match {
						dest = r.dest
						break
					}
				}
				if dest == "" {
					glog.Fatalf("did not find match (even default!) for %s in %s", p, w)
				}
				h.Workflows[dest].Append(p)
				processed++
			}
		}
		glog.V(2).Infof("Iteration % 4d: Processed %d parts", i, processed)
		if processed == 0 {
			break
		}
		i++
	}
}

func (h *Heap) Sum(workflow string) (rv int) {
	for _, p := range h.Workflows[workflow].Queue {
		for _, v := range *p {
			rv += v
		}
	}
	return
}

// Node represents a branch from a workflow rule match.
type Node struct {
	// Workflow and rule index that this node was created from
	W *Workflow
	R int
}

func (h *Heap) BuildGraph() {
	for _, w := range h.Workflows {
		for _, r := range w.Rules {
			r.DestW = h.Workflows[r.dest]
		}
	}
	for _, w := range h.Workflows {
		rules := []Rule{}
		for _, r := range w.Rules {
			rules = append(rules, *r)
		}
		rSizes, _ := calcAttrSize(rules, false)
		aSizes := map[string]int64{"x": 4000, "m": 4000, "a": 4000, "s": 4000}
		for n, size := range rSizes {
			w.Rules[n].Size = size
			aSizes[w.Rules[n].attr] -= size
			for k, v := range aSizes {
				w.Rules[n].RemainingSpace[k] = v
			}
		}
		/*if len(w.Rules) > 0 {
			defIdx := len(w.Rules) - 1
			for a, size := range aSizes {
				w.Rules[defIdx].RemainingSpace[a] = 4000 - size
			}
		}*/
		/*		for _, a := range []string{"x", "m", "a", "s"} {
					aR := []*Rule{}
					for _, r := range w.Rules {
						if r.attr != a {
							continue
						}
						aR = append(aR, r)
					}
					sort.Slice(aR, func(i, j int) bool {
						return aR[i].val < aR[j].val
					})
					aSize := int64(4000)
					for n, r := range aR {
						switch r.op {
						case ">":
							max := 4000
							if n < len(aR)-1 {
								max = aR[n+1].val
							}
							r.Size = int64(max - (r.val + 1))
						case "<":
							bottom := 0
							if n > 0 {
								bottom = aR[n-1].val
							}
							r.Size = int64(r.val - 1 - bottom)
						default:
							glog.Fatalf("bad rule (%s) in %s", r, w)
						}
						aSize -= r.Size
					}
					defaultSize += aSize
				}
				if len(w.Rules) > 1 {
					w.Rules[len(w.Rules)-1].Size = defaultSize
				} else if len(w.Rules) == 1 {
					w.Rules[0].Size = 4000
				}*/
	}
}

func calcAttrSize(rules []Rule, useLimits bool) (rSizes []int64, aSizes map[string]int64) {
	aSizes = make(map[string]int64)

	for _, a := range []string{"x", "m", "a", "s"} {
		clamp := 0
		aR := []*Rule{}
		for n, r := range rules {
			if useLimits && r.RemainingSpace[a] != 0 {
				if clamp == 0 {
					clamp = int(r.RemainingSpace[a])
				} else {
					//glog.Infof("multiple defaults for " + a)
					clamp = Min(clamp, int(r.RemainingSpace[a]))
				}
			}
			if r.attr != a {
				continue
			}
			aR = append(aR, &rules[n])
		}
		sort.Slice(aR, func(i, j int) bool {
			return aR[i].val < aR[j].val
		})
		sum := int64(0)
		for n, r := range aR {
			switch r.op {
			case ">":
				max := 4000
				if n < len(aR)-1 {
					max = aR[n+1].val
				}
				r.Size = int64(max - (r.val))
			case "<":
				bottom := 0
				if n > 0 {
					bottom = aR[n-1].val
				}
				r.Size = int64(r.val - 1 - bottom)
			default:
				glog.Fatalf("bad rule (%s)", r)
			}
			sum += r.Size
		}
		if clamp == 0 {
			aSizes[a] = sum
		} else if sum == 0 {
			aSizes[a] = int64(clamp)
		} else {
			aSizes[a] = int64(Min(int(sum), clamp))
		}
	}

	rSizes = make([]int64, len(rules))
	for n, r := range rules {
		rSizes[n] = r.Size
	}
	return
}

type Path struct {
	N string
	R *Rule
}

func (h *Heap) Combinations() *big.Int {
	in := h.Workflows["in"]
	in.walk([]Path{})

	sum := big.NewInt(0)
	A := h.Workflows["A"]
	for _, r := range A.Result {
		sum = sum.Add(sum, r)
		glog.Infof("+ %s", r)
	}
	return sum
}

var Neg1 = big.NewInt(-1)

func (w *Workflow) walk(path []Path) *big.Int {
	rules := []Rule{}
	for n := range path {
		rules = append(rules, *path[n].R)
	}
	_, aSizes := calcAttrSize(rules, true)

	if len(w.Rules) == 0 {
		product := big.NewInt(1)
		dbg := []string{}
		for a, s := range aSizes {
			if s == 0 {
				s = 4000
			}
			n := big.NewInt(s)
			product = product.Mul(product, n)
			dbg = append(dbg, fmt.Sprintf("%s:%d", a, s))
		}
		path = append(path, Path{N: w.Name})
		glog.Infof("         %s", strings.Join(dbg, " x "))
		LogPath(w, path, product)
		w.Result = append(w.Result, product)
		//glog.Infof("%#v", rSizes)
		//glog.Infof("%#v", aSizes)
		return Neg1
	}

	rv := big.NewInt(0)
	for _, r := range w.Rules {
		below := r.DestW.walk(append(path, Path{w.Name, r}))
		if below.Cmp(Neg1) == 0 {
			rv = rv.Add(rv, big.NewInt(r.Size))
		} else {
			sz := big.NewInt(r.Size)
			rv = rv.Add(rv, sz.Mul(sz, below))
		}
	}
	LogPath(w, path, rv)
	return rv
}

func LogPath(w *Workflow, path []Path, n *big.Int) {
	s := []string{}
	for _, p := range path {
		s = append(s, fmt.Sprintf("%s,%s", p.N, p.R))
	}
	glog.Infof("% 5s: %s == %d", w.Name, strings.Join(s, " => "), n)
}

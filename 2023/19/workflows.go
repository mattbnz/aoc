package day19

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
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

func NewRule(s string) (rv Rule) {
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

	Rules   []Rule
	Default string

	Queue []*Part
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

// Copyright (C) 2022 Matt Brown

// Advent of Code 2022 - Day 22, Puzzle 2.
// Monkey Map.  Path following on a cube!.

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

func Max(a, b int) int {
	if a == -1 {
		return b
	}
	if b == -1 {
		return a
	}
	if a > b {
		return a
	}
	return b
}

func Int(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("%s is not an int: %v", s, err)
	}
	return v
}

func GCD(a, b int) int {
	rv := 1
	for i := 1; i <= a && i <= b; i++ {
		if a%i == 0 && b%i == 0 {
			rv = i
		}
	}
	return rv
}

// Basic position type, used in three ways
// 1) To reference cells in the input (called "Absolute position")
// 2) To reference cells on a side (called "Relative position")
// 3) To reference sides in the input (called "Side position")
type Pos struct {
	row, col int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d,%d", p.row, p.col)
}

// What can go in each position
type Content int

const (
	VOID int = iota
	OPEN
	WALL
)

// Directions of movement
type Direction int

const (
	D_RIGHT Direction = iota
	D_DOWN
	D_LEFT
	D_UP
)

func (d Direction) String() string {
	return DIR_S[d]
}

var DIR_S = map[Direction]string{
	D_RIGHT: ">",
	D_DOWN:  "v",
	D_LEFT:  "<",
	D_UP:    "^",
}

const STR_LEFT = "L"
const STR_RIGHT = "R"

var TURN_MAP = map[Direction]map[string]Direction{
	D_RIGHT: {STR_LEFT: D_UP, STR_RIGHT: D_DOWN},
	D_DOWN:  {STR_LEFT: D_RIGHT, STR_RIGHT: D_LEFT},
	D_LEFT:  {STR_LEFT: D_DOWN, STR_RIGHT: D_UP},
	D_UP:    {STR_LEFT: D_LEFT, STR_RIGHT: D_RIGHT},
}

// Rotations that can be made when loading the cube
type Rotation int

const (
	R_NONE Rotation = iota
	R_ANTI_CW
	R_CW
	R_FLIPR
	R_FLIPC
)

var R_STR = map[Rotation]string{
	R_NONE:    "unmodified",
	R_ANTI_CW: "rotate anti-clockwise",
	R_CW:      "rotate clockwise",
	R_FLIPR:   "flip rows",
	R_FLIPC:   "flip cols",
}

// how side adjacencies move when rotating
var SIDE_MAP = map[Direction]map[Rotation]Direction{
	D_UP:    {R_ANTI_CW: D_LEFT, R_CW: D_RIGHT},
	D_LEFT:  {R_ANTI_CW: D_DOWN, R_CW: D_UP, R_FLIPC: D_RIGHT},
	D_DOWN:  {R_ANTI_CW: D_LEFT, R_CW: D_RIGHT, R_FLIPC: D_LEFT},
	D_RIGHT: {R_ANTI_CW: D_UP, R_CW: D_DOWN},
}

func (r Rotation) String() string {
	return R_STR[r]
}

// A side of the cube
type Side struct {
	SidePos Pos // Side position in the input grid

	// zero based relative co-ordinates on this size, use Abs() to get
	// back to input co-ords.
	c    map[Pos]int
	Size int

	Orientation CubeSide
	Transform   Rotation

	next map[Direction]*Side
}

func (s Side) String() string {
	return fmt.Sprintf("%s side (from %s %s)", s.Orientation, s.SidePos, s.Transform)
}

func (s Side) Print() {
	fmt.Println(s)
	for row := 0; row < s.Size; row++ {
		for col := 0; col < s.Size; col++ {
			c := s.C(Pos{row, col})
			if c == VOID {
				fmt.Printf(" ")
			} else if c == OPEN {
				fmt.Printf(".")
			} else if c == WALL {
				fmt.Printf("#")
			}
		}
		fmt.Println()
	}
}

func (s Side) C(p Pos) int {
	return s.c[p] // default val if missing == VOID
}

// Returns the absolute co-ordinates of this sides top-left point
func (s Side) AbsBase() Pos {
	r, c := s.SidePos.row, s.SidePos.col
	gr := (r * s.Size) + 1
	gc := (c * s.Size) + 1
	return Pos{gr, gc}
}

func (s Side) Abs(rel Pos) Pos {
	// TODO: Needs to handle transformed sides...
	base := s.AbsBase()
	return Pos{base.row + rel.row, base.col + rel.col}
}

func (s *Side) Rotate(rotation Rotation) {
	if rotation == R_NONE {
		return
	}
	fmt.Printf("Rotating %s %s\n", s, rotation)
	s.Print()
	sc := map[Pos]int{}
	if rotation == R_ANTI_CW {
		for p, what := range s.c {
			sc[Pos{s.Size - 1 - p.col, p.row}] = what
		}
	} else if rotation == R_CW {
		for p, what := range s.c {
			sc[Pos{p.col, s.Size - (s.Size - 1 - p.row)}] = what
		}
	} else if rotation == R_FLIPC {
		for p, what := range s.c {
			sc[Pos{p.row, s.Size - 1 - p.col}] = what
		}
	}
	ns := Side{SidePos: s.SidePos, c: sc, Size: s.Size, Transform: rotation}
	ns.next = map[Direction]*Side{}
	for d, os := range s.next {
		if nd, found := SIDE_MAP[d][rotation]; found {
			ns.next[nd] = os
		}
	}
	ns.Print()
	*s = ns
}

// Converts c to relative co-ords and returns a side of them
func NewSide(origin Pos, c map[Pos]int, size int) Side {
	s := Side{SidePos: origin, Size: size}
	s.next = map[Direction]*Side{}

	base := s.AbsBase()
	relc := map[Pos]int{}
	for p := range c {
		np := Pos{p.row - base.row, p.col - base.col}
		relc[np] = c[p]
	}
	s.c = relc
	s.Print()
	return s
}

// Descriptions of the side of a cube
type CubeSide int

var CubeSide_S = map[CubeSide]string{
	CS_TOP:    "Top",
	CS_BOTTOM: "Bottom",
	CS_LEFT:   "Left",
	CS_RIGHT:  "Right",
	CS_FRONT:  "Front",
	CS_BACK:   "Back",
}

const (
	CS_TOP CubeSide = iota
	CS_BOTTOM
	CS_LEFT
	CS_RIGHT
	CS_FRONT
	CS_BACK
)

func (sp CubeSide) String() string {
	return CubeSide_S[sp]
}

// A cube.
type Cube struct {
	// Stores sides in normalized "cross" form, with BACK at the top, BOTTOM on the farm left, TOP in the center.
	// e.g.                 BACK
	//       BOTTOM | LEFT | TOP | RIGHT
	//                      FRONT
	sides map[CubeSide]*Side

	Size int
}

func (c Cube) Print() {
	fmt.Println("Cube:")
	for _, side := range c.sides {
		side.Print()
	}
}

// Used to construct a cube from a 2D representation, by folding sides in
func (c *Cube) Fold(sides map[Pos]*Side, side Pos, position CubeSide, rotation Rotation) {
	s := sides[side]
	fmt.Printf("Trying to put %s into %s of cube\n", s, position)
	if _, taken := c.sides[position]; taken {
		c.Print()
		s.Print()
		fmt.Printf("%s is already populated, cannot put %s there!\n", position, side)
		return
	}

	remSides := map[Pos]*Side{}
	for p, s := range sides {
		if p != side {
			remSides[p] = s
		}
	}

	s.Rotate(rotation)
	s.Orientation = position
	c.sides[position] = s

	for dir := range DIR_S {
		fmt.Printf(" - (%s) Checking %s: ", position, dir)
		if _, present := s.next[dir]; present {
			fmt.Printf(" already present!\n")
			continue
		}
		// TODO: Need to normalize, transform!
		opos := OriginPos(side, dir /*I_MAP[position][dir]*/)
		if next, exists := remSides[opos]; exists {
			wrap := C_MAP[position][dir]
			fmt.Printf(" needs to fold %s (%s) to make %s side\n", next, wrap.Rotation, wrap.Side)
			c.Fold(remSides, opos, wrap.Side, wrap.Rotation)
			continue
		} /*else {
			log.Fatalf("Cannot find %s side for %")
		}*/
		fmt.Println()
	}
	return
}

// Spec of what happens when we wrap round a particular edge (assuming normalized form above)
type WrapSpec struct {
	Side      CubeSide
	Transform Pos // one of the T_ constants below for each of row,col
	Heading   Direction
	Rotation  Rotation
}

const T_MIN = 0           // use min row, col for position
const T_MAX = math.MaxInt // use max row, col for position
const T_MATCHR = -1       // use row value from current position
const T_MATCHC = -2       // use col value from current position
const T_MATCHR_I = -3     // use inverse row value from current position
const T_MATCHC_I = -4     // use inverse col value from current position

// Assumes normalized sides
var C_MAP = map[CubeSide]map[Direction]WrapSpec{
	CS_TOP: {
		D_UP:    {CS_BACK, Pos{T_MAX, T_MATCHC}, D_UP, R_NONE},
		D_LEFT:  {CS_LEFT, Pos{T_MATCHR, T_MAX}, D_LEFT, R_NONE},
		D_RIGHT: {CS_RIGHT, Pos{T_MATCHR, T_MIN}, D_RIGHT, R_NONE},
		D_DOWN:  {CS_FRONT, Pos{T_MIN, T_MATCHC}, D_DOWN, R_NONE},
	},
	CS_BOTTOM: {
		D_UP:    {CS_BACK, Pos{T_MIN, T_MATCHC}, D_DOWN, R_NONE},
		D_LEFT:  {CS_RIGHT, Pos{T_MATCHR, T_MAX}, D_LEFT, R_NONE},
		D_RIGHT: {CS_LEFT, Pos{T_MATCHR, T_MIN}, D_RIGHT, R_NONE},
		D_DOWN:  {CS_FRONT, Pos{T_MAX, T_MATCHC}, D_UP, R_NONE},
	},

	CS_FRONT: {
		D_UP:    {CS_TOP, Pos{T_MAX, T_MATCHC}, D_UP, R_NONE},
		D_LEFT:  {CS_LEFT, Pos{T_MAX, T_MATCHR_I}, D_UP, R_ANTI_CW},
		D_RIGHT: {CS_RIGHT, Pos{T_MAX, T_MATCHR}, D_UP, R_ANTI_CW},
		D_DOWN:  {CS_BOTTOM, Pos{T_MAX, T_MATCHC_I}, D_UP, R_FLIPR},
	},
	CS_BACK: {
		D_UP:    {CS_BOTTOM, Pos{T_MIN, T_MATCHC}, D_DOWN, R_FLIPR},
		D_LEFT:  {CS_LEFT, Pos{T_MIN, T_MATCHR}, D_DOWN, R_ANTI_CW},
		D_RIGHT: {CS_RIGHT, Pos{T_MIN, T_MATCHR_I}, D_DOWN, R_CW},
		D_DOWN:  {CS_TOP, Pos{T_MIN, T_MATCHC}, D_DOWN, R_NONE},
	},

	CS_LEFT: {
		D_UP:    {CS_BACK, Pos{T_MATCHC, T_MIN}, D_RIGHT, R_CW},
		D_LEFT:  {CS_BOTTOM, Pos{T_MATCHR, T_MIN}, D_RIGHT, R_NONE}, // bottom already rotated in this state
		D_RIGHT: {CS_TOP, Pos{T_MATCHR, T_MIN}, D_RIGHT, R_NONE},
		D_DOWN:  {CS_FRONT, Pos{T_MATCHC_I, T_MIN}, D_RIGHT, R_ANTI_CW},
	},
	CS_RIGHT: {
		D_UP:    {CS_BACK, Pos{T_MATCHC_I, T_MAX}, D_LEFT, R_ANTI_CW},
		D_LEFT:  {CS_TOP, Pos{T_MATCHR, T_MAX}, D_LEFT, R_NONE},
		D_RIGHT: {CS_BOTTOM, Pos{T_MATCHR, T_MAX}, D_LEFT, R_FLIPC},
		D_DOWN:  {CS_FRONT, Pos{T_MATCHC, T_MAX}, D_LEFT, R_CW},
	},
}

func (c Cube) C(p CubePos) int {
	return p.Side.C(p.RelPos) // default val if missing == VOID
}

func (c *Cube) Next(from CubePos) CubePos {
	p := Pos{}
	if from.Heading == D_RIGHT {
		p = Pos{from.RelPos.row, from.RelPos.col + 1}
	} else if from.Heading == D_DOWN {
		p = Pos{from.RelPos.row + 1, from.RelPos.col}
	} else if from.Heading == D_LEFT {
		p = Pos{from.RelPos.row, from.RelPos.col - 1}
	} else if from.Heading == D_UP {
		p = Pos{from.RelPos.row - 1, from.RelPos.col}
	}
	return CubePos{Side: from.Side, Heading: from.Heading, RelPos: p}
}

func (c *Cube) TransformPos(p Pos, t Pos) Pos {
	np := Pos{}
	//fmt.Println(p, t)
	if t.row == T_MIN {
		np.row = 0
	} else if t.row == T_MAX {
		np.row = c.Size - 1
	} else if t.row == T_MATCHC {
		np.row = p.col
	} else if t.row == T_MATCHR {
		np.row = p.row
	} else if t.row == T_MATCHC_I {
		np.row = c.Size - 1 - p.col
	} else {
		log.Fatalf("Unknown row transform from %s to %s", p, t)
	}

	if t.col == T_MIN {
		np.col = 0
	} else if t.col == T_MAX {
		np.col = c.Size - 1
	} else if t.col == T_MATCHC {
		np.col = p.col
	} else if t.col == T_MATCHC_I {
		np.col = c.Size - 1 - p.col
	} else if t.col == T_MATCHR {
		np.col = p.row
	} else if t.row == T_MATCHR_I {
		np.col = c.Size - 1 - p.row
	} else {
		log.Fatalf("Unknown col transform from %s to %s", p, t)
	}

	return np
}

func (c *Cube) Wrap(from CubePos) CubePos {
	next := C_MAP[from.Side.Orientation][from.Heading]
	fmt.Println(from.Side.Orientation, from.Heading, next.Side, next.Heading)
	return CubePos{c.sides[next.Side], c.TransformPos(from.RelPos, next.Transform), next.Heading}
}

func (c *Cube) Nav(from CubePos, steps int) CubePos {
	at := from
	for i := 0; i < steps; i++ {
		next := c.Next(at)
		if c.C(next) == OPEN {
			fmt.Println(at, " ==> ", next, " OK")
			at = next
			continue
		}
		if c.C(next) == WALL {
			fmt.Println(at, " ==> ", next, " BLOCKED")
			return at
		}
		if c.C(next) == VOID {
			fmt.Println(at, " is at the edge!")
			next = c.Wrap(at)
			if c.C(next) == OPEN {
				fmt.Println(at, " ==>", next, " OK")
				at = next
				//i++
				continue
			}
			if c.C(next) == WALL {
				fmt.Println(at, " ==>", next, " BLOCKED")
				return at
			}
		}
	}
	return at
}

// Position on the cube
type CubePos struct {
	Side    *Side
	RelPos  Pos // 1,1 based (aka relative), co-ordinates on Side
	Heading Direction
}

func (cp CubePos) String() string {
	return fmt.Sprintf("%s/%s going %s", cp.Side, cp.RelPos, cp.Heading)
}

// Returns an absolute (aka on input) coordinate matching fp
func (cp CubePos) Abs() Pos {
	return cp.Side.Abs(cp.RelPos)
}

// Give a side co-ordinate from the input grid, and a direction, return where we
// expect to find a side in that direction
// TODO: Merge with Next above (they're the same...)
func OriginPos(side Pos, dir Direction) Pos {
	if dir == D_LEFT {
		return Pos{side.row, side.col - 1}
	} else if dir == D_UP {
		return Pos{side.row - 1, side.col}
	} else if dir == D_RIGHT {
		return Pos{side.row, side.col + 1}
	} else if dir == D_DOWN {
		return Pos{side.row + 1, side.col}
	}
	return Pos{}
}

func SidePoints(input map[Pos]int, size int, row, col int) map[Pos]int {
	rv := map[Pos]int{}
	for r := row; r < row+size; r++ {
		for c := col; c < col+size; c++ {
			p := Pos{r, c}
			if v, found := input[p]; found {
				rv[p] = v
				continue
			}
			log.Fatalf("Misisng %s for side (%d, %d) from cube of size %d", p, row, col, size)
		}
	}
	if len(rv) != size*size {
		log.Fatalf("Side (%d, %d) did not have %d points (%d)", row, col, size*size, len(rv))
	}
	return rv
}

// Parses the 2D input grid into a 3D cube with sides
func NewCube(input map[Pos]int) (Cube, CubePos) {
	rows := 0
	cols := 0
	for p := range input {
		rows = Max(rows, p.row)
		cols = Max(cols, p.col)
	}
	size := GCD(rows, cols)
	fmt.Printf("Found input for a %dx%d cube from %d columns and %d rows\n", size, size, cols/size, rows/size)

	start := CubePos{RelPos: Pos{0, 0}, Heading: D_RIGHT}
	sides := map[Pos]*Side{}
	for r := 0; r < rows/size; r++ {
		for c := 0; c < cols/size; c++ {
			gr := (r * size) + 1
			gc := (c * size) + 1
			if input[Pos{gr, gc}] != VOID {
				origin := Pos{r, c}
				side := NewSide(origin, SidePoints(input, size, gr, gc), size)
				sides[origin] = &side
				if start.Side == nil {
					start.Side = &side
				}
			}
		}
	}
	if len(sides) != 6 {
		log.Fatalf("Found %d sides, expected 6!", len(sides))
	}

	// Find 3 consecutive sides, and make the middle one our "top"
	top := Pos{-1, -1}
COLS:
	for c := 0; c < cols/size; c++ {
		cnt := 0
		for r := 0; r < rows/size; r++ {
			p := Pos{r, c}
			_, exists := sides[p]
			if !exists {
				continue COLS
			}
			cnt++
			fmt.Printf("%s = %d\n", p, cnt)
			if cnt == 3 {
				p.row--
				fmt.Printf("Making %s the top\n", p)
				top = p
				break COLS
			}
		}
	}
	if top.row == -1 { // didn't find in cols, check via rows
	ROWS:
		for r := 0; r < rows/size; r++ {
			cnt := 0
			for c := 0; c < cols/size; c++ {
				p := Pos{r, c}
				_, exists := sides[p]
				if !exists {
					continue ROWS
				}
				cnt++
				fmt.Printf("%s = %d\n", p, cnt)
				if cnt == 3 {
					p.col--
					fmt.Printf("Making %s the top\n", p)
					top = p
					break ROWS
				}
			}
		}

	}
	if top.row == -1 {
		log.Fatal("Didn't find 3 consecutive sides in input!")
	}

	fmt.Println()
	fmt.Println("First side is ", top)
	fmt.Println()

	c := Cube{Size: size}
	c.sides = map[CubeSide]*Side{}
	c.Fold(sides, top, CS_TOP, R_NONE)
	c.Print()
	return c, start
}

var I_RE = regexp.MustCompile(`([LR])?(\d+)`)

func main() {
	s := bufio.NewScanner(os.Stdin)

	input := map[Pos]int{}

	var firstPos Pos
	row := 1
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		spacePreSize := 0
		for col, c := range s.Text() {
			if c == ' ' {
				continue
			} else {
				if spacePreSize == 0 {
					spacePreSize = col
				}
			}
			if c == '.' {
				p := Pos{row, col + 1}
				input[p] = OPEN
				if firstPos.row == 0 {
					firstPos = p
				}
			} else if c == '#' {
				input[Pos{row, col + 1}] = WALL
			} else {
				log.Fatalf("Bad input at %d, unknwon char '%c': %s", col, c, s.Text())
			}
		}
		row++
	}
	if !s.Scan() {
		log.Fatal("Couldn't read instruction line")
	}
	instructions := s.Text()
	//c.Print()
	fmt.Println()
	fmt.Println(instructions)
	fmt.Println()

	c, start := NewCube(input)
	fmt.Println(c)

	if start.Abs() != firstPos {
		log.Fatalf("Expected %s to be %s!", start, firstPos)
	}

	pos := start
	for step, i := range I_RE.FindAllStringSubmatch(instructions, -1) {
		fmt.Printf("* Step %d, at %s\n", step, pos)
		if i[1] != "" {
			pos.Heading = TURN_MAP[pos.Heading][i[1]]
			fmt.Printf("  - Turn %s, Heading: %s\n", i[1], pos.Heading)
		}
		fmt.Printf("  - Moving %s steps\n", i[2])
		pos = c.Nav(pos, Int(i[2]))
		fmt.Println()
	}
	abs := pos.Abs()
	fmt.Printf("Finished at: %s, aka %s\n", pos, abs)
	fmt.Printf("Password: %d * 1000 + %d * 4 + %d = %d\n", abs.row, abs.col, pos.Heading,
		abs.row*1000+abs.col*4+int(pos.Heading))
}

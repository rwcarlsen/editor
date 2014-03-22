package main

import (
	"fmt"
	"strings"
)

type Tabber struct {
	Line []rune
	Tabwidth int
	VisLen int
	// ChToX returns the effective position of a rune indexed by the key if all tabs
	// in the line were expanded to spaces of the Tabber's Tabwidth.
	ChToX []int
	// XToCh does the reverse of ChToX
	XToCh []int
}

func NewTabber(line []rune, tabw int) *Tabber {
	vislen := len(line) + strings.Count(string(line), "\t")*(tabw-1)
	t := &Tabber{
		Line: line,
		Tabwidth: tabw,
		VisLen: vislen,
		ChToX: make([]int, len(line)),
		XToCh: make([]int, vislen),
	}

	n := 0
	for i, r := range line {
		t.ChToX[i] = n
		t.XToCh[n] = i
		if r == '\t' {
			for j := 0; j < t.Tabwidth; j++ {
				t.XToCh[n + j] = i
			}
			n += t.Tabwidth
		} else {
			n++
		}
	}

	return t
}

type View interface {
	Render() Surface
	SetSize(w, h int)
	SetBuf(b *Buffer)
	SetRef(line, char int, x, y int)
	SetTabwidth(n int)
}

type Surface interface {
	Char(x, y int) int
	Line(x, y int) int
	Rune(x, y int) rune
	X(line, char int) int
	Y(line, char int) int
}

func Contains(s Surface, line, char int) bool {
	x, y := RenderPos(s, line, char)
	return x != -1 && y != -1
}

func RenderPos(s Surface, line, char int) (x, y int) {
	return s.X(line, char), s.Y(line, char)
}

func DataPos(s Surface, x, y int) (line, char int) {
	return s.Line(x, y), s.Char(x, y)
}

type WrapView struct {
	w, h           int
	b              *Buffer
	startl, startc int
	startx, starty int
	tabw int
}

func (v *WrapView) Render() Surface {
	surf := &WrapSurf{}
	surf.init(v.w, v.h, v.b, v.startl, v.starty, v.tabw)
	return surf
}

func (v *WrapView) SetSize(w, h int) { v.w, v.h = w, h }
func (v *WrapView) SetTabwidth(n int) { v.tabw = n }
func (v *WrapView) SetBuf(b *Buffer) { v.b = b }
func (v *WrapView) SetRef(line, char int, x, y int) {
	v.startx, v.starty = x, y
	v.startl, v.startc = line, char
}

type PosMap map[int]map[int]int

type WrapSurf struct {
	lines PosMap // map[y]map[x]line#
	chars PosMap // map[y]map[x]char#
	xs    PosMap // map[line#]map[char#]x
	ys    PosMap // map[line#]map[char#]y
	b     *Buffer
}

func (c *WrapSurf) Rune(x, y int) rune {
	l, ch := DataPos(c, x, y)
	if l == -1 || ch == -1 {
		return ' '
	}
	return c.b.Rune(l, ch)
}

func (c *WrapSurf) Char(x, y int) int {
	return c.chars[y][x]
}

func (c *WrapSurf) Line(x, y int) int {
	return c.lines[y][x]
}

func (c *WrapSurf) X(line, char int) int {
	if v, ok := c.xs[line][char]; ok {
		return v
	}
	return -1
}

func (c *WrapSurf) Y(line, char int) int {
	if v, ok := c.ys[line][char]; ok {
		return v
	}
	return -1
}

func (c *WrapSurf) init(w, h int, b *Buffer, startl, starty int, tabw int) {
	c.b = b
	c.lines = PosMap{}
	c.chars = PosMap{}
	c.xs = PosMap{}
	c.ys = PosMap{}

	// figure out line+char for top left corner of canvas
	l, ch := startl, 0
	y := starty
	for l > 0 && y > 0 {
		l--
		line := b.Line(l)
		t := NewTabber(line, tabw)
		dy := t.VisLen/w + 1
		y -= dy
		ch = 0
		if y < 0 && t.VisLen > w {
			ch = t.XToCh[w*-1*y]
		}
	}

	lg.Printf("startl=%v, starty=%v, l=%v, ch=%v", startl, starty, l, ch)
	// draw from start line and char down
	for y := 0; y < h; y++ {
		if c.xs[l] == nil {
			c.xs[l] = map[int]int{}
			c.ys[l] = map[int]int{}
		}
		if c.chars[y] == nil {
			c.chars[y] = map[int]int{}
			c.lines[y] = map[int]int{}
		}

		var line []rune
		if l < b.Nlines() {
			line = b.Line(l)
		}
		startch := ch
		t := NewTabber(line[startch:], tabw)
		for x := 0; x < w; x++ {
			if x < t.VisLen {
				ch = startch + t.XToCh[x]
			} else {
				ch = len(line)
			}

			if _, ok := c.xs[l][ch]; !ok {
				c.xs[l][ch] = x
				c.ys[l][ch] = y
			}
			c.chars[y][x] = ch
			c.lines[y][x] = l

			if l >= b.Nlines() {
				c.lines[y][x] = -1
			}
			if x >= t.VisLen {
				c.chars[y][x] = -1
			}
		}
		lg.Printf("%+v", c.ys[l])

		if ch == len(line) { // if we drew entire line
			ch = 0
			l++    // go to next line
		}
	}
}

type LineNumView struct {
	View
	b       *Buffer
	w, h    int
	ndigits int
}

func (v *LineNumView) Render() Surface {
	linenums := map[int]map[int]rune{}
	v.ndigits = len(fmt.Sprint(v.b.Nlines())) + 1
	surf := v.View.Render()

	prev := -1
	for y := 0; y < v.h; y++ {
		line := surf.Line(0, y)
		nums := ""
		if line == -1 {
			break
		} else if line != prev {
			nums = fmt.Sprint(line + 1)
		}
		prev = line
		nums = strings.Repeat(" ", v.ndigits-1-len(nums)) + nums + " "
		for n := 0; n < v.ndigits; n++ {
			if _, ok := linenums[n]; !ok {
				linenums[n] = map[int]rune{}
			}
			linenums[n][y] = rune(nums[n])
		}
	}

	return &LineNumSurf{
		Surface: surf,
		ndigits: v.ndigits,
		nums:    linenums,
	}
}

func (v *LineNumView) SetSize(w, h int) {
	v.w, v.h = w, h
	v.View.SetSize(w-v.ndigits, h)
}
func (v *LineNumView) SetBuf(b *Buffer) {
	v.b = b
	v.View.SetBuf(b)
}
func (v *LineNumView) SetRef(line, char int, x, y int) {
	v.View.SetRef(line, char, x-v.ndigits, y)
}

type LineNumSurf struct {
	Surface
	ndigits int
	nums    map[int]map[int]rune
}

func (s *LineNumSurf) Char(x, y int) int {
	return s.Surface.Char(x-s.ndigits, y)
}
func (s *LineNumSurf) Line(x, y int) int {
	return s.Surface.Line(x-s.ndigits, y)
}
func (s *LineNumSurf) Rune(x, y int) rune {
	if x < s.ndigits {
		return s.nums[x][y]
	} else {
		return s.Surface.Rune(x-s.ndigits, y)
	}
}
func (s *LineNumSurf) X(line, char int) int {
	return s.Surface.X(line, char) + s.ndigits
}

package view

import (
	"strings"

	"github.com/rwcarlsen/editor/util"
)

type View interface {
	Render() Surface
	SetSize(w, h int)
	SetBuf(b *util.Buffer)
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

func Draw(s Surface, x, y int) {
	for y := 0; y < s.H; y++ {
		for x := 0; x < s.W; x++ {
			termbox.SetCell(x, y, surf.Rune(x, y), 0, 0)
		}
	}
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

type Wrap struct {
	w, h           int
	b              *util.Buffer
	startl, startc int
	startx, starty int
	tabw           int
}

func (v *Wrap) Render() Surface {
	surf := &WrapSurf{}
	surf.init(v.w, v.h, v.b, v.startl, v.starty, v.tabw)
	return surf
}

func (v *Wrap) SetSize(w, h int)      { v.w, v.h = w, h }
func (v *Wrap) SetTabwidth(n int)     { v.tabw = n }
func (v *Wrap) SetBuf(b *util.Buffer) { v.b = b }
func (v *Wrap) SetRef(line, char int, x, y int) {
	v.startx, v.starty = x, y
	v.startl, v.startc = line, char
}

type PosMap map[int]map[int]int

type WrapSurf struct {
	lines PosMap // map[y]map[x]line#
	chars PosMap // map[y]map[x]char#
	xs    PosMap // map[line#]map[char#]x
	ys    PosMap // map[line#]map[char#]y
	b     *util.Buffer
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

func (c *WrapSurf) init(w, h int, b *util.Buffer, startl, starty int, tabw int) {
	c.b = b
	c.lines = PosMap{}
	c.chars = PosMap{}
	c.xs = PosMap{}
	c.ys = PosMap{}

	// figure out line+char for top left corner of canvas
	l, nextch := FindStart(b, w, startl, starty, tabw)

	// draw from start line and char down
	for y := 0; y < h; y++ {
		if c.xs[l] == nil {
			c.xs[l], c.ys[l] = map[int]int{}, map[int]int{}
		}
		if c.chars[y] == nil {
			c.chars[y], c.lines[y] = map[int]int{}, map[int]int{}
		}

		var line []rune
		if l < b.Nlines() {
			line = b.Line(l)
		}

		var chs []int
		chs, nextch = RenderLine(line, nextch, w, tabw)
		for x := 0; x < w; x++ {
			ch := chs[x]
			c.chars[y][x] = ch
			c.lines[y][x] = l
			if l >= b.Nlines() {
				c.lines[y][x] = -1
			}

			if _, ok := c.xs[l][ch]; !ok {
				c.xs[l][ch] = x
				c.ys[l][ch] = y
			}
		}

		if nextch >= len(line) { // if we drew entire line
			nextch = 0
			l++ // go to next line
		}
	}
}

// RenderLine renders the given line starting at the rune indexed by startch
// and returns a slice of length w of indices that index into line for each x
// tile. Indices of -1 indicate that no character is drawn.  nextch indexes the
// first rune in line that didn't fit across the screen.
func RenderLine(line []rune, startch int, w, tabw int) (chs []int, nextch int) {
	chs = make([]int, w)
	nextch = startch

	t := NewTabber(line[startch:], tabw)
	for x := 0; x < w; x++ {
		if x < len(t.XToCh) {
			chs[x] = startch + t.XToCh[x]
			nextch = chs[x] + 1
		} else {
			chs[x] = -1
		}
	}
	return chs, nextch
}

func FindStart(b *util.Buffer, w int, startl, starty int, tabw int) (line, char int) {
	line, char = startl, 0
	y := starty
	for line > 0 && y > 0 {
		line--
		l := b.Line(line)
		t := NewTabber(l, tabw)
		dy := t.VisLen/w + 1
		y -= dy
		char = 0
		if y < 0 && t.VisLen > w {
			char = t.XToCh[w*-1*y]
		}
	}
	return line, char
}

type Tabber struct {
	Line     []rune
	Tabwidth int
	VisLen   int
	// ChToX returns the effective position of a rune indexed by the key if all tabs
	// in the line were expanded to spaces of the Tabber's Tabwidth.
	ChToX []int
	// XToCh does the reverse of ChToX
	XToCh []int
}

func NewTabber(line []rune, tabw int) *Tabber {
	vislen := len(line) + strings.Count(string(line), "\t")*(tabw-1)
	t := &Tabber{
		Line:     line,
		Tabwidth: tabw,
		VisLen:   vislen,
		ChToX:    make([]int, len(line)),
		XToCh:    make([]int, vislen),
	}

	n := 0
	for i, r := range line {
		t.ChToX[i] = n
		t.XToCh[n] = i
		if r == '\t' {
			for j := 0; j < t.Tabwidth; j++ {
				t.XToCh[n+j] = i
			}
			n += t.Tabwidth
			t.XToCh[n-1] = i
			t.ChToX[i] = n - 1
		} else {
			n++
		}
	}

	return t
}


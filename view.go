package main

type View interface {
	Render() Surface
	SetSize(w, h int)
	SetBuf(b *Buffer)
	SetRef(line, char int, x, y int)
}

type Surface interface {
	Char(x, y int) int
	Line(x, y int) int
	X(line, char int) int
	Y(line, char int) int
}

type WrapView struct {
	w, h           int
	b              *Buffer
	startl, startc int
	startx, starty int
}

func (v *WrapView) Render() Surface {
	surf := &WrapSurf{}
	surf.init(v.w, v.h, v.b, v.startl, v.starty)
	return surf
}

func (v *WrapView) SetSize(w, h int) { v.w, v.h = w, h }
func (v *WrapView) SetBuf(b *Buffer) { v.b = b }
func (v *WrapView) SetRef(line, char int, x, y int) {
	v.startx, v.starty = x, y
	v.startl, v.startc = line, char
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

type PosMap map[int]map[int]int

type WrapSurf struct {
	lines PosMap // map[y]map[x]line#
	chars PosMap // map[y]map[x]char#
	xs    PosMap // map[line#]map[char#]x
	ys    PosMap // map[line#]map[char#]y
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

func (c *WrapSurf) init(w, h int, b *Buffer, startl, starty int) {
	c.lines = PosMap{}
	c.chars = PosMap{}
	c.xs = PosMap{}
	c.ys = PosMap{}

	// figure out line+char for top left corner of canvas
	l, ch := startl, 0
	if starty > 0 {
		y := starty - 1
		for l > 0 {
			l--
			line := b.Line(l)
			dy := len(line)/w + 1
			if dy > y && dy == 1 {
				ch = 0
				break
			} else if dy > y && dy > 1 {
				ch = len(line) % w
				break
			}
			y -= dy
		}
	}

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
		for x := 0; x < w; x++ {
			c.xs[l][ch] = x
			c.ys[l][ch] = y
			c.chars[y][x] = ch
			c.lines[y][x] = l
			if ch >= len(line) {
				c.chars[y][x] = -1
			}
			if l >= b.Nlines() {
				c.lines[y][x] = -1
			}
			ch++
		}

		if ch >= len(line) { // if we drew entire line
			l++    // go to next line
			ch = 0 // at first char
		}
	}
}

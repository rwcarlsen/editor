package main

type View interface {
	Line(x, y int) int
	Char(x, y int) int
	X(line, char int) int
	Y(line, char int) int
}

type PosMap map[int]map[int]int

type WrapView struct {
	lines PosMap // map[y]map[x]line#
	chars PosMap // map[y]map[x]char#
	xs    PosMap // map[line#]map[char#]x
	ys    PosMap // map[line#]map[char#]y
}

func RenderPos(v View, line, char int) (x, y int) {
	return v.X(line, char), v.Y(line, char)
}

func DataPos(v View, x, y int) (line, char int) {
	return v.Line(x, y), v.Char(x, y)
}

func NewWrapView(w, h int, content [][]rune, startl, starty int) *WrapView {
	c := &WrapView{}
	c.init(w, h, content, startl, starty)
	return c
}

func (c *WrapView) Char(x, y int) int {
	return c.chars[y][x]
}

func (c *WrapView) Line(x, y int) int {
	return c.lines[y][x]
}

func (c *WrapView) X(line, char int) int {
	if v, ok := c.xs[line][char]; ok {
		return v
	}
	return -1
}

func (c *WrapView) Y(line, char int) int {
	if v, ok := c.ys[line][char]; ok {
		return v
	}
	return -1
}

func (c *WrapView) init(w, h int, content [][]rune, startl, starty int) {
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
			line := content[l]
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
		lg.Printf("l=%v, y=%v\n", l, y)
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
		if l < len(content) {
			line = content[l]
		}
		for x := 0; x < w; x++ {
			c.xs[l][ch] = x
			c.ys[l][ch] = y
			c.chars[y][x] = ch
			c.lines[y][x] = l
			if ch >= len(line) {
				c.chars[y][x] = -1
			}
			if l >= len(content) {
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

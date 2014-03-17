package main

type View interface {
	Line(x, y int) int
	Char(x, y int) int
	X(line, char int) int
	Y(line, char int) int
}

type PosMap map[int]map[int]int

type WrapCanvas struct {
	lines PosMap // map[y]map[x]line#
	chars PosMap // map[y]map[x]char#
	xs    PosMap // map[line#]map[char#]x
	ys    PosMap // map[line#]map[char#]y
}

func NewWrapCanvas(w, h int, content [][]rune, startl, starty int) *WrapCanvas {
	c := &WrapCanvas{}
	c.init(w, h, content, startl, starty)
	return c
}

func (c *WrapCanvas) Char(x, y int) int {
	return c.chars[y][x]
}

func (c *WrapCanvas) Line(x, y int) int {
	return c.lines[y][x]
}

func (c *WrapCanvas) X(line, char int) int {
	return c.xs[line][char]
}

func (c *WrapCanvas) Y(line, char int) int {
}

func (c *WrapCanvas) init(w, h int, content [][]rune, startl, starty int) {
	c.lines = PosMap{}
	c.chars = PosMap{}
	c.xs = PosMap{}
	c.ys = PosMap{}

	// figure out line+char for top left corner of canvas
	l, ch := startl-1, 0
	y := starty
	for l >= 0 {
		line := content[l]
		dy := len(line) / w + 1
		if y - dy < 0 && dy == 1 {
			ch = 0
			break
		} else if y - dy < 0 && dy > 1 {
			ch = len(line) % w
			break
		}
		l--
		y -= dy
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

package main

type PosMap map[int]int

type Canvas struct {
	Line PosMap // map[y]line#
	Char map[int]PosMap // map[y]map[x]char#
	X map[int]PosMap // map[line#]map[char#]x
	Y map[int]PosMap // map[line#]map[char#]y
}

func NewCanvas(w, h int, content [][]rune, startline int) *Canvas {
	c := &Canvas{}
	c.init(w, h, content, startline)
	return c
}

func (c *Canvas) DataPos(x, y int) (line, char int) {
	return c.Line[y], c.Char[y][x]
}

func (c *Canvas) RenderPos(line, char int) (x, y int) {
	return c.X[line][char], c.Y[line][char]
}

func (c *Canvas) init(w, h int, content [][]rune, startline int) {
	c.Line = PosMap{}
	c.Char = map[int]PosMap{}
	c.X = map[int]PosMap{}
	c.Y = map[int]PosMap{}

	l, ch := startline, 0
	for y := 0; y < h; y++ {
		if c.X[l] == nil {
			c.X[l] = PosMap{}
			c.Y[l] = PosMap{}
		}
		if c.Char[y] == nil {
			c.Char[y] = PosMap{}
		}

		c.Line[y] = -1
		var line []rune
		if l < len(content) {
			line = content[l]
			c.Line[y] = l
		}
		for x := 0; x < w; x++ {
			c.X[l][ch] = x
			c.Y[l][ch] = y
			c.Char[y][x] = ch
			if ch >= len(line) {
				c.Char[y][x] = -1
			}
			ch++
		}

		if ch >= len(line) { // if we drew entire line
			l++   // go to next line
			ch = 0 // at first char
		}
	}
}

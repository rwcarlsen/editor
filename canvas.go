package main

type PosMap map[int]int

type Canvas struct {
	Line PosMap // map[y]line#
	Char map[int]PosMap // map[y]map[x]char#
	X map[int]PosMap // map[line#]map[char#]x
	Y map[int]PosMap // map[line#]map[char#]y
	cont map[int]bool // map[line#]bool
}

func NewCanvas(w, h int, content [][]rune, startline int) *Canvas {
	c := &Canvas{}
	c.init(w, h, content, startline)
	return c
}

func (c *Canvas) DataPos(x, y int) (line, char int) {
	return c.Line[y], c.Char[y][x]
}

func (c *Canvas) Contains(line int) bool {
	return c.cont[line]
}

func (c *Canvas) RenderPos(line, char int) (x, y int) {
	return c.X[line][char], c.Y[line][char]
}

func (c *Canvas) init(w, h int, content [][]rune, startline int) {
	c.Line = PosMap{}
	c.Char = map[int]PosMap{}
	c.X = map[int]PosMap{}
	c.Y = map[int]PosMap{}
	c.cont = map[int]bool{}

	l, ch := startline, 0
	for y := 0; y < h; y++ {
		c.Line[y] = -1
		var line []rune
		if l < len(content) {
			line = content[l]
			c.Line[y] = l
			c.cont[l] = true
		}
		for x := 0; x < w; x++ {
			if ch >= len(line) {
				if c.Char[y] == nil {
					c.Char[y] = PosMap{}
				}
				c.Char[y][x] = -1
				continue
			}

			if c.X[l] == nil {
				c.X[l] = PosMap{}
				c.Y[l] = PosMap{}
			}
			c.X[l][ch] = x
			c.Y[l][ch] = y

			if c.Char[y] == nil {
				c.Char[y] = PosMap{}
			}
			c.Char[y][x] = ch
			ch++
		}

		if ch >= len(line) { // if we drew entire line
			l++   // go to next line
			ch = 0 // at first char
		}
	}
}

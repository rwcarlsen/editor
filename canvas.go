package main

type PosMap map[int]map[int]int

type Canvas struct {
	Line PosMap // map[x]map[y]line#
	Char PosMap // map[x]map[y]char#
	X    PosMap // map[line#]map[char#]x
	Y    PosMap // map[line#]map[char#]y
}

func NewCanvas(w, h int, content [][]rune, startline int) *Canvas {
	c := &Canvas{}
	c.init(w, h, content, startline)
	return c
}

func (c *Canvas) DataPos(x, y int) (line, char int) {
	return c.Line[x][y], c.Char[x][y]
}

func (c *Canvas) Contains(line, char int) bool {
	if c.X[line] == nil {
		return false
	} else if _, ok := c.X[line][char]; !ok {
		return false
	} else if c.Y[line] == nil {
		return false
	} else if _, ok := c.Y[line][char]; !ok {
		return false
	}
	return true
}

func (c *Canvas) RenderPos(line, char int) (x, y int) {
	return c.X[line][char], c.Y[line][char]
}

func (c *Canvas) init(w, h int, content [][]rune, startline int) {
	c.Line = PosMap{}
	c.Char = PosMap{}
	c.X = PosMap{}
	c.Y = PosMap{}

	l, c := startline, 0
	wrapCount := 0
	for y := 0; y < h; y++ {
		var line []rune
		if l < len(s.Content) {
			line = s.Content[l]
		}
		for x := 0; x < w; x++ {
			if c >= len(line) {
				if c.Line[x] == nil {
					c.Line[x] = map[int]int{}
					c.Char[x] = map[int]int{}
				}
				c.Line[x][y] = -1
				c.Char[x][y] = -1
				continue
			}

			if c.X[l] == nil {
				c.X[l] = map[int]int{}
				c.Y[l] = map[int]int{}
			}
			c.X[l][c] = x
			c.Y[l][c] = y

			if c.Line[x] == nil {
				c.Line[x] = map[int]int{}
				c.Char[x] = map[int]int{}
			}
			c.Line[x][y] = l
			c.Char[x][y] = c
			c++
		}

		if c >= len(line) { // if we drew entire line
			l++   // go to next line
			c = 0 // at first char
			wrapCount = 0
		} else {
			wrapCount++
		}
	}
}

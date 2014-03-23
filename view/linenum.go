package view

import (
	"fmt"
	"strings"

	"github.com/rwcarlsen/editor/util"
)

type LineNum struct {
	View
	b       *util.Buffer
	w, h    int
	ndigits int
}

func (v *LineNum) Render() Surface {
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

func (v *LineNum) SetSize(w, h int) {
	v.w, v.h = w, h
	v.View.SetSize(w-v.ndigits, h)
}
func (v *LineNum) SetBuf(b *util.Buffer) {
	v.b = b
	v.View.SetBuf(b)
}
func (v *LineNum) SetRef(line, char int, x, y int) {
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

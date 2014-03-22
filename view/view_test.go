package view

import (
	"strings"
	"testing"

	"github.com/rwcarlsen/editor/util"
)

type tabbertest struct {
	tabw     int
	line     []rune
	expectx  []int
	expectch []int
}

var tabbertests = []tabbertest{
	tabbertest{
		tabw:     1,
		line:     []rune("\t"),
		expectx:  []int{0},
		expectch: []int{0},
	},
	tabbertest{
		tabw:     1,
		line:     []rune("\t "),
		expectx:  []int{0, 1},
		expectch: []int{0, 1},
	},
	tabbertest{
		tabw:     1,
		line:     []rune(" \t"),
		expectx:  []int{0, 1},
		expectch: []int{0, 1},
	},
	tabbertest{
		tabw:     2,
		line:     []rune("\t "),
		expectx:  []int{0, 2},
		expectch: []int{0, 0, 1},
	},
	tabbertest{
		tabw:     2,
		line:     []rune(" \t"),
		expectx:  []int{0, 1},
		expectch: []int{0, 1, 1},
	},
	tabbertest{
		tabw:     3,
		line:     []rune(" \t\t \t"),
		expectx:  []int{0, 1, 4, 7, 8},
		expectch: []int{0, 1, 1, 1, 2, 2, 2, 3, 4, 4, 4},
	},
}

func TestTabber(t *testing.T) {
	fmtstr := "test %v %v (line='%v', tabw=%v): expected %+v, got %+v"
	for i, test := range tabbertests {
		tabber := NewTabber(test.line, test.tabw)
		str := strings.Replace(string(test.line), "\t", "\\t", -1)

		if len(test.expectx) != len(tabber.ChToX) {
			t.Errorf(fmtstr, i, "ChToX", str, test.tabw, test.expectx, tabber.ChToX)
		} else {
			for j := range test.expectx {
				if test.expectx[j] != tabber.ChToX[j] {
					t.Errorf(fmtstr, i, "ChToX", str, test.tabw, test.expectx,
						tabber.ChToX)
					break
				}
			}
		}

		if len(test.expectch) != len(tabber.XToCh) {
			t.Errorf(fmtstr, i, "XToCh", str, test.tabw, test.expectch, tabber.XToCh)
		} else {
			for j := range test.expectch {
				if test.expectch[j] != tabber.XToCh[j] {
					t.Errorf(fmtstr, i, "XToCh", str, test.tabw, test.expectch,
						tabber.XToCh)
					break
				}
			}
		}
	}
}

type viewtest struct {
	text       string
	tabw       int
	w, h       int
	l, c, x, y int
	expectch   [][]int
	expectl    [][]int
	expectx    [][]int
	expecty    [][]int
}

var viewtests = []viewtest{
	viewtest{
		text: "abc\ndef",
		tabw: 1, w: 3, h: 2,
		l: 0, c: 0, x: 0, y: 0,
		expectch: [][]int{
			[]int{0, 1, 2},
			[]int{0, 1, 2},
		},
		expectl: [][]int{
			[]int{0, 0, 0},
			[]int{1, 1, 1},
		},
		expectx: [][]int{
			[]int{0, 1, 2},
			[]int{0, 1, 2},
		},
		expecty: [][]int{
			[]int{0, 0, 0},
			[]int{1, 1, 1},
		},
	},
	viewtest{
		text: "abc\ndef",
		tabw: 1, w: 4, h: 2,
		l: 0, c: 0, x: 0, y: 0,
		expectch: [][]int{
			[]int{0, 1, 2, -1},
			[]int{0, 1, 2, -1},
		},
		expectl: [][]int{
			[]int{0, 0, 0, 0},
			[]int{1, 1, 1, 1},
		},
		expectx: [][]int{
			[]int{0, 1, 2},
			[]int{0, 1, 2},
		},
		expecty: [][]int{
			[]int{0, 0, 0},
			[]int{1, 1, 1},
		},
	},
	viewtest{
		text: "abcd\nef",
		tabw: 1, w: 3, h: 3,
		l: 0, c: 0, x: 0, y: 0,
		expectch: [][]int{
			[]int{0, 1, 2},
			[]int{3, -1, -1},
			[]int{0, 1, -1},
		},
		expectl: [][]int{
			[]int{0, 0, 0},
			[]int{0, 0, 0},
			[]int{1, 1, 1},
		},
		expectx: [][]int{
			[]int{0, 1, 2, 0},
			[]int{0, 1},
		},
		expecty: [][]int{
			[]int{0, 0, 0, 1},
			[]int{2, 2},
		},
	},
}

func TestWrap(t *testing.T) {
	v := &Wrap{}

	for i, tst := range viewtests {
		printtxt := strings.Replace(tst.text, "\n", "\\n", -1)
		printtxt = strings.Replace(printtxt, "\t", "\\t", -1)
		t.Logf("* test %v: text='%v' tabw=%v (w,h) = (%v,%v)"+
			"(l,c,x,y) = (%v,%v,%v,%v",
			i, printtxt, tst.tabw, tst.w, tst.h, tst.l, tst.c, tst.x, tst.y)

		b := util.NewBuffer([]byte(tst.text))
		v.SetBuf(b)
		v.SetSize(tst.w, tst.h)
		v.SetRef(tst.l, tst.c, tst.x, tst.y)
		v.SetTabwidth(tst.tabw)
		surf := v.Render()

		t.Log("\t* checking Surface.Char(x, y)")
		for y, row := range tst.expectch {
			for x, ch := range row {
				if got := surf.Char(x, y); ch != got {
					t.Errorf("\t\t[**] for x,y = %v,%v: expected %v, got %v",
						x, y, ch, got)
					printSurf(t, surf, tst, b)
				} else {
					t.Logf("\t\t[OK] for x,y = %v,%v: expected %v, got %v",
						x, y, ch, got)
				}
			}
		}

		t.Log("\t* checking Surface.Line(x, y)")
		for y, row := range tst.expectl {
			for x, l := range row {
				if got := surf.Line(x, y); l != got {
					t.Errorf("\t\t[**] for x,y = %v,%v: expected %v, got %v",
						x, y, l, got)
					printSurf(t, surf, tst, b)
				} else {
					t.Logf("\t\t[OK] for x,y = %v,%v: expected %v, got %v",
						x, y, l, got)
				}
			}
		}

		t.Log("\t* checking Surface.X(l, ch)")
		for l, row := range tst.expectx {
			for ch, x := range row {
				if got := surf.X(l, ch); x != got {
					t.Errorf("\t\t[**] for l,c = %v,%v: expected %v, got %v",
						l, ch, x, got)
					printSurf(t, surf, tst, b)
				} else {
					t.Logf("\t\t[OK] for l,c = %v,%v: expected %v, got %v",
						l, ch, x, got)
				}
			}
		}

		t.Log("\t* checking Surface.Y(l, ch)")
		for l, row := range tst.expecty {
			for ch, y := range row {
				if got := surf.Y(l, ch); y != got {
					t.Errorf("\t\t[**] for l,c = %v,%v: expected %v, got %v",
						l, ch, y, got)
					printSurf(t, surf, tst, b)
				} else {
					t.Logf("\t\t[OK] for l,c = %v,%v: expected %v, got %v",
						l, ch, y, got)
				}
			}
		}
	}
}

func printSurf(t *testing.T, surf Surface, tst viewtest, b *util.Buffer) {
	t.Log("")
	for y := 0; y < tst.h; y++ {
		got := ""
		expect := ""
		pretxt := "   "
		midtxt := "      "
		if y == 0 {
			pretxt = "got"
			midtxt = "expect"
		}
		for x := 0; x < tst.w; x++ {
			got += string(surf.Rune(x, y))
			l, ch := tst.expectl[y][x], tst.expectch[y][x]
			if l == -1 || ch == -1 {
				expect += " "
			} else {
				expect += string(b.Rune(l, ch))
			}
		}
		t.Logf("\t\t\t%v   |%v|   %v   |%v|", pretxt, got, midtxt, expect)
	}
	t.Log("")
}

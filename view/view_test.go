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
		expectx:  []int{1, 2},
		expectch: []int{0, 0, 1},
	},
	tabbertest{
		tabw:     2,
		line:     []rune(" \t"),
		expectx:  []int{0, 2},
		expectch: []int{0, 1, 1},
	},
	tabbertest{
		tabw:     3,
		line:     []rune(" \t\t \t"),
		expectx:  []int{0, 3, 6, 7, 10},
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
	name       string
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
		name: "no wrap, full screen",
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
		name: "no wrap, short lines",
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
		name: "simple wrap",
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
	viewtest{
		name: "no wrap, simple tab",
		text: "a\tb",
		tabw: 2, w: 4, h: 1,
		l: 0, c: 0, x: 0, y: 0,
		expectch: [][]int{
			[]int{0, 1, 1, 2},
		},
		expectl: [][]int{
			[]int{0, 0, 0, 0},
		},
		expectx: [][]int{
			[]int{0, 2, 3},
		},
		expecty: [][]int{
			[]int{0, 0, 0},
		},
	},
	viewtest{
		name: "no wrap, empty line",
		text: "a\n\nb",
		tabw: 1, w: 2, h: 3,
		l: 0, c: 0, x: 0, y: 0,
		expectch: [][]int{
			[]int{0, -1},
			[]int{-1, -1},
			[]int{0, -1},
		},
		expectl: [][]int{
			[]int{0, 0},
			[]int{1, 1},
			[]int{2, 2},
		},
		expectx: [][]int{
			[]int{0},
			[]int{-1},
			[]int{0},
		},
		expecty: [][]int{
			[]int{0},
			[]int{1},
			[]int{2},
		},
	},
}

func TestWrap0(t *testing.T) { testWrap(t, 0) }
func TestWrap1(t *testing.T) { testWrap(t, 1) }
func TestWrap2(t *testing.T) { testWrap(t, 2) }
func TestWrap3(t *testing.T) { testWrap(t, 3) }
func TestWrap4(t *testing.T) { testWrap(t, 4) }

func testWrap(t *testing.T, i int) {
	v := &Wrap{}
	tst := viewtests[i]

	printtxt := strings.Replace(tst.text, "\n", "\\n", -1)
	printtxt = strings.Replace(printtxt, "\t", "\\t", -1)
	t.Logf("* test %v (%v): text='%v' tabw=%v (w,h) = (%v,%v)"+
		" (l,c,x,y) = (%v,%v,%v,%v",
		i, tst.name, printtxt, tst.tabw, tst.w, tst.h,
		tst.l, tst.c, tst.x, tst.y)

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
				t.Errorf("\t\t[--] for x,y = %v,%v: expected %v, got %v",
					x, y, ch, got)
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
				t.Errorf("\t\t[--] for x,y = %v,%v: expected %v, got %v",
					x, y, l, got)
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
				t.Errorf("\t\t[--] for l,c = %v,%v: expected %v, got %v",
					l, ch, x, got)
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
				t.Errorf("\t\t[--] for l,c = %v,%v: expected %v, got %v",
					l, ch, y, got)
			} else {
				t.Logf("\t\t[OK] for l,c = %v,%v: expected %v, got %v",
					l, ch, y, got)
			}
		}
	}

	printSurf(t, surf, tst, b)
}

func printSurf(t *testing.T, surf Surface, tst viewtest, b *util.Buffer) {
	t.Log("")
	for y := 0; y < tst.h; y++ {
		got := ""
		expect := ""
		pretxt := "        "
		midtxt := "   "
		if y == 0 {
			pretxt = "expected"
			midtxt = "got"
		}
		for x := 0; x < tst.w; x++ {
			if r := surf.Rune(x, y); r != '\t' {
				got += string(r) + " "
			} else {
				got += "\\t"
			}

			l, ch := tst.expectl[y][x], tst.expectch[y][x]
			if l == -1 || ch == -1 {
				expect += "  "
			} else if r := b.Rune(l, ch); r != '\t' {
				expect += string(r) + " "
			} else {
				expect += "\\t"
			}
		}
		expect = strings.Replace(expect, "\t", " ", -1)
		got = strings.Replace(got, "\t", " ", -1)
		t.Logf("\t%v   |%v|   %v   |%v|", pretxt, expect, midtxt, got)
	}
	t.Log("")
}

type renderlinetest struct {
	line             string
	startch, w, tabw int
	expectchs        []int
	expectnextch     int
}

var renderlinetests = []renderlinetest{
	renderlinetest{
		line:    "abc",
		startch: 0, w: 3, tabw: 1,
		expectchs:    []int{0, 1, 2},
		expectnextch: 3,
	},
	renderlinetest{
		line:    "",
		startch: 0, w: 3, tabw: 1,
		expectchs:    []int{-1, -1, -1},
		expectnextch: 0,
	},
	renderlinetest{
		line:    "abc",
		startch: 0, w: 5, tabw: 1,
		expectchs:    []int{0, 1, 2, -1, -1},
		expectnextch: 3,
	},
	renderlinetest{
		line:    "abc\t",
		startch: 0, w: 6, tabw: 3,
		expectchs:    []int{0, 1, 2, 3, 3, 3},
		expectnextch: 4,
	},
	renderlinetest{
		line:    "abc\t",
		startch: 0, w: 5, tabw: 3,
		expectchs:    []int{0, 1, 2, 3, 3},
		expectnextch: 4,
	},
	renderlinetest{
		line:    "ab\tc",
		startch: 0, w: 5, tabw: 3,
		expectchs:    []int{0, 1, 2, 2, 2},
		expectnextch: 3,
	},
	renderlinetest{
		line:    "a\tb\tc",
		startch: 0, w: 7, tabw: 3,
		expectchs:    []int{0, 1, 1, 1, 2, 3, 3},
		expectnextch: 4,
	},
}

func TestRenderLine(t *testing.T) {
	for i, tst := range renderlinetests {
		str := strings.Replace(tst.line, "\t", "\\t", -1)
		t.Logf("* test %v: line='%v', startch=%v, w=%v, tabw=%v",
			i, str, tst.startch, tst.w, tst.tabw)

		chs, nextch := RenderLine([]rune(tst.line), tst.startch, tst.w, tst.tabw)
		if nextch != tst.expectnextch {
			t.Errorf("\texpected nextch = %+v, got %+v", tst.expectnextch, nextch)
		}
		if len(chs) != len(tst.expectchs) {
			t.Fatalf("\texpected chs = %+v, got %+v", tst.expectchs, chs)
		} else {
			for j := range chs {
				if chs[j] != tst.expectchs[j] {
					t.Fatalf("\texpected chs = %+v, got %+v", tst.expectchs, chs)
				}
			}
		}
		if !t.Failed() {
			t.Logf("\tPASSED")
		}
	}
}


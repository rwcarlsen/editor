package main

import (
	"strings"
	"testing"
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

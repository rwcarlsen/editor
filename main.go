package main

import (
	"flag"
	"fmt"
	termbox "github.com/nsf/termbox-go"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var ErrQuit = fmt.Errorf("Quit")
var flog *os.File
var lg *log.Logger

func main() {
	flag.Parse()
	log.SetFlags(0)

	flog, err := os.Create("editor.log")
	if err != nil {
		log.Fatal(err)
	}
	defer flog.Close()
	lg = log.New(flog, "", 0)

	// start terminal
	err = termbox.Init()
	if err != nil {
		log.Print(err)
		return
	}
	defer termbox.Close()

	// start session
	fname := flag.Arg(0)
	s, err := NewSession(fname)
	if err != nil {
		log.Print(err)
		return
	}

	// run ...
	err = s.Run()
	if err != ErrQuit {
		log.Print(err)
	}
}

type Screen struct {
	lineshift int     // which line is the first we draw
	lineMap   [][]int // [screenY][screenx]line#
	charMap   [][]int // [screenY][screenX]char#
	LineNums  bool
}

func (s *Screen) Scroll(n int, lines []string) {
	_, h := termbox.Size()
	s.lineshift += n
	if s.lineshift < 0 {
		s.lineshift = 0
	} else if s.lineshift > len(lines)-h {
		s.lineshift = len(lines) - h
	}
}

func (s *Screen) Loc(x, y int) (line, char int) {
	return s.lineMap[y][x], s.charMap[y][x]
}

func (s *Screen) Draw(xpos, ypos int, lines [][]rune) {
	s.clear()
	defer termbox.Flush()

	ndigits := len(fmt.Sprint(len(lines))) + 1
	if s.LineNums {
		xpos += ndigits
	}

	w, h := termbox.Size()
	x, y := 0, s.lineshift // char#, line#
	for i := ypos; i < h; i++ {
		if y >= len(lines) {
			break
		}
		line := lines[y]
		for j := xpos; j < w; j++ {
			if x >= len(line) {
				termbox.SetCell(j, i, ' ', 0, 0)
			} else {
				termbox.SetCell(j, i, line[x], 0, 0)
				s.lineMap[i][j] = y
				s.charMap[i][j] = x
				x++
			}

		}

		if s.LineNums {
			nums := fmt.Sprint(y + 1)
			nums = strings.Repeat(" ", ndigits - 1 - len(nums)) + nums + " "
			for n := 0; n < ndigits; n++ {
				termbox.SetCell(xpos - ndigits + n, i, rune(nums[n]), 0, 0)
			}
		}

		if x >= len(line) { // if we drew entire line
			y++   // go to next line
			x = 0 // at first char
		}
	}
}

func (s *Screen) clear() {
	w, y := termbox.Size()

	s.lineMap = make([][]int, y)
	s.charMap = make([][]int, y)
	for i := range s.lineMap {
		s.lineMap[i] = make([]int, w)
		s.charMap[i] = make([]int, w)
		for j := range s.lineMap[i] {
			s.lineMap[i][j] = -1
			s.charMap[i][j] = -1
		}
	}
}

type Session struct {
	CursorX int
	CursorY int
	File    string
	Lines   [][]rune
	scr     *Screen
}

func NewSession(fname string) (*Session, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	slines := strings.Split(string(data), "\n")
	lines := make([][]rune, len(slines))
	for i, l := range slines {
		lines[i] = []rune(l)
	}
	return &Session{File: fname, Lines: lines, scr: &Screen{LineNums: true}}, nil
}

func (s *Session) Run() error {
	s.scr.Draw(10, 20, s.Lines)

	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			if err := s.HandleKey(ev); err != nil {
				return err
			}
		}
		termbox.Flush()
	}
}

func (s *Session) HandleKey(ev termbox.Event) error {
	if ev.Ch != 0 {
		termbox.SetCell(s.CursorX, s.CursorY, ev.Ch, 0, 0)
		return nil
	}

	switch ev.Key {
	case termbox.KeyEsc:
		return ErrQuit
	}
	return nil
}


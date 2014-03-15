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
	LineShift int // which line is the first we draw
	LineNums  bool
	W, H      int
	X, Y      int // upper corner of screen
	Content   [][]rune
	lineMap   [][]int // [screenY][screenx]line#
	charMap   [][]int // [screenY][screenX]char#
}

func (s *Screen) Loc(x, y int) (line, char int) {
	return s.lineMap[y-s.Y][x-s.X], s.charMap[y-s.Y][x-s.X]
}

func (s *Screen) Draw() {
	s.clear()

	ndigits := len(fmt.Sprint(len(s.Content))) + 1

	xpos, ypos := s.X, s.Y
	if s.LineNums {
		xpos += ndigits
	}

	x, y := 0, s.LineShift // char#, line#
	wrapCount := 0
	for i := ypos; i < s.H; i++ {
		if y >= len(s.Content) {
			break
		}
		line := s.Content[y]
		for j := xpos; j < s.W; j++ {
			if x >= len(line) {
				termbox.SetCell(j, i, ' ', 0, 0)
			} else {
				termbox.SetCell(j, i, line[x], 0, 0)
				s.lineMap[i][j] = y
				s.charMap[i][j] = x
				x++
			}

		}

		if s.LineNums && wrapCount == 0{
			nums := fmt.Sprint(y + 1)
			nums = strings.Repeat(" ", ndigits-1-len(nums)) + nums + " "
			for n := 0; n < ndigits; n++ {
				termbox.SetCell(s.X+n, i-wrapCount, rune(nums[n]), 0, 0)
			}
		}


		if x >= len(line) { // if we drew entire line
			y++   // go to next line
			x = 0 // at first char
			wrapCount = 0
		} else {
			wrapCount++
		}
	}
}

func (s *Screen) clear() {
	s.lineMap = make([][]int, s.H)
	s.charMap = make([][]int, s.H)
	for i := range s.lineMap {
		s.lineMap[i] = make([]int, s.W)
		s.charMap[i] = make([]int, s.W)
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

	// initialize and draw start screen
	w, h := termbox.Size()
	scr := &Screen{
		LineNums: true,
		W: w,
		H: h,
		Content: lines,
	}

	return &Session{File: fname, Lines: lines, scr: scr}, nil
}

func (s *Session) Run() error {
	for {
		termbox.Clear(0, 0)
		s.scr.Draw()
		termbox.Flush()

		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			if err := s.HandleKey(ev); err != nil {
				return err
			}
		case termbox.EventResize:
			s.scr.W, s.scr.H = ev.Width, ev.Height
			s.scr.Draw()
		}
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

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
	CursorX   int     // cursor line#
	CursorY   int     // cursor char#
}

func (s *Screen) MovCursorX(n int) {
	line := s.Content[s.CursorY]
	s.CursorX = min(s.CursorX + n, len(line))
	s.CursorX = max(s.CursorX, 0)
}

func (s *Screen) MovCursorY(n int) {
	s.CursorY = min(s.CursorY + n, len(s.Content))
	s.CursorY = max(s.CursorY, 0)
	s.MoveCursorX(0)
}

// loc gives the line and char coordinates of the x and y (absolute) screen
// coordinates
func (s *Screen) loc(x, y int) (line, char int) {
	return s.lineMap[y][x], s.charMap[y][x]
}

func (s *Screen) Insert(ch rune) {
	s.loc(s.CursorX, s.CursorY+1)
	l, c := s.CursorY, s.CursorX
	line := s.Content[l]

	if ch == '\n' {
		head := line[:c]
		tail := append([]rune{'\n'}, line[c:]...)
		s.Content[l] = head
		s.Content = append(s.Content[:l+1], append([][]rune{tail}, s.Content[l+1:]...)...)
		s.CursorY++
		s.CursorX = 0
	} else {
		s.Content[l] = append(line[:c], append([]rune{ch}, line[c:]...)...)
		s.CursorX++
	}

	if s.CursorY > s.H {
		s.LineShift++
	}
}

func (s *Screen) Draw() {
	s.clear()

	xpos, ypos := s.X, s.Y
	ndigits := len(fmt.Sprint(len(s.Content))) + 1
	if s.LineNums {
		xpos += ndigits
	}

	// draw cursor
	termbox.SetCursor(s.X+s.CursorX+ndigits, s.Y+s.CursorY)

	// draw content
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

		if s.LineNums && wrapCount == 0 {
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
	File string
	scr  *Screen
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
		W:        w,
		H:        h,
		Content:  lines,
	}

	return &Session{
		File: fname,
		scr:  scr,
	}, nil
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
		case termbox.EventMouse:
		case termbox.EventError:
			return ev.Err
		}
	}
}

func (s *Session) HandleKey(ev termbox.Event) error {
	if ev.Ch != 0 {
		s.scr.Insert(ev.Ch)
		return nil
	}

	switch ev.Key {
	case termbox.KeyEnter:
		s.scr.Insert('\n')
	case termbox.KeySpace:
		s.scr.Insert(' ')
	case termbox.KeyArrowUp:
		s.scr.ShiftCursor(0, -1)
	case termbox.KeyArrowDown:
		s.scr.ShiftCursor(0, 1)
	case termbox.KeyArrowLeft:
		s.scr.ShiftCursor(-1, 0)
	case termbox.KeyArrowRight:
		s.scr.ShiftCursor(1, 0)
	case termbox.KeyEsc:
		return ErrQuit
	}
	return nil
}

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
	LineNums  bool
	W, H      int
	X, Y      int // upper corner of screen
	Content   [][]rune
	CursorX   int                 // cursor line#
	CursorY   int                 // cursor char#
	ypivot int
}

func (s *Screen) MovCursorX(n int) {
	line := s.Content[s.CursorY]
	s.CursorX = min(s.CursorX+n, len(line))
	s.CursorX = max(s.CursorX, 0)
}

func (s *Screen) MovCursorY(n int) {
	if s.CursorY + n >= len(s.Content) {
		s.CursorY = len(s.Content) - 1
	} else if s.CursorY + n < 0 {
		s.CursorY = 0
	} else {
		s.CursorY += n
	}

	// keep x cursor pos on text for new line
	s.MovCursorX(0)

	cv := NewWrapView(s.W - s.ndigits(), s.H, s.Content, s.CursorY, s.ypivot)
	if cv.Y(s.CursorY, 0) > 0 && cv.Y(s.CursorY, 0) < s.H-1 {
		s.ypivot = cv.Y(s.CursorY, s.CursorX)
	}
}

func (s *Screen) ndigits() int {
	if s.LineNums {
		return len(fmt.Sprint(len(s.Content))) + 1
	}
	return 0
}

func (s *Screen) Insert(ch rune) {
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
}

func (s *Screen) Resize(w, h int) {
	s.W = w
	s.H = h
}

func (s *Screen) Draw() {
	cv := NewWrapView(s.W-s.ndigits(), s.H, s.Content, s.CursorY, s.ypivot)

	// draw cursor
	x, y := RenderPos(cv, s.CursorY, s.CursorX)
	termbox.SetCursor(s.X+s.ndigits()+x, s.Y+y)

	// draw content
	for y := 0; y < s.H; y++ {
		for x := 0; x < s.W-s.ndigits(); x++ {
			line, char := DataPos(cv, x, y)
			if char != -1 {
				termbox.SetCell(s.X + x + s.ndigits(), s.Y + y, s.Content[line][char], 0, 0)
			}
		}
	}

	// draw line number
	if s.LineNums {
		prev := -1
		for y := 0; y < s.H; y++ {
			line := cv.Line(0, y)
			lg.Println(line)
			if line == -1 {
				break
			} else if line == prev {
				continue
			}

			prev = line
			nums := fmt.Sprint(line + 1)
			nums = strings.Repeat(" ", s.ndigits()-1-len(nums)) + nums + " "
			for n := 0; n < s.ndigits(); n++ {
				termbox.SetCell(s.X+n, s.Y + y, rune(nums[n]), 0, 0)
			}
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
			s.scr.Resize(ev.Width, ev.Height)
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
		s.scr.MovCursorY(-1)
	case termbox.KeyArrowDown:
		s.scr.MovCursorY(1)
	case termbox.KeyArrowLeft:
		s.scr.MovCursorX(-1)
	case termbox.KeyArrowRight:
		s.scr.MovCursorX(1)
	case termbox.KeyEsc:
		return ErrQuit
	}
	return nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

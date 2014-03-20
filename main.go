package main

import (
	"flag"
	"fmt"
	termbox "github.com/nsf/termbox-go"
	"io/ioutil"
	"log"
	"os"
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

	// start termbox
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
		lg.Print(err)
		return
	}

	// run ...
	err = s.Run()
	if err != ErrQuit {
		lg.Print(err)
	}
}

type Session struct {
	File     string
	W, H     int  // size of terminal window
	Buf      *Buffer
	View     View
	CursorL  int // cursor line#
	CursorC  int // cursor char#
	ypivot   int
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
	b := NewBuffer(data)

	w, h := termbox.Size()
	v := &LineNumView{View: &WrapView{}}
	//v := &WrapView{}
	v.SetBuf(b)
	v.SetSize(w, h)
	return &Session{
		File:     fname,
		W:        w,
		H:        h,
		Buf:      b,
		View:     v,
	}, nil
}

func (s *Session) Run() error {
	for {
		termbox.Clear(0, 0)
		s.Draw()
		termbox.Flush()

		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			if err := s.HandleKey(ev); err != nil {
				return err
			}
		case termbox.EventResize:
			s.W, s.H = ev.Width, ev.Height
			s.View.SetSize(s.W, s.H)
		case termbox.EventMouse:
		case termbox.EventError:
			return ev.Err
		}
	}
}

func (s *Session) HandleKey(ev termbox.Event) error {
	if ev.Ch != 0 {
		s.Insert(ev.Ch)
		return nil
	}

	switch ev.Key {
	case termbox.KeyEnter:
		s.Newline()
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		s.Backspace()
	case termbox.KeySpace:
		s.Insert(' ')
	case termbox.KeyArrowUp:
		s.MovCursorY(-1)
	case termbox.KeyArrowDown:
		s.MovCursorY(1)
	case termbox.KeyArrowLeft:
		s.MovCursorX(-1)
	case termbox.KeyArrowRight:
		s.MovCursorX(1)
	case termbox.KeyEsc:
		return ErrQuit
	}
	return nil
}

func (s *Session) MovCursorX(n int) {
	line := s.Buf.Line(s.CursorL)
	s.CursorC = min(s.CursorC+n, len(line))
	s.CursorC = max(s.CursorC, 0)
}

func (s *Session) MovCursorY(n int) {
	s.View.SetRef(s.CursorL, 0, 0, s.ypivot)
	cv := s.View.Render()

	if s.CursorL+n >= s.Buf.Nlines() {
		s.CursorL = s.Buf.Nlines() - 1
	} else if s.CursorL+n < 0 {
		s.CursorL = 0
	} else {
		s.CursorL += n
	}

	// keep x cursor pos on text for new line
	s.MovCursorX(0)

	// if new cursor position is on prev screen render,
	// move the cursor draw location to that screen loc
	// (i.e. don't scroll the screen)
	if Contains(cv, s.CursorL, s.CursorC) {
		s.ypivot = cv.Y(s.CursorL, s.CursorC)
	}
}

func (s *Session) Newline() {
	l, c := s.CursorL, s.CursorC
	s.Buf.Insert(s.Buf.Offset(l, c), []rune{'\n'})
	s.MovCursorY(1)
	s.CursorC = 0
}

func (s *Session) Backspace() {
	l, c := s.CursorL, s.CursorC
	offset := s.Buf.Offset(l, c)
	s.Buf.Delete(offset-1, offset)
	s.CursorL, s.CursorC = s.Buf.Pos(offset - 1)
	s.MovCursorY(0) // force refresh of scroll reference
}

func (s *Session) Insert(chs ...rune) {
	l, c := s.CursorL, s.CursorC
	s.Buf.Insert(s.Buf.Offset(l, c), chs)
	s.CursorC += len(chs)
}

func (s *Session) Draw() {
	s.View.SetRef(s.CursorL, 0, 0, s.ypivot)
	cv := s.View.Render()

	// draw cursor
	x, y := RenderPos(cv, s.CursorL, s.CursorC)
	termbox.SetCursor(x, y)

	// draw content
	for y := 0; y < s.H; y++ {
		for x := 0; x < s.W; x++ {
			termbox.SetCell(x, y, cv.Rune(x, y), 0, 0)
		}
	}
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

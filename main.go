package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	termbox "github.com/nsf/termbox-go"
	"github.com/rwcarlsen/editor/util"
	"github.com/rwcarlsen/editor/view"
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

	v := &view.LineNum{View: &view.Wrap{}}
	//v := &view.Wrap{}
	s := &Session{
		File: flag.Arg(0),
		View: v,
		ExpandTabs: false,
		Tabwidth: 4,
	}

	// run ...
	err = s.Run()
	if err != ErrQuit {
		lg.Print(err)
	}
}

type Session struct {
	File       string
	w, h       int // size of terminal window
	buf        *util.Buffer
	View       view.View
	CursorL    int // cursor line#
	CursorC    int // cursor char#
	ExpandTabs bool
	Tabwidth   int
	ypivot     int
}

func (s *Session) Run() error {
	data, err := ioutil.ReadFile(s.File)
	if err != nil {
		return err
	}
	s.buf = util.NewBuffer(data)
	s.w, s.h = termbox.Size()
	s.View.SetBuf(s.buf)
	s.View.SetSize(s.w, s.h)
	s.View.SetTabwidth(s.Tabwidth)

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
			s.w, s.h = ev.Width, ev.Height
			s.View.SetSize(s.w, s.h)
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
	case termbox.KeyTab:
		if s.ExpandTabs {
			s.Insert([]rune(strings.Repeat(" ", s.Tabwidth))...)
		} else {
			s.Insert('\t')
		}
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
	line := s.buf.Line(s.CursorL)
	s.CursorC = util.Min(s.CursorC+n, len(line))
	s.CursorC = util.Max(s.CursorC, 0)
}

func (s *Session) MovCursorY(n int) {
	s.View.SetRef(s.CursorL, 0, 0, s.ypivot)
	 surf := s.View.Render()

	if s.CursorL+n >= s.buf.Nlines() {
		s.CursorL = s.buf.Nlines() - 1
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
	if view.Contains(surf, s.CursorL, s.CursorC) {
		s.ypivot = surf.Y(s.CursorL, s.CursorC)
	}
}

func (s *Session) Newline() {
	l, c := s.CursorL, s.CursorC
	s.buf.Insert(s.buf.Offset(l, c), []rune{'\n'})
	s.MovCursorY(1)
	s.CursorC = 0
}

func (s *Session) Backspace() {
	l, c := s.CursorL, s.CursorC
	offset := s.buf.Offset(l, c)
	s.buf.Delete(offset-1, offset)
	s.CursorL, s.CursorC = s.buf.Pos(offset - 1)
	s.MovCursorY(0) // force refresh of scroll reference
}

func (s *Session) Insert(chs ...rune) {
	l, c := s.CursorL, s.CursorC
	s.buf.Insert(s.buf.Offset(l, c), chs)
	s.CursorC += len(chs)
}

func (s *Session) Draw() {
	s.View.SetRef(s.CursorL, 0, 0, s.ypivot)
	surf := s.View.Render()

	// draw cursor
	x, y := view.RenderPos(surf, s.CursorL, s.CursorC)
	termbox.SetCursor(x, y)

	// draw content
	for y := 0; y < s.h; y++ {
		for x := 0; x < s.w; x++ {
			termbox.SetCell(x, y, surf.Rune(x, y), 0, 0)
		}
	}
}


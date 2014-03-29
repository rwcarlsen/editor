package session

import (
	"fmt"
	"io/ioutil"
	"regexp"

	termbox "github.com/nsf/termbox-go"
	"github.com/rwcarlsen/editor/util"
	"github.com/rwcarlsen/editor/view"
)

var ErrQuit = fmt.Errorf("Quit")

type Mode interface {
	HandleKey(*Session, termbox.Event) (Mode, error)
}

type Session struct {
	File        string
	mode        Mode
	W, H        int // size of terminal window
	Buf         *util.Buffer
	View        view.View
	CursorL     int // cursor line#
	CursorC     int // cursor char#
	ExpandTabs  bool
	SmartIndent bool
	Search      *regexp.Regexp
	Matches     [][]int // regexp search matches
	Tabwidth    int
	Ypivot      int
}

func (s *Session) UpdSearch() {
	if s.Search == nil {
		return
	}
	s.Matches = s.Search.FindAllIndex(s.Buf.Bytes(), -1)
}

func (s *Session) Run() error {
	s.mode = &ModeEdit{}
	data, err := ioutil.ReadFile(s.File)
	if err != nil {
		return err
	}
	s.Buf = util.NewBuffer(data)
	s.W, s.H = termbox.Size()
	s.H--
	s.View.SetBuf(s.Buf)
	s.View.SetSize(s.W, s.H)
	s.View.SetTabwidth(s.Tabwidth)

	for {
		s.Draw()
		termbox.Flush()
		termbox.Clear(0, 0)

		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			s.mode, err = s.mode.HandleKey(s, ev)
			if err != nil {
				return err
			}
		case termbox.EventResize:
			s.W, s.H = ev.Width, ev.Height-1
			s.View.SetSize(s.W, s.H)
		case termbox.EventMouse:
		case termbox.EventError:
			return ev.Err
		}
	}
}

func (s *Session) SetCursor(line, char int) {
	if char < 0 {
		char = s.CursorC
	}
	if line < 0 {
		line = s.CursorL
	}

	s.View.SetRef(s.CursorL, 0, 0, s.Ypivot)
	surf := s.View.Render()

	line = util.Min(line, s.Buf.Nlines()-1)
	line = util.Max(line, 0)

	l := s.Buf.Line(line)
	char = util.Min(char, len(l)-1)
	char = util.Max(char, 0)

	if view.Contains(surf, line, char) {
		s.Ypivot = surf.Y(line, char) // don't scroll
	} else if line > s.CursorL {
		s.Ypivot = s.H - 1 // draw cursor at bottom & scroll
	} else if line < s.CursorL {
		s.Ypivot = 0 // draw cursor at top & scroll
	}
	s.CursorL = line
	s.CursorC = char
}

func (s *Session) Delete(n int) {
	offset := s.Buf.Offset(s.CursorL, s.CursorC)
	nb := s.Buf.Delete(offset, n)
	s.SetCursor(s.Buf.Pos(offset - nb))
	s.UpdSearch()
}

func (s *Session) Insert(chs ...rune) {
	offset := s.Buf.Offset(s.CursorL, s.CursorC)
	n := s.Buf.Insert(offset, chs...)
	s.SetCursor(s.Buf.Pos(offset + n))
	s.UpdSearch()
}

func (s *Session) Draw() {
	s.View.SetRef(s.CursorL, 0, 0, s.Ypivot)
	surf := s.View.Render()

	// draw cursor
	x, y := view.RenderPos(surf, s.CursorL, s.CursorC)
	termbox.SetCursor(x, y)

	// draw content
	view.Draw(surf, 0, 0)
}

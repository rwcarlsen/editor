package session

import (
	"fmt"
	"io/ioutil"

	termbox "github.com/nsf/termbox-go"
	"github.com/rwcarlsen/editor/util"
	"github.com/rwcarlsen/editor/view"
)

var ErrQuit = fmt.Errorf("Quit")

type Mode interface {
	HandleKey(*Session, termbox.Event) (Mode, error)
}

type Session struct {
	File       string
	mode       Mode
	W, H       int // size of terminal window
	Buf        *util.Buffer
	View       view.View
	CursorL    int // cursor line#
	CursorC    int // cursor char#
	ExpandTabs bool
	Tabwidth   int
	Ypivot     int
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
		termbox.Clear(0, 0)
		s.Draw()
		termbox.Flush()

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

func (s *Session) MovCursorX(n int) {
	line := s.Buf.Line(s.CursorL)
	s.CursorC = util.Min(s.CursorC+n, len(line)-1)
	s.CursorC = util.Max(s.CursorC, 0)
}

func (s *Session) MovCursorY(n int) {
	s.View.SetRef(s.CursorL, 0, 0, s.Ypivot)
	surf := s.View.Render()

	s.CursorL += n
	if s.CursorL >= s.Buf.Nlines() {
		s.CursorL = s.Buf.Nlines() - 1
	} else if s.CursorL < 0 {
		s.CursorL = 0
	}

	// keep x cursor pos on text for new line
	s.MovCursorX(0)

	// if new cursor position is on prev screen render,
	// move the cursor draw location to that screen loc
	// (i.e. don't scroll the screen)
	if view.Contains(surf, s.CursorL, s.CursorC) {
		s.Ypivot = surf.Y(s.CursorL, s.CursorC)
	}
}

func (s *Session) Newline() {
	l, c := s.CursorL, s.CursorC
	s.Buf.Insert(s.Buf.Offset(l, c), '\n')
	s.MovCursorY(1)
	s.CursorC = 0
}

func (s *Session) Backspace() {
	l, c := s.CursorL, s.CursorC
	offset := s.Buf.Offset(l, c)
	s.Buf.Delete(offset, -1)
	s.CursorL, s.CursorC = s.Buf.Pos(offset - 1)
	s.MovCursorY(0) // force refresh of scroll reference
}

func (s *Session) Insert(chs ...rune) {
	l, c := s.CursorL, s.CursorC
	s.Buf.Insert(s.Buf.Offset(l, c), chs...)
	s.CursorC += len(chs)
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

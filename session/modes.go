package session

import (
	"strings"
	"io/ioutil"

	termbox "github.com/nsf/termbox-go"
	"github.com/rwcarlsen/editor/view"
)

type ModeInsert struct{
	s *Session
}

func (m *ModeInsert) HandleKey(s *Session, ev termbox.Event) (Mode, error) {
	m.s = s
	if ev.Ch != 0 {
		s.Insert(ev.Ch)
		return m, nil
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
	case termbox.KeyCtrlS:
		err := ioutil.WriteFile(s.File, s.Buf.Bytes(), 0666)
		if err != nil {
			return m, err
		}
	case termbox.KeyEsc:
		s.MovCursorX(-1)
		return &ModeEdit{}, nil
	case termbox.KeyCtrlQ:
		return m, ErrQuit
	}
	return m, nil
}

type ModeSearch struct{
	s *Session
	view view.View
	b *util.Buffer
	pos int
}

func (m *ModeSearch) HandleKey(s *Session, ev termbox.Event) (Mode, error) {
	if m.s == nil {
		m.b = util.NewBuffer([]byte{})
		m.view = &view.Wrap{}
		m.view.SetBuf(m.b)
		m.view.SetSize(s.W, 1)
		m.view.SetTabwidth(1)
	}

	if ev.Ch != 0 {
		b.Insert(m.pos, ev.Ch)
		m.pos++
	}
}

type ModeEdit struct{
	s *Session
	prevkey rune
}

func (m *ModeEdit) HandleKey(s *Session, ev termbox.Event) (Mode, error) {
	m.s = s
	if ev.Ch != 0 {
		switch ev.Ch {
		case 'i':
			return &ModeInsert{}, nil
		case 'j':
			s.MovCursorY(1)
		case 'k':
			s.MovCursorY(-1)
		case 'l':
			s.MovCursorX(1)
		case 'h':
			s.MovCursorX(-1)
		case 'g':
			if m.prevkey == 'g' {
				m.prevkey = 0
				s.MovCursorX(-s.CursorC)
				s.MovCursorY(-s.CursorL)
			} else {
				m.prevkey = 'g'
			}
		case 'G':
			s.MovCursorX(-s.CursorC)
			s.MovCursorY(s.Buf.Nlines()-1-s.CursorL)
			s.Ypivot=s.H-1
		}
	}

	switch ev.Key {
	case termbox.KeyEnter:
		s.MovCursorY(1)
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		s.MovCursorX(-1)
	case termbox.KeySpace:
		s.MovCursorX(1)
	case termbox.KeyArrowUp:
		s.MovCursorY(-1)
	case termbox.KeyArrowDown:
		s.MovCursorY(1)
	case termbox.KeyArrowLeft:
		s.MovCursorX(-1)
	case termbox.KeyArrowRight:
		s.MovCursorX(1)
	case termbox.KeyCtrlS:
		err := ioutil.WriteFile(s.File, s.Buf.Bytes(), 0666)
		if err != nil {
			return m, err
		}
	case termbox.KeyCtrlQ:
		return m, ErrQuit
	}
	return m, nil
}


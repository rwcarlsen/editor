package session

import (
	"io/ioutil"
	"regexp"
	"strings"

	termbox "github.com/nsf/termbox-go"
	"github.com/rwcarlsen/editor/util"
	"github.com/rwcarlsen/editor/view"
)

type ModeInsert struct {
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
		space := BuildSmartIndent(s, s.CursorL)
		s.Insert('\n')
		s.Insert(space...)
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		s.Delete(-1)
	case termbox.KeySpace:
		s.Insert(' ')
	case termbox.KeyTab:
		if s.ExpandTabs {
			s.Insert([]rune(strings.Repeat(" ", s.Tabwidth))...)
		} else {
			s.Insert('\t')
		}
	case termbox.KeyArrowUp:
		s.SetCursor(s.CursorL-1, -1)
	case termbox.KeyArrowDown:
		s.SetCursor(s.CursorL+1, -1)
	case termbox.KeyArrowLeft:
		s.SetCursor(-1, s.CursorC-1)
	case termbox.KeyArrowRight:
		s.SetCursor(-1, s.CursorC+1)
	case termbox.KeyCtrlS:
		err := ioutil.WriteFile(s.File, s.Buf.Bytes(), 0666)
		if err != nil {
			return m, err
		}
	case termbox.KeyEsc:
		s.SetCursor(-1, s.CursorC-1)
		return &ModeEdit{}, nil
	case termbox.KeyCtrlQ:
		return m, ErrQuit
	}
	return m, nil
}

func BuildSmartIndent(s *Session, line int) (space []rune) {
	if !s.SmartIndent {
		return []rune{}
	}

	l := s.Buf.Line(line)
	for _, ch := range l {
		if ch == ' ' || ch == '\t' {
			space = append(space, ch)
		} else {
			break
		}
	}
	return space
}

type ModeSearch struct {
	s    *Session
	view view.View
	b    *util.Buffer
	pos  int
}

func (m *ModeSearch) HandleKey(s *Session, ev termbox.Event) (Mode, error) {
	if m.s == nil {
		m.s = s
		m.b = util.NewBuffer([]byte{})
		m.view = &view.Wrap{}
		m.view.SetBuf(m.b)
		m.view.SetSize(s.W-1, 1)
		m.view.SetTabwidth(1)
	}

	var err error
	if ev.Ch != 0 {
		m.b.Insert(m.pos, ev.Ch)
		m.pos++
	}
	switch ev.Key {
	case termbox.KeyEnter:
		s.Search, err = regexp.Compile(string(m.b.Bytes()))
		if err != nil {
			msg := err.Error()
			for i, ch := range msg {
				termbox.SetCell(i, s.H, ch, 0, 0)
			}
			return &ModeEdit{}, nil
		}

		s.UpdSearch()
		if len(s.Matches) > 0 {
			cursor := s.Buf.Offset(s.CursorL, s.CursorC)
			n := 0
			for i, match := range s.Matches {
				offset := match[0]
				if offset >= cursor {
					n = i
					break
				}

			}
			offset := s.Matches[n][0]
			s.SetCursor(s.Buf.Pos(offset))
		}
		return &ModeEdit{}, nil
	case termbox.KeySpace:
		m.b.Insert(m.pos, ' ')
		m.pos++
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		m.b.Delete(m.pos, -1)
		m.pos--
	case termbox.KeyEsc:
		return &ModeEdit{}, nil
	}

	surf := m.view.Render()
	termbox.SetCell(0, s.H, '/', 0, 0)
	view.Draw(surf, 1, s.H)
	return m, nil
}

type ModeEdit struct {
	s       *Session
	prevkey rune
}

func (m *ModeEdit) HandleKey(s *Session, ev termbox.Event) (Mode, error) {
	m.s = s

	switch m.prevkey {
	case 'g':
		switch ev.Ch {
		case 'g':
			if m.prevkey == 'g' {
				m.prevkey = 0
				s.SetCursor(0, 0)
			}
		default:
			m.prevkey = 0
		}
	case 0:
		switch ev.Ch {
		case 'i':
			return &ModeInsert{}, nil
		case 'j':
			s.SetCursor(s.CursorL+1, -1)
		case 'k':
			s.SetCursor(s.CursorL-1, -1)
		case 'l':
			s.SetCursor(-1, s.CursorC+1)
		case 'h':
			s.SetCursor(-1, s.CursorC-1)
		case 'o':
			l := s.Buf.Line(s.CursorL)
			s.SetCursor(-1, len(l)-1)
			space := BuildSmartIndent(s, s.CursorL)
			s.Insert('\n')
			s.Insert(space...)
			return &ModeInsert{}, nil
		case 'x':
			s.Delete(1)
			s.SetCursor(-1, s.CursorC+1)
		case 'g':
			m.prevkey = 'g'
		case 'G':
			s.SetCursor(s.Buf.Nlines()-1, 0)
		case '/':
			termbox.SetCell(0, s.H, '/', 0, 0)
			return &ModeSearch{}, nil
		case 'n':
			if len(s.Matches) == 0 {
				break
			}

			cursor := s.Buf.Offset(s.CursorL, s.CursorC)
			n := 0
			for i, match := range s.Matches {
				offset := match[0]
				if offset > cursor {
					n = i
					break
				}

			}
			offset := s.Matches[n][0]
			s.SetCursor(s.Buf.Pos(offset))
		}
	}

	switch ev.Key {
	case termbox.KeyArrowUp:
		s.SetCursor(s.CursorL-1, -1)
	case termbox.KeyArrowDown, termbox.KeyEnter:
		s.SetCursor(s.CursorL+1, -1)
	case termbox.KeyArrowLeft, termbox.KeyBackspace, termbox.KeyBackspace2:
		s.SetCursor(-1, s.CursorC-1)
	case termbox.KeyArrowRight, termbox.KeySpace:
		s.SetCursor(-1, s.CursorC+1)
	case termbox.KeyCtrlS:
		err := ioutil.WriteFile(s.File, s.Buf.Bytes(), 0666)
		if err != nil {
			return m, err
		}
	case termbox.KeyEsc:
		m.prevkey = 0
	case termbox.KeyCtrlQ:
		return m, ErrQuit
	}
	return m, nil
}

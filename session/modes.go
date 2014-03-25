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

	if ev.Ch != 0 {
		m.b.Insert(m.pos, ev.Ch)
		m.pos++
	}
	switch ev.Key {
	case termbox.KeyEnter:
		re, err := regexp.Compile(string(m.b.Bytes()))
		if err != nil {
			msg := err.Error()
			for i, ch := range msg {
				termbox.SetCell(i, s.H, ch, 0, 0)
			}
			return &ModeEdit{}, nil
		}

		matches := re.FindAllIndex(s.Buf.Bytes(), -1)
		n := 0
		if len(matches) > 0 {
			cursor := s.Buf.Offset(s.CursorL, s.CursorC)
			for i, match := range matches {
				offset := match[0]
				if offset >= cursor {
					n = i
					break
				}

			}
			offset := matches[n][0]
			l, c := s.Buf.Pos(offset)
			s.MovCursorY(-s.CursorL + l)
			s.MovCursorX(-s.CursorC + c)
		}
		return &ModeEdit{Search: matches, SearchN: n}, nil
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
	Search  [][]int
	SearchN int
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
				s.MovCursorX(-s.CursorC)
				s.MovCursorY(-s.CursorL)
			}
		default:
			m.prevkey = 0
		}
	case 0:
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
			s.MovCursorY(s.Buf.Nlines() - 1 - s.CursorL)
			s.Ypivot = s.H - 1
		case '/':
			termbox.SetCell(0, s.H, '/', 0, 0)
			return &ModeSearch{}, nil
		case 'n':
			if len(m.Search) == 0 {
				break
			} else if m.SearchN++; m.SearchN >= len(m.Search) {
				m.SearchN = 0
			}
			offset := m.Search[m.SearchN][0]
			l, c := s.Buf.Pos(offset)
			s.MovCursorY(-s.CursorL + l)
			s.MovCursorX(-s.CursorC + c)
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
	case termbox.KeyEsc:
		m.prevkey = 0
	case termbox.KeyCtrlQ:
		return m, ErrQuit
	}
	return m, nil
}

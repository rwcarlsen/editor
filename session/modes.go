package session

import (
	"strings"
	"io/ioutil"

	termbox "github.com/nsf/termbox-go"
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
		err := ioutil.WriteFile(s.File, s.buf.Bytes(), 0666)
		if err != nil {
			return m, err
		}
	case termbox.KeyCtrlQ:
		return m, ErrQuit
	}
	return m, nil
}


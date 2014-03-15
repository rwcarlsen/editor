package main

import (
	"flag"
	"fmt"
	termbox "github.com/nsf/termbox-go"
	"log"
	"os"
)

var ErrQuit = fmt.Errorf("Quit")

func main() {
	flag.Parse()
	log.SetFlags(0)

	// start terminal
	err := termbox.Init()
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
	Cells []termbox.Cell
	Offsets []int
	Lines []string
	Xshift int
	Yshift int
}

func (s *Screen) Draw() error {
	termbox.Clear()
	termbox.Flush()
	w, y := termbox.Size()
	for i, line := range s.Lines[Yshift:Yshift + h] {

		for j, line := range line[Xshift:Xshift + w] {
		}
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
	}

}

	Lines []rune

type Session struct {
	CursorX int
	CursorY int
	File string
	Lines []string
}

func NewSession(fname string) (*Session, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	text := string(ioutil.ReadFile(fname))
	lines := strings.Split(text, "\n")
	return &Session{File: fname, Lines: lines}, nil
}

func (s *Session) Run() error {
	defer s.f.Close()
	s.LoadFile()

	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			if err := HandleKey(ev); err != nil {
				return err
			}
		}
		termbox.Flush()
	}
}

func (s *Session) FillScreen() error {
	w, h := termbox.Size()

	f.Read(

	for err != io.EOF

	data := make([]byte, w * h)
	err := io.ReadFull(data, s.f)
	if err != nil {
		return err
	}
	x, y := 0, 0
	for i, b := range data {
		termbox.SetCell(s.CursorX, s.CursorY, ev.Ch, 0, 0)
	}
}

func (s *Session) HandleKey(ev termbox.Event) error {
	if ev.Ch != 0 {
		termbox.SetCell(s.CursorX, s.CursorY, ev.Ch, 0, 0)
		return nil
	}

	switch ev.Key {
	case termbox.KeyEsc:
		return ErrQuit
	}
	return nil
}


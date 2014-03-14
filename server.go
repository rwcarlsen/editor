package server

import (
	"code.google.com/p/goncurses"
	"fmt"
	"os"
	"math/rand"
)

type Screen struct {
	Chars  [][]goncurses.Char
}

type ScreenUpd struct {
	X, Y int
	Char goncurses.Char
}

func NewScreen(w, h int) *Screen {
	s := &Screen{}
	s.Chars = make([][]goncurses.Char, h)
	for i := range s.Chars {
		s.Chars[i] = make([]goncurses.Char, w)
	}
	return s
}

type Client struct {
}

type Server struct {
	SessionAddr string
	Sessions map[string]*Session
}

func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go s.handle(c)
	}
}

type ServMsg struct {
	Type string
	File string
}

func (s *Server) handle(c net.Conn) {
	var m ServMsg
	err := DecodeMsg(&m, c)
	if err != nil {
		log.Print(err)
		return
	}

	switch m.Type {
	case "new-client"
		name := m.File
		s, ok := s.Sessions[m.File]
		if ok {
			name = fmt.Sprintf("%v-%v", name, rand.Int())
			s = &Session{serv: s, Name: name, File: m.File}
			if err := s.Open(); err != nil {
				log.Print(err)
				return
			}
			s.Sessions[m.File] = s
		}
		s.client = c
	}
}

func (s *Server) ListenAndServe() error {
}

type Session struct {
	Name string
	File string
	f *os.File
	client net.Conn
	serv *Server
	W, H int
}

func DecodeMsg(v interface{}, r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}

type ClientMsg struct {
	Type string
	W, H int
	Key goncurses.Key
}

func (s *Session) listenAndServe() {
	for {
		var m ClientMsg
		err := DecodeMsg(&m, s.client)
		if err == io.EOF {
			log.Print("connection closed by client unexpectedly")
			return
		} else if err != nil {
			log.Print(err)
			continue
		}

		// send key+other events to server - server should return a list of
		// screen updates if any, and a list of file updates if any
	}
}

func (s *Session) UpdateScreen(scr *Screen) error {
}

func (s *Session) UpdateFile(data []byte, offset int) error {
	f.WriteAt(data, offset)
	return f.Sync()
}

func (s *Session) Open(client net.Conn) (err error) {
	s.client = client
	var m ClientMsg
	err := DecodeMsg(&m, client)
	if err != nil {
		return err
	}
	s.W, s.H = m.W, m.H

	s.f, err = os.Open(s.File)
	if err != nil {
		return err
	}

	go s.listenAndServe()
	return nil
}

func (s *Session) Close() error {
	s.client.Close()
	return s.f.Close()
}


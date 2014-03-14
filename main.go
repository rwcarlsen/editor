package main

import (
	"code.google.com/p/goncurses"
	"fmt"
	"os"
)

func main() {
	win, err := goncurses.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer goncurses.End()

	win.Print("hello from ncurses!")
	win.Move(10, 20)
	win.Print("hello from a moved cursor!")
	win.Move(20, 30)
	win.AttrOn(goncurses.A_BOLD)
	win.Print("hello in bold!")
	win.Move(21, 30)
	win.Print("hello in bold again!")
	win.Move(25, 30)
	win.AddChar('z')
	_ = win.GetChar()
	win.Refresh()
}

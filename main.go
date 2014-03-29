package main

import (
	"flag"
	"log"
	"os"

	termbox "github.com/nsf/termbox-go"
	"github.com/rwcarlsen/editor/session"
	"github.com/rwcarlsen/editor/view"
)

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
	s := &session.Session{
		File:        flag.Arg(0),
		View:        v,
		ExpandTabs:  false,
		SmartIndent: true,
		Tabwidth:    4,
	}

	// run ...
	err = s.Run()
	if err != session.ErrQuit {
		lg.Print(err)
	}
}

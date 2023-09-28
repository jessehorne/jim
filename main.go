package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jessehorne/jim/jim"
	"log"
)

func handleInput(s tcell.Screen) {
L:
	for {
		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			fmt.Println(ev.Key())
			if ev.Key() == tcell.KeyEscape {
				break L
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalln(err)
	}

	if err = s.Init(); err != nil {
		log.Fatalln(err)
	}

	s.SetStyle(tcell.StyleDefault.Foreground(jim.ColorWhite).Background(jim.ColorBlack))
	s.Clear()

	s.Show()
	s.Sync()

	// create folder/file viewer
	newFv := jim.NewFv(s)
	newFv.Refresh()
	newFv.Redraw()

	// begin loop
	handleInput(s)

	s.Fini()

	fmt.Println("Thanks.")
}

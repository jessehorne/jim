package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jessehorne/jim/jim"
	"log"
)

func handleInput(fv *jim.Fv) {
L:
	for {
		ev := fv.Screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				break L
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			buttons := ev.Buttons()
			fv.ButtonEvent(x, y, buttons)
		case *tcell.EventResize:
			fv.Width, fv.Height = fv.Screen.Size()
			fv.RefreshTree()
			fv.DrawBackground()
			fv.PrintTree()
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
	s.EnableMouse()

	s.Clear()
	s.Show()
	s.Sync()

	// create folder/file viewer
	newFv := jim.NewFv(s)
	newFv.ExpandDir(nil) // expand ./ in the main tree view
	newFv.RefreshTree()
	newFv.PrintTree()

	// begin loop
	handleInput(newFv)

	s.Fini()

	fmt.Println("Thanks.")
}

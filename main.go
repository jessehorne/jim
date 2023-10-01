package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/jim/jim"
	"log"
)

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalln(err)
	}

	if err = s.Init(); err != nil {
		log.Fatalln(err)
	}

	s.SetStyle(tcell.StyleDefault.Background(jim.ColorDark).Foreground(jim.ColorWhite))
	s.EnableMouse()
	s.Clear()
	s.Show()
	s.Sync()

	newManager := jim.NewManager(s)
	newManager.Init()
	newManager.HandleInput()

	s.Fini()

	fmt.Println("Thanks.")
}

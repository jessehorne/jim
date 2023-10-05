package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/jim/jim"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Please provide a path to the directory you'd like to open.")
		return
	}

	dir := os.Args[1]

	f, err := os.Open(dir)
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if !fileInfo.IsDir() {
		log.Fatalln("You must provide a path to a directory. You can't edit single files yet and also please check your path to make sure it exists.")
	}

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
	newManager.Init(dir)
	newManager.HandleInput()

	s.Fini()

	fmt.Println("Thanks.")
}

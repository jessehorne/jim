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
		fmt.Println("Provide a path to a directory or file.")
		fmt.Println("Try '--help' for more information.")
		return
	}

	dir := os.Args[1]

	if dir == "--help" {
		fmt.Println(jim.HelpMessage)
		return
	}

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
		fmt.Println("Invalid directory path.")
		fmt.Println("Try `--help` for more information.")
		return
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

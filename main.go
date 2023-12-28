package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/jessehorne/jim/jim"
	"github.com/rivo/tview"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Must specify a directory to open.")
		os.Exit(1)
	}

	dir := os.Args[1]

	dirFile, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer dirFile.Close()

	info, err := dirFile.Stat()
	if err != nil {
		panic(err)
	}

	if !info.IsDir() {
		panic("Not a directory.")
	}

	app := tview.NewApplication()

	editor := jim.NewEditor(app, dirFile)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			return nil
		} else if event.Key() == tcell.KeyCtrlS {
			editor.SaveCurrentTab()
		}

		return event
	})

	if err := app.SetRoot(editor.Grid, true).SetFocus(editor.Grid).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

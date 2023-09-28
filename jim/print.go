package jim

import "github.com/gdamore/tcell"

func Print(scr tcell.Screen, st tcell.Style, x int, y int, str string) {
	for i := 0; i < len(str); i++ {
		scr.SetCell(x+i, y, st, rune(str[i]))
	}
}

package jim

import (
	"github.com/gdamore/tcell/v2"
	"os"
	"strconv"
)

const (
	CursorDirUp = iota
	CursorDirDown
	CursorDirLeft
	CursorDirRight
)

var Typeables = []rune{
	'A',
}

type Tab struct {
	Screen    tcell.Screen
	Manager   *Manager
	File      *File
	Active    bool
	Width     int
	Height    int
	ScrollY   int // the number of bytes scrolled down in the document
	ScrollX   int // the number of characters scrolled right
	OffsetX   int
	OffsetY   int
	Content   []string
	LineCount int
	CursorX   int
	CursorY   int
}

func NewTab(screen tcell.Screen, m *Manager) *Tab {
	return &Tab{
		Screen:  screen,
		Manager: m,
	}
}

func (t *Tab) SetContent(s string) {
	var temp string
	var lineCount int
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			t.Content = append(t.Content, temp)
			temp = ""
			lineCount++
		} else {
			temp = temp + string(s[i])
		}
	}
	t.LineCount = lineCount + 1
}

func (t *Tab) TypeCharacter(c rune) {
	return
}

func (t *Tab) LoadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	t.SetContent(string(data))

	return nil
}

func (t *Tab) GetCursorLine() int {
	return t.CursorY + t.OffsetY - 1
}

func (t *Tab) ScrollUp() {
	if t.CursorY > 0 {
		if t.OffsetY < 1 {
			t.CursorY--
			t.OffsetY++
		} else {
			t.CursorY--
		}

		// get above line
		if t.GetCursorLine() > 0 {
			if t.CursorX > len(t.Content[t.GetCursorLine()]) {
				t.CursorX = len(t.Content[t.GetCursorLine()])
			}
		}

		t.Redraw()
	}
}

func (t *Tab) ScrollDown() {
	if t.CursorY <= t.Height-3 {
		t.CursorY++
	} else {
		if t.CursorY < len(t.Content) {
			t.OffsetY--
			t.CursorY++
			t.Redraw()
		}
	}

	if t.GetCursorLine() < len(t.Content) {
		if t.CursorX > len(t.Content[t.GetCursorLine()]) {
			t.CursorX = len(t.Content[t.GetCursorLine()])
			t.Redraw()
		}
	}
}

func (t *Tab) ScrollRight() {
	if t.CursorX < len(t.Content[t.GetCursorLine()]) {
		t.CursorX++
	} else {
		if t.GetCursorLine()+1 < len(t.Content) {
			t.CursorY++
			t.CursorX = 0
		}
	}
}

func (t *Tab) ScrollLeft() {
	if t.CursorX > 0 {
		t.CursorX--
	} else {
		if t.GetCursorLine() > 0 {
			t.CursorY--
			t.CursorX = len(t.Content[t.GetCursorLine()])
		}
	}
}

func (t *Tab) MoveCursor(dir int) {
	if dir == CursorDirUp {
		t.ScrollUp()
	} else if dir == CursorDirDown {
		t.ScrollDown()
	} else if dir == CursorDirRight {
		t.ScrollRight()
	} else if dir == CursorDirLeft {
		t.ScrollLeft()
	}

	t.Redraw()
	t.Screen.ShowCursor(t.OffsetX+t.CursorX, t.OffsetY+t.CursorY)
	t.Screen.Sync()
}

func (t *Tab) Redraw() {
	t.Manager.ClearEditor()

	x, y := t.OffsetX, t.OffsetY

	for lineNumber := 0 - t.OffsetY; lineNumber < t.Height-t.OffsetY; lineNumber++ {
		// draw line number
		lineString := strconv.FormatInt(int64(lineNumber), 10)

		// draw bg of line number
		for i := 0; i < 7; i++ {
			t.Screen.SetCell(x+i-7, y+lineNumber, StyleLineNumber, ' ')
		}

		for i := len(lineString) - 1; i >= 0; i-- {
			t.Screen.SetCell(x-i-2, y+lineNumber, StyleLineNumber, rune(lineString[len(lineString)-i-1]))
		}

		if lineNumber >= 0 && lineNumber < len(t.Content) {
			line := t.Content[lineNumber]

			var word string
			for i := 0; i < len(line); i++ {
				t.Screen.SetCell(x+i, y+lineNumber, StyleEditor, rune(line[i]))

				c := line[i]

				if c == ' ' || c == '\t' {
					if len(word) > 0 {
						PrintWord(t.Screen, word, x+i-len(word), y+lineNumber)
						word = ""
					}

					t.Screen.SetCell(x+i, y+lineNumber, StyleEditor, rune(c))
				} else {
					word = word + string(c)
				}
			}
		}
	}

	// draw cursor
	t.Screen.SetCursorStyle(tcell.CursorStyleBlinkingBar)
	t.Screen.ShowCursor(t.OffsetX+t.CursorX, t.OffsetY+t.CursorY)

	t.Screen.Sync()
}

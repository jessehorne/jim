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
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-', '=',
	'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+',
	'[', ']', '\\', ';', '\'', ',', '.', '/',
	'{', '}', '|', ':', '"', '<', '>', '?',
	'`', '~',
	' ', '\t',
}

var TypeablesMap map[rune]bool

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
	Edited    bool
}

func NewTab(screen tcell.Screen, m *Manager) *Tab {
	return &Tab{
		Screen:  screen,
		Manager: m,
	}
}

func (t *Tab) SaveToFile() {
	f, err := os.Create(t.File.FullPath)
	if err != nil {
		return
	}

	for _, l := range t.Content {
		f.WriteString(l + "\n")
	}

	if t.Content[len(t.Content)-1] != "" {
		t.Content = append(t.Content, "")
		t.Redraw()
	}

	if err := f.Close(); err != nil {
		return
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

	t.Content = append(t.Content, temp)

	t.LineCount = lineCount + 1
}

func (t *Tab) TypeCharacter(c rune) {
	t.Edited = true

	_, ok := TypeablesMap[c]

	if !ok {
		return
	}

	t.CursorX++

	newLine := t.Content[t.GetCursorLine()][:t.CursorX-1] + string(c) + t.Content[t.GetCursorLine()][t.CursorX-1:]
	t.Content[t.GetCursorLine()] = newLine
	t.Redraw()
}

func (t *Tab) Backspace() {
	if t.CursorX > 0 {
		firstPart := t.Content[t.GetCursorLine()][:t.CursorX-1]
		nextPart := t.Content[t.GetCursorLine()][t.CursorX:]
		t.Content[t.GetCursorLine()] = firstPart + nextPart
		t.CursorX--
		t.Redraw()
	} else {
		if t.GetCursorLine() > 0 {
			t.MoveCursor(CursorDirUp)
			t.CursorX = len(t.Content[t.GetCursorLine()])
			t.Content[t.GetCursorLine()] = t.Content[t.GetCursorLine()] + t.Content[t.GetCursorLine()+1]
			t.Content = append(t.Content[:t.GetCursorLine()+1], t.Content[t.GetCursorLine()+2:]...)
			t.Redraw()
		}
	}
}

func (t *Tab) ScrollTo(pos int) {
	t.ScrollY = pos
	t.Redraw()
}

func (t *Tab) PageUp() {
	if t.ScrollY < 0 {
		t.ScrollTo(t.ScrollY + t.Height - 4)
		if t.ScrollY > 0 {
			t.ScrollTo(0)
		}
	}
}

func (t *Tab) PageDown() {
	if t.ScrollY > -len(t.Content) {
		t.ScrollTo(t.ScrollY - t.Height + 4)
		if t.ScrollY < -len(t.Content) {
			t.ScrollTo(-len(t.Content))
		}
	}
}

func (t *Tab) Newline() {
	index := t.GetCursorLine()

	if len(t.Content) == index {
		t.Content = append(t.Content, " ")
	}

	dataFirst := t.Content[t.GetCursorLine()][:t.CursorX]
	dataSecond := t.Content[t.GetCursorLine()][t.CursorX:]

	firstHalf := t.Content[:t.GetCursorLine()+1]
	secondHalf := t.Content[t.GetCursorLine():]

	t.Content = append(firstHalf, secondHalf...)

	t.Content[t.GetCursorLine()] = dataFirst
	t.Content[t.GetCursorLine()+1] = dataSecond

	t.CursorX = 0
	t.MoveCursor(CursorDirDown)
	t.Redraw()
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
	return -t.ScrollY + t.CursorY
}

func (t *Tab) ScrollUp() {
	if t.GetCursorLine() == 0 {
		return
	}

	if t.CursorY == 0 {
		if t.ScrollY < 0 {
			t.ScrollY++
			t.Redraw()
		}
	} else {
		t.CursorY--
	}

	if t.GetCursorLine() > 0 {
		if t.CursorX > len(t.Content[t.GetCursorLine()]) {
			t.CursorX = len(t.Content[t.GetCursorLine()])
		}
	}
}

func (t *Tab) ScrollDown() {
	if t.GetCursorLine() == len(t.Content)-1 {
		return
	}

	if t.CursorY < t.Height-2 {
		t.CursorY++
	} else {
		if t.GetCursorLine() < len(t.Content)-1 {
			t.ScrollY--
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

func (t *Tab) RedrawLine(line int) {
	// draw clear line
	l := t.Content[t.GetCursorLine()]

	x, y := t.OffsetX+t.ScrollX, t.OffsetY

	for i := 0; i < t.Width; i++ {
		t.Screen.SetCell(x+i, y+t.GetCursorLine(), StyleEditor, ' ')
	}

	// draw letters
	var word string
	for i := 0; i < len(l); i++ {
		t.Screen.SetCell(x+i, y+t.GetCursorLine(), StyleEditor, rune(l[i]))

		c := l[i]

		if c == ' ' || c == '\t' {
			if len(word) > 0 {
				PrintWord(t.Screen, word, x+i-len(word), y+t.GetCursorLine())
				word = ""
			}

			t.Screen.SetCell(x+i, y+t.GetCursorLine(), StyleEditor, rune(c))
		} else {
			word = word + string(c)
		}
	}

	// draw cursor
	t.Screen.SetCursorStyle(tcell.CursorStyleBlinkingBar)
	t.Screen.ShowCursor(t.OffsetX+t.CursorX, t.OffsetY+t.CursorY)

	t.Screen.Sync()
}

func (t *Tab) Redraw() {
	t.Manager.ClearEditor()

	x, y := t.OffsetX, t.OffsetY+t.ScrollY

	var max int

	if len(t.Content) < -t.ScrollY+t.Height-1 {
		max = len(t.Content)
	} else {
		max = -t.ScrollY + t.Height - 1
	}

	for lineNumber := -t.ScrollY; lineNumber < max; lineNumber++ {
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

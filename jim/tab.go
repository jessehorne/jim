package jim

import (
	"github.com/gdamore/tcell/v2"
)

type Tab struct {
	Screen  tcell.Screen
	Manager *Manager
	File    *File
	Active  bool
	Width   int
	Height  int
	ScrollY int // the number of bytes scrolled down in the document
	ScrollX int // the number of characters scrolled right
	OffsetX int
	OffsetY int
	Content string
	CursorX int
	CursorY int
}

func NewTab(screen tcell.Screen, m *Manager) *Tab {
	return &Tab{
		Screen:  screen,
		Manager: m,
	}
}

func (t *Tab) SetContent(s string) {
	t.Content = s
}

func (t *Tab) LoadFile(path string) error {
	return nil
}

func (t *Tab) Redraw() {
	if t.ScrollY >= len(t.Content) {
		return
	}

	t.Manager.ClearEditor()

	style := tcell.StyleDefault.Background(ColorDark).Foreground(ColorWhite)

	x, y := t.OffsetX, t.OffsetY
	for _, c := range t.Content {
		t.Screen.SetCell(x, y, style, c)

		if c == '\n' {
			y++
			x = t.OffsetX
		} else {
			x++
		}
	}

	// draw cursor
	t.Screen.SetCursorStyle(tcell.CursorStyleBlinkingBar)
	t.Screen.ShowCursor(t.OffsetX+t.CursorX, t.OffsetY+t.CursorY)

	t.Screen.Sync()
}

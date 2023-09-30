package jim

import "github.com/gdamore/tcell"

type Tab struct {
	Screen  tcell.Screen
	File    *File
	Active  bool
	Width   int
	Height  int
	ScrollY int // the number of bytes scrolled down in the document
	ScrollX int // the number of characters scrolled right
	OffsetX int
	OffsetY int
	Content string
}

func NewTab(screen tcell.Screen) *Tab {
	return &Tab{
		Screen: screen,
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

	style := tcell.StyleDefault.Background(ColorBlack).Foreground(ColorWhite)

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
}

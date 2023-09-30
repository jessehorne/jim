package jim

import (
	"fmt"
	"github.com/gdamore/tcell"
	"os"
)

type Manager struct {
	Screen tcell.Screen
	Fv     *Fv
	Tabs   []*Tab
}

func NewManager(s tcell.Screen) *Manager {
	return &Manager{
		Screen: s,
		Tabs:   []*Tab{},
	}
}

func (m *Manager) Init() {
	w, h := m.Screen.Size()

	// create folder/file viewer
	m.Fv = NewFv(m.Screen)
	m.Fv.ExpandDir(nil) // expand ./ in the main tree view
	m.Fv.RefreshTree()

	// create initial welcome tab
	data, err := os.ReadFile("./welcome.txt")

	if err != nil {
		return
	}
	newTab := NewTab(m.Screen, m)
	newTab.Width = w - m.Fv.Width
	newTab.Height = h
	newTab.OffsetX = m.Fv.Width
	newTab.OffsetY = 1
	newTab.SetContent(string(data))
	newTab.Redraw()

	m.RedrawTabs()
}

func (m *Manager) ClearScreen() {
	w, h := m.Screen.Size()
	blackStyle := tcell.StyleDefault.Background(ColorBlack).Foreground(ColorWhite)

	for y := 1; y < h; y++ {
		for x := m.Fv.WallX; x < w; x++ {
			m.Screen.SetCell(x, y, blackStyle, ' ')
		}
	}
}

func (m *Manager) RedrawTabs() {
	w, _ := m.Screen.Size()

	// draw black line up top
	blackStyle := tcell.StyleDefault.Background(ColorDarkBlack).Foreground(ColorBlack)
	for i := m.Fv.WallX; i < w; i++ {
		m.Screen.SetCell(i, 0, blackStyle, ' ')
	}

	currentX := m.Fv.WallX
	style := tcell.StyleDefault.Background(ColorWhite).Foreground(ColorGrey)
	for _, t := range m.Tabs {
		label := fmt.Sprintf(" %s ", t.File.Name)
		for _, l := range label {
			m.Screen.SetCell(currentX, 0, style, l)
			currentX++
		}
		currentX++                                // add one more space for padding between tabs
		m.Screen.SetCell(currentX, 0, style, ' ') // draw black space between tabs
	}

	m.Screen.Sync()
}

func (m *Manager) OpenTab(f *File) {
	var t *Tab
	for i := 0; i < len(m.Tabs); i++ {
		if m.Tabs[i].File == f {
			t = m.Tabs[i]
			break
		}
	}

	if t == nil {
		w, h := m.Screen.Size()

		data, err := os.ReadFile(f.FullPath)
		if err != nil {
			return
		}

		// open new tab and print content from 0 scroll y
		newTab := NewTab(m.Screen, m)
		newTab.Width = w - m.Fv.WallX
		newTab.Height = h
		newTab.OffsetX = m.Fv.WallX
		newTab.OffsetY = 1
		newTab.File = f
		newTab.SetContent(string(data))
		newTab.Redraw()

		m.Tabs = append(m.Tabs, newTab)

		m.RedrawTabs()
	} else {
		t.Redraw()
		m.RedrawTabs()
	}

	// switch to current tab
}

func (m *Manager) ButtonEvent(x int, y int, buttons tcell.ButtonMask) {
	switch buttons {
	case tcell.Button1:
		f := m.Fv.GetFileAtXY(x, y)

		if f == nil {
			return
		}

		if f.Type == FileTypeDir {
			if f.Expanded {
				m.Fv.UnexpandDir(f)
			} else {
				m.Fv.ExpandDir(f)
			}
			m.Fv.DrawBackground()
			m.Fv.RefreshTree()
			m.Fv.PrintTree()
		} else if f.Type == FileTypeFile {
			m.OpenTab(f)
		}
	}
}

func (m *Manager) HandleInput() {
L:
	for {
		ev := m.Screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				break L
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			buttons := ev.Buttons()
			m.ButtonEvent(x, y, buttons)
		case *tcell.EventResize:
			m.Fv.Width, m.Fv.Height = m.Screen.Size()
			m.Fv.RefreshTree()
			m.Fv.DrawBackground()
			m.Fv.PrintTree()
		}
	}
}

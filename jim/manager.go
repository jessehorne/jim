package jim

import (
	"fmt"
	"github.com/gdamore/tcell"
	"os"
)

type Manager struct {
	Screen        tcell.Screen
	Fv            *Fv
	Tabs          []*Tab
	LastActiveTab *Tab
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
	newTab.OffsetX = m.Fv.WallX
	newTab.OffsetY = 1
	newTab.SetContent(string(data))
	newTab.Redraw()

	m.RedrawTabLabels()
}

func (m *Manager) ClearEditor() {
	w, h := m.Screen.Size()
	blackStyle := tcell.StyleDefault.Background(ColorDark).Foreground(ColorWhite)

	for y := 1; y < h; y++ {
		for x := m.Fv.WallX; x < w; x++ {
			m.Screen.SetCell(x, y, blackStyle, ' ')
		}
	}

	m.Screen.Sync()
}

func (m *Manager) RedrawTabLabels() {
	w, _ := m.Screen.Size()

	// draw black line up top
	blackStyle := tcell.StyleDefault.Background(ColorGrey).Foreground(ColorWhite)
	for i := m.Fv.WallX; i < w; i++ {
		m.Screen.SetCell(i, 0, blackStyle, ' ')
	}

	currentX := m.Fv.WallX
	style := tcell.StyleDefault.Background(ColorGrey).Foreground(ColorWhite)
	activeStyle := tcell.StyleDefault.Background(ColorGrey).Foreground(ColorLightBlue)
	for _, t := range m.Tabs {
		label := fmt.Sprintf(" %s ", t.File.Name)
		for _, l := range label {
			if t.Active {
				m.Screen.SetCell(currentX, 0, activeStyle, l)
			} else {
				m.Screen.SetCell(currentX, 0, style, l)
			}
			currentX++
		}
		m.Screen.SetCell(currentX, 0, style, ' ') // draw black space between tabs
		currentX += 2                             // add one more space for padding between tabs
	}

	m.Screen.Sync()
}

func (m *Manager) CloseTab(t *Tab) {
	var removedIndex = -1
	for i := 0; i < len(m.Tabs); i++ {
		if m.Tabs[i] == t {
			m.Tabs = append(m.Tabs[:i], m.Tabs[i+1:]...)
			removedIndex = i
			break
		}
	}

	if len(m.Tabs) > 0 {
		if removedIndex > 0 {
			m.OpenTab(m.Tabs[removedIndex-1].File)
		} else {
			m.ClearEditor()
			m.LastActiveTab = nil
		}
	} else {
		m.ClearEditor()
		m.RedrawTabLabels()
		m.LastActiveTab = nil
	}
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
		newTab.Active = true
		newTab.Redraw()

		if m.LastActiveTab != nil {
			m.LastActiveTab.Active = false
		}
		m.LastActiveTab = newTab

		m.Tabs = append(m.Tabs, newTab)

		m.RedrawTabLabels()
	} else {
		if m.LastActiveTab != nil {
			m.LastActiveTab.Active = false
		}
		m.LastActiveTab = t

		t.Active = true
		t.Redraw()
		m.RedrawTabLabels()
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
			} else if ev.Key() == tcell.KeyCtrlW {
				m.CloseTab(m.LastActiveTab)
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

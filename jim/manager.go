package jim

import "github.com/gdamore/tcell"

type Manager struct {
	Screen tcell.Screen
	Fv     *Fv
	Tabs   []*Tab
}

func NewManager(s tcell.Screen) *Manager {
	return &Manager{
		Screen: s,
	}
}

func (m *Manager) Init() {
	w, h := m.Screen.Size()

	// create folder/file viewer
	m.Fv = NewFv(m.Screen)
	m.Fv.ExpandDir(nil) // expand ./ in the main tree view
	m.Fv.RefreshTree()

	// create initial welcome tab
	newTab := NewTab(m.Screen)
	newTab.Width = w - m.Fv.Width
	newTab.Height = h
	newTab.OffsetX = m.Fv.Width
	newTab.OffsetY = 0
	newTab.SetContent(WelcomeTabMessage)
	newTab.Redraw()
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
			m.Fv.ButtonEvent(x, y, buttons)
		case *tcell.EventResize:
			m.Fv.Width, m.Fv.Height = m.Screen.Size()
			m.Fv.RefreshTree()
			m.Fv.DrawBackground()
			m.Fv.PrintTree()
		}
	}
}

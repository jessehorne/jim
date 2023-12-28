package jim

import (
	"fmt"
	"os"
	"path"

	"github.com/rivo/tview"
)

type Editor struct {
	App     *tview.Application
	DirFile *os.File

	Pages     *tview.Pages
	PagesOpen []string

	Tree    *Tree
	Grid    *tview.Grid
	Bottom  *tview.TextView
	Bottom2 *tview.TextView

	Tabs []*Tab
}

func NewEditor(app *tview.Application, d *os.File) *Editor {
	e := &Editor{}

	treeView := NewTree(e, d.Name())
	treeView.Tree.SetBorder(true)

	bottom := tview.NewTextView()
	bottom.SetBorder(true)
	bottom.SetDynamicColors(true)
	bottom.SetTextAlign(tview.AlignCenter)
	bottom.SetText("jim v0.0.1")

	bottom2 := tview.NewTextView()
	bottom2.SetBorder(true)
	bottom2.SetDynamicColors(true)
	bottom2.SetTextAlign(tview.AlignRight)

	pages := tview.NewPages()
	pages.SetBorder(false)

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(0).
		SetBorders(false)

	grid.AddItem(treeView.Tree, 0, 0, 12, 3, 0, 0, false)
	grid.AddItem(pages, 0, 4, 12, 15, 0, 0, false)
	grid.AddItem(bottom, 12, 0, 1, 5, 0, 0, false)
	grid.AddItem(bottom2, 12, 6, 1, 14, 0, 0, false)

	e.App = app
	e.DirFile = d
	e.Pages = pages
	e.Tree = treeView
	e.Grid = grid
	e.Bottom = bottom
	e.Bottom2 = bottom2
	e.Tabs = []*Tab{}

	return e
}

func (e *Editor) OpenTab(p string) {
	t := NewTab(p)
	if t == nil {
		return
	}

	d, _ := os.ReadFile(p)

	// if tab isn't already open, we should open the tab and fill it with the content then set the line numbers
	if !e.IsFileATab(p) {
		e.Tabs = append(e.Tabs, t)

		tv := NewTextArea(path.Ext(p)).SetWordWrap(false).SetWrap(false)
		tv.SetBorder(false)
		tv.SetText(string(d), false)
		updateInfos := func() {
			fromRow, fromColumn, toRow, toColumn := tv.GetCursor()
			if fromRow == toRow && fromColumn == toColumn {
				e.Bottom2.SetText(fmt.Sprintf("Row: [yellow]%d[white], Column: [yellow]%d ", fromRow, fromColumn))
			} else {
				e.Bottom2.SetText(fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn))
			}
		}

		tv.SetMovedFunc(updateInfos)
		updateInfos()

		e.Pages.AddPage(p, tv, true, true)
		e.Pages.SwitchToPage(p)

		return
	}

	// if we get here, it means the tab should exist already. we should switch to it
	e.Pages.SwitchToPage(p)
}

func (e *Editor) SetTabActive(p string) {
	if !e.IsFileATab(p) {
		return
	}

}

func (e *Editor) IsFileATab(p string) bool {
	for _, t := range e.Tabs {
		if t.Path == p {
			return true
		}
	}

	return false
}

func (e *Editor) SaveCurrentTab() {
	// get current page
	name, page := e.Pages.GetFrontPage()
	if page == nil {
		return
	}

	tv := page.(*TextArea)
	if tv == nil {
		return
	}

	err := os.WriteFile(name, []byte(tv.GetText()), 0555)
	if err != nil {
		return
	}
}

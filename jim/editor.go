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
	TabsBox *tview.Box
	Bottom  *tview.TextView

	Tabs []*Tab
}

func NewEditor(app *tview.Application, d *os.File) *Editor {
	e := &Editor{}

	treeView := NewTree(e, d.Name())
	treeView.Tree.SetBorder(true)

	tabs := tview.NewBox()

	name := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("\njim v0.0.1")

	bottom := tview.NewTextView()
	bottom.SetBorder(true)
	bottom.SetDynamicColors(true)
	bottom.SetTextAlign(tview.AlignRight)

	pages := tview.NewPages()
	pages.SetBorder(false)

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(0).
		SetBorders(false)

	grid.AddItem(name, 0, 0, 1, 4, 0, 0, false)
	grid.AddItem(tabs, 0, 4, 1, 15, 0, 0, false)
	grid.AddItem(treeView.Tree, 1, 0, 11, 3, 0, 0, false)
	grid.AddItem(pages, 1, 4, 11, 15, 0, 0, false)
	grid.AddItem(bottom, 12, 0, 1, 19, 0, 0, false)

	e.App = app
	e.DirFile = d
	e.Pages = pages
	e.Tree = treeView
	e.Grid = grid
	e.TabsBox = tabs
	e.Bottom = bottom
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
				e.Bottom.SetText(fmt.Sprintf("Row: [yellow]%d[white], Column: [yellow]%d ", fromRow, fromColumn))
			} else {
				e.Bottom.SetText(fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn))
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

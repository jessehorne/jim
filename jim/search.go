package jim

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

type SearchReplace struct {
	Parent *Editor
	Form   tview.Primitive
}

func NewSearchReplace(e *Editor) *SearchReplace {
	var searchTerm string

	lastIndex := 0

	form := tview.NewForm()
	form.SetHorizontal(true)
	form.AddInputField("", "", 20, nil, func(text string) {
		searchTerm = text
		lastIndex = 0
	})
	//form.AddInputField("R", "", 10, nil, nil)

	form.AddButton("Find", func() {
		x, y, i := findXYOfWord(e, searchTerm, lastIndex)
		e.SetCursor(x, y)
		lastIndex = i
	})

	form.SetButtonStyle(tcell.StyleDefault.Background(tview.Styles.ContrastSecondaryTextColor))
	//form.AddButton("R", nil)
	//form.AddButton("RA", nil)

	return &SearchReplace{
		Parent: e,
		Form:   form,
	}
}

// finds the X and Y position in a document where a word starts
// begin looking at startIndex
func findXYOfWord(e *Editor, word string, startIndex int) (int, int, int) {
	// get current page
	name, p := e.Pages.GetFrontPage()
	if name == "" {
		return -1, -1, 0
	}

	// loop through content and find first occurence, saving row/col
	x, y := 0, 0
	data := p.(*TextArea).GetText()
	current := ""

	for i, c := range data {
		if i <= startIndex {
			if c == '\n' {
				y++
				x = 0
			} else if c == '\t' {
				x += 4
			} else {
				x++
			}

			continue
		}

		if len(current) > len(word) {
			current = ""
		}
		// current: c
		// word: cat
		if c == rune(word[len(current)]) {
			current += string(c)
		} else {
			current = ""
		}
		if current == word {
			return x - len(current) + 1, y, i
		}

		if c == '\n' {
			y++
			x = 0
		} else if c == '\t' {
			x += 4
		} else {
			x++
		}
	}

	return -1, -1, 0
}

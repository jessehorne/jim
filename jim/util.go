package jim

import (
	"regexp"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	// Regular expression used to escape style/region tags.
	nonEscapePattern = regexp.MustCompile(`(\[[a-zA-Z0-9_,;: \-\."#]+\[*)\]`)

	// The number of colors available in the terminal.
	availableColors = 256
)

// printWithStyle works like Print() but it takes a style instead of just a
// foreground color. The skipWidth parameter specifies the number of cells
// skipped at the beginning of the text. It returns the start index, end index
// (exclusively), and screen width of the text actually printed. If
// maintainBackground is "true", the existing screen background is not changed
// (i.e. the style's background color is ignored).
func printWithStyle(screen tcell.Screen, text string, x, y, skipWidth, maxWidth, align int, style tcell.Style, maintainBackground bool) (start, end, printedWidth int) {
	totalWidth, totalHeight := screen.Size()
	if maxWidth <= 0 || len(text) == 0 || y < 0 || y >= totalHeight {
		return 0, 0, 0
	}

	// If we don't overwrite the background, we use the default color.
	if maintainBackground {
		style = style.Background(tcell.ColorDefault)
	}

	// Skip beginning and measure width.
	var (
		state     *stepState
		textWidth int
	)
	str := text
	for len(str) > 0 {
		_, str, state = step(str, state, stepOptionsStyle)
		if skipWidth > 0 {
			skipWidth -= state.Width()
			if skipWidth <= 0 {
				text = str
				style = state.Style()
			}
			start += state.GrossLength()
		} else {
			textWidth += state.Width()
		}
	}

	// Reduce all alignments to AlignLeft.
	if align == tview.AlignRight {
		// Chop off characters on the left until it fits.
		state = nil
		for len(text) > 0 && textWidth > maxWidth {
			_, text, state = step(text, state, stepOptionsStyle)
			textWidth -= state.Width()
			start += state.GrossLength()
			style = state.Style()
		}
		x, maxWidth = x+maxWidth-textWidth, textWidth
	} else if align == tview.AlignCenter {
		// Chop off characters on the left until it fits.
		state = nil
		subtracted := (textWidth - maxWidth) / 2
		for len(text) > 0 && subtracted > 0 {
			_, text, state = step(text, state, stepOptionsStyle)
			subtracted -= state.Width()
			textWidth -= state.Width()
			start += state.GrossLength()
			style = state.Style()
		}
		if textWidth < maxWidth {
			x, maxWidth = x+maxWidth/2-textWidth/2, textWidth
		}
	}

	// Draw left-aligned text.
	end = start
	rightBorder := x + maxWidth
	state = &stepState{
		unisegState: -1,
		style:       style,
	}
	for len(text) > 0 && x < rightBorder && x < totalWidth {
		var c string
		c, text, state = step(text, state, stepOptionsStyle)
		if c == "" {
			break // We don't care about the style at the end.
		}
		width := state.Width()

		if width > 0 {
			finalStyle := state.Style()
			if maintainBackground {
				_, backgroundColor, _ := finalStyle.Decompose()
				if backgroundColor == tcell.ColorDefault {
					_, _, existingStyle, _ := screen.GetContent(x, y)
					_, background, _ := existingStyle.Decompose()
					finalStyle = finalStyle.Background(background)
				}
			}
			for offset := width - 1; offset >= 0; offset-- {
				// To avoid undesired effects, we populate all cells.
				runes := []rune(c)
				if offset == 0 {
					screen.SetContent(x+offset, y, runes[0], runes[1:], finalStyle)
				} else {
					screen.SetContent(x+offset, y, ' ', nil, finalStyle)
				}
			}
		}

		x += width
		end += state.GrossLength()
		printedWidth += width
	}

	return
}


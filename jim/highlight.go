package jim

import "github.com/gdamore/tcell/v2"

var keywords = []string{
	"break", "case", "chan", "const", "continue",
	"default", "defer", "else", "fallthrough", "for",
	"func", "go", "goto", "if", "import",
	"interface", "map", "package", "range", "return",
	"select", "struct", "switch", "type", "var",
}

var KeywordMap map[string]tcell.Style

var StyleNone = tcell.StyleDefault.Background(ColorDark).Foreground(ColorWhite)
var StyleKeyword = tcell.StyleDefault.Background(ColorDarkGrey).Foreground(ColorOrange)
var StyleLineNumber = tcell.StyleDefault.Background(ColorLightGrey).Foreground(ColorGrey)
var StyleEditor = tcell.StyleDefault.Background(ColorDarkGrey).Foreground(ColorWhite)
var StyleTab = tcell.StyleDefault.Background(ColorTabBg).Foreground(ColorWhite)
var StyleTabActive = tcell.StyleDefault.Background(ColorTabBgActive).Foreground(ColorWhite)


func InitKeywords() {
	KeywordMap = map[string]tcell.Style{}

	for _, k := range keywords {
		KeywordMap[k] = StyleKeyword
	}
}

func PrintWord(scr tcell.Screen, word string, x int, y int) {
	kw, ok := KeywordMap[word]

	if ok {
		for i := 0; i < len(word); i++ {
			scr.SetCell(x+i, y, kw, rune(word[i]))
		}
	} else {
		for i := 0; i < len(word); i++ {
			scr.SetCell(x+i, y, StyleEditor, rune(word[i]))
		}
	}
}

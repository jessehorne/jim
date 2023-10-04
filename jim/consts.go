package jim

import "github.com/gdamore/tcell/v2"

var (
	ColorDarkBlack = tcell.NewRGBColor(0, 0, 0)
	ColorBlack     = tcell.NewRGBColor(41, 45, 52)
	ColorGreen     = tcell.NewRGBColor(144, 176, 97)
	ColorOrange    = tcell.NewRGBColor(190, 138, 89)
	ColorRed       = tcell.NewRGBColor(193, 98, 102)
	ColorGrey      = tcell.NewRGBColor(157, 163, 157)
	ColorWhite     = tcell.NewRGBColor(240, 240, 240)

	ColorDark        = tcell.NewHexColor(0x06070E)
	ColorDarkGreen   = tcell.NewHexColor(0x29524A)
	ColorLightGreen  = tcell.NewHexColor(0x94A187)
	ColorLightBlue   = tcell.NewHexColor(0x3777FF)
	ColorLightGrey   = tcell.NewHexColor(0x323232)
	ColorDarkGrey    = tcell.NewHexColor(0x2b2b2b)
	ColorLighterGrey = tcell.NewHexColor(0x363636)

	ColorTabBg       = tcell.NewHexColor(0x3c3f41)
	ColorTabBgActive = tcell.NewHexColor(0x4d4a52)
)

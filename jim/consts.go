package jim

import "github.com/gdamore/tcell"

var (
	ColorBlack  = tcell.NewRGBColor(41, 45, 52)
	ColorGreen  = tcell.NewRGBColor(144, 176, 97)
	ColorOrange = tcell.NewRGBColor(190, 138, 89)
	ColorRed    = tcell.NewRGBColor(193, 98, 102)
	ColorGrey   = tcell.NewRGBColor(157, 163, 157)
	ColorWhite  = tcell.NewRGBColor(240, 240, 240)
)

const WelcomeTabMessage = `Welcome to Jim!

Jim is a simple code editor written in Go. For more information, please check out the GitHub repository...
https://github.com/jessehorne/jim

To support the project, consider making issues in GitHub for any bugs or enhancement ideas you have.
Maybe even fork it and assign an issue to yourself. :)

more coming soon...
`

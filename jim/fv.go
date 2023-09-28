package jim

import (
	"github.com/gdamore/tcell"
	"os"
	"sort"
	"strings"
)

// Fv stands for folder/file viewer. I couldn't think of a better name.
// It's the thing on the left that shows directories and files.
type Fv struct {
	Screen     tcell.Screen
	Width      int
	Height     int
	Bg         tcell.Color
	Files      []*File
	ShowHidden bool
}

func NewFv(s tcell.Screen) *Fv {
	_, h := s.Size()

	return &Fv{
		Screen:     s,
		Width:      20,
		Height:     h,
		Bg:         ColorBlack,
		Files:      []*File{},
		ShowHidden: false,
	}
}

func (fv *Fv) UnexpandDir(i int) {
	if i >= len(fv.Files) || i < 0 {
		return
	}

	f := fv.Files[i]

	if f.Type != FileTypeDir {
		return
	}

	f.Expanded = false
	f.Files = []*File{}
}

func (fv *Fv) ExpandDir(parent *File) {
	var path string
	if parent == nil {
		path = "."
	} else {
		path = parent.FullPath
	}

	var finalFiles []*File

	files, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, f := range files {
		var fileType int
		if f.IsDir() {
			fileType = FileTypeDir
		} else {
			fileType = FileTypeFile
		}

		if f.Name()[0] == '.' {
			continue
		}

		newFile := NewFile()
		newFile.Name = f.Name()
		newFile.FullPath = path + "/" + f.Name()
		newFile.Level = strings.Count(newFile.FullPath, "/")
		newFile.Type = fileType
		newFile.Expanded = false
		newFile.Parent = parent


		finalFiles = append(finalFiles, newFile)
	}

	// sort by dir / alphabetically
	sort.Slice(finalFiles, func(i, j int) bool {
		byDir := finalFiles[i].Type == FileTypeDir && finalFiles[j].Type == FileTypeDir
		byFile := finalFiles[i].Type == FileTypeDir && finalFiles[j].Type == FileTypeFile
		fileByFile := finalFiles[i].Type == FileTypeFile && finalFiles[j].Type == FileTypeFile

		if byDir {
			return strings.ToLower(finalFiles[i].FullPath) < strings.ToLower(finalFiles[j].FullPath)
		} else if byFile {
			return true
		} else if fileByFile {
			return strings.ToLower(finalFiles[i].FullPath) < strings.ToLower(finalFiles[j].FullPath)
		}

		return false
	})

	if parent == nil {
		fv.Files = finalFiles
	} else {
		parent.Files = finalFiles
		parent.Expanded = true
	}
}

// Refresh traverses the directory to update the list of files
func (fv *Fv) Refresh() {
	fv.ExpandDir(nil) // expand ./ in the main tree view
	fv.ExpandDir(fv.Files[0]) // expand ./jim in the main tree view
	fv.ExpandDir(fv.Files[0].Files[0]) // only works if jim is expanded but its for ./jim/shit
}

func (fv *Fv) PrintDir(count int, files []*File) int {
	bgStyle := tcell.StyleDefault.Background(fv.Bg).Foreground(tcell.ColorWhite)
	dirStyle := tcell.StyleDefault.Background(fv.Bg).Foreground(ColorOrange)

	offsetX := 2
	tabSize := 3

	for i := 0; i < len(files); i++ {
		f := files[i]

		fx := offsetX + (tabSize * f.Level) - tabSize

		if f.Type == FileTypeDir {
			var char rune
			if f.Expanded {
				char = 'v'
			} else {
				char = '>'
			}
			Print(fv.Screen, dirStyle, fx-1, count, string(char))
		}

		// we want to cut off the path when it gets to the wall
		currentX := fx+2
		maxX := fv.Width - currentX - 1 // max size string can be

		if maxX > 0 {
			if maxX > len(f.Name) {
				maxX = len(f.Name)
			}
			Print(fv.Screen, bgStyle, fx, count, f.Name[:maxX])
		}

		count += 1

		if f.Type == FileTypeDir && f.Expanded {
			count = fv.PrintDir(count, f.Files)
		}
	}

	return count
}

func (fv *Fv) Redraw() {
	bgStyle := tcell.StyleDefault.Background(fv.Bg).Foreground(tcell.ColorWhite)

	for y := 0; y < fv.Height; y++ {
		for x := 0; x < fv.Width; x++ {
			var char rune
			if x == fv.Width-1 {
				char = '│'
			}
			fv.Screen.SetCell(x, y, bgStyle, char)
		}
	}

	fv.PrintDir(0, fv.Files)

	fv.Screen.Sync()
}

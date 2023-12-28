package jim

import "os"

type Tab struct {
	Path string
	File *os.File
	HasChanged bool
}

func NewTab(path string) *Tab {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}

	stat, err := f.Stat()
	if err != nil {
		return nil
	}

	if stat.IsDir() {
		return nil
	}
	
	return &Tab{
		Path:       path,
		File:       f,
		HasChanged: false,
	}
}

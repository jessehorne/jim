package jim

const (
	FileTypeDir = iota
	FileTypeFile
)

// File is a directory for file represented in the file viewer
type File struct {
	Name     string
	FullPath string
	Type     int
	Level    int     // used for indentation
	Expanded bool    // if the directory is expanded
	Files    []*File //if its a directory, list its files
	Parent *File
}

func NewFile() *File {
	return &File{
		Files: []*File{},
	}
}

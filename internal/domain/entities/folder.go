package entities

type Folder struct {
	Path string
}

func NewFolder(path string) Folder {
	return Folder{
		Path: path,
	}
}

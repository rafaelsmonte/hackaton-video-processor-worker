package entities

type Folder struct {
	Path string
	Name string
}

func NewFolder(path, name string) Folder {
	return Folder{
		Path: path,
		Name: name,
	}
}

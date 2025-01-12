package entities

type File struct {
	Path string
	Name string
	Id   string
}

func NewFile(id, name, path string) File {
	return File{
		Path: path,
		Name: name,
		Id:   id,
	}
}

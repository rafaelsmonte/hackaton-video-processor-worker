package entities

type File struct {
	Path    string
	Id      string
	Content []byte
}

func NewFile(id, path string, content []byte) File {
	return File{
		Path:    path,
		Id:      id,
		Content: content,
	}
}

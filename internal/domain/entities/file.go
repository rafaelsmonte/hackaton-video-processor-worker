package entities

type File struct {
	Path    string
	UserId  string
	Id      string
	Content []byte
	Name    string
}

func NewFile(id, path, userId, name string, content []byte) File {
	return File{
		Path:    path,
		Id:      id,
		Content: content,
		UserId:  userId,
		Name:    name,
	}
}

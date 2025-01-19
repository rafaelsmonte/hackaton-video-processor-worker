package entities

type Folder struct {
	Path string
	Id   string
}

func NewFolder(path, id string) Folder {
	return Folder{
		Path: path,
		Id:   id,
	}
}

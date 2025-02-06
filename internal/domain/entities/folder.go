package entities

type Folder struct {
	Path   string
	Name   string
	Id     string
	UserId string
}

func NewFolder(path, name, id, userId string) Folder {
	return Folder{
		Path:   path,
		Name:   name,
		Id:     id,
		UserId: userId,
	}
}

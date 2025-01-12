package entities

type Video struct {
	Id   string
	Path string
	Name string
}

func NewVideo(name, id, path string) Video {
	return Video{
		Id:   id,
		Path: path,
		Name: name,
	}
}

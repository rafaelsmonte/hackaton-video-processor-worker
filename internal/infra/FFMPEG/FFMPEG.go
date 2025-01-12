package FFMPEG

import (
	"bytes"
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type FFMPEG struct {
}

// ConvertToImages implements adapters.IVideoProcessorConverter.
func (f *FFMPEG) ConvertToImages(file entities.File) (entities.Folder, error) {

	err := os.Mkdir(file.Id, os.ModeDir)
	if err != nil {
		log.Printf("Ao criar pasta destino: %v\n", err)
		return entities.NewFolder(""), err
	}
	outputPattern := filepath.Join(file.Id, "/", file.Name+"-frame-%04d.jpg")
	//TODO: remover o caminho absoluto
	cmd := exec.Command("C:\\ffmpeg\\ffmpeg.exe", "-i", file.Path+file.Name, "-vf", "fps=1", outputPattern)

	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	err = cmd.Run()
	if err != nil {
		log.Printf("Erro ao executar ffmpeg: %v\n", err)
		return entities.NewFolder(""), err
	}

	log.Println("Convertando arquivo para imagens", file)
	dir, _ := os.Getwd()
	return entities.NewFolder(dir + "\\" + file.Id), nil

}

func NewFFMPEG() adapters.IVideoProcessorConverter {
	return &FFMPEG{}
}

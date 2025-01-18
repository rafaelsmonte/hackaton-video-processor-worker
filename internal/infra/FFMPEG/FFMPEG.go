package FFMPEG

import (
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type FFMPEG struct {
}

// ConvertToImages implements adapters.IVideoProcessorConverter.
func (f *FFMPEG) ConvertToImages(file entities.File) (entities.Folder, error) {
	// destinationFolder := filepath.Join("/home/ra056172/fiap/hackaton-video-processor-worker/ffmpeg", file.Id+time.Now().Format("20060102150405"))
	// os.Mkdir(destinationFolder, os.ModeTemporary)

	// outputPattern := filepath.Join(destinationFolder, "-frame-%04d.jpg")
	// inputFile := filepath.Join(file.Path, file.Name)

	// cmd := exec.Command("ffmpeg", "-i", inputFile, "-vf", "fps=1", outputPattern)

	// cmd.Stdout = os.Stderr
	// cmd.Stderr = os.Stderr

	// err := cmd.Run()
	// if err != nil {
	// 	log.Printf("Erro ao executar ffmpeg: %v\n", err)
	// 	return entities.NewFolder(""), err

	// }
	// return entities.NewFolder(destinationFolder), nil
	fileName := strings.TrimSuffix(file.Name, filepath.Ext(file.Name))

	destinationFolder := filepath.Join("/home/ra056172/fiap/hackaton-video-processor-worker/ffmpeg", fileName+time.Now().Format("20060102150405"))
	err := os.Mkdir(destinationFolder, os.ModePerm)
	if err != nil {
		log.Printf("Ao criar pasta destino: %v\n", err)
		return entities.NewFolder(""), err
	}
	outputPattern := filepath.Join(destinationFolder, "frame-%04d.jpg")
	inputFile := filepath.Join(file.Path, file.Name)

	cmd := exec.Command("ffmpeg", "-i", inputFile, "-vf", "fps=1", outputPattern)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Printf("Erro ao executar ffmpeg: %v\n", err)
		return entities.NewFolder(""), err
	}

	log.Println("Convertando arquivo para imagens", file)
	return entities.NewFolder(destinationFolder), nil

}

func NewFFMPEG() adapters.IVideoProcessorConverter {
	return &FFMPEG{}
}

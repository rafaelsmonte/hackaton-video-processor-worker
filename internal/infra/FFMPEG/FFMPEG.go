package FFMPEG

import (
	"bytes"
	"fmt"
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FFMPEG struct {
}

func (f *FFMPEG) ConvertToImages(file entities.File) (entities.Folder, error) {
	inputFileData := file.Content
	pr, pw, _ := os.Pipe()
	folderName := strings.ReplaceAll(file.Name, ".mp4", "")

	outputDir := fmt.Sprintf("./%s/", folderName)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return entities.Folder{}, err

	}
	outputPattern := filepath.Join(outputDir, "output%d.jpg")

	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-vf", "fps=1", outputPattern)

	cmd.Stdin = pr
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting ffmpeg: %v\n", err)
		return entities.Folder{}, err
	}

	go func() {
		defer pw.Close()
		pw.Write(inputFileData)
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for ffmpeg: %v\n", err)
		stderr := cmd.Stderr.(*bytes.Buffer)
		fmt.Printf("FFMPEG error: %s\n", stderr.String()) // Imprime o erro detalhado
		fmt.Printf("Error waiting for ffmpeg: %v\n", err)
		return entities.Folder{}, err

	}

	return entities.NewFolder(outputDir, folderName, file.Id, file.UserId), nil

}

func NewFFMPEG() adapters.IVideoProcessorConverter {
	return &FFMPEG{}
}

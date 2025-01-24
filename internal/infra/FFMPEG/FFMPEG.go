package FFMPEG

import (
	"bytes"
	"fmt"
	"hackaton-video-processor-worker/internal/domain/adapters"
	"hackaton-video-processor-worker/internal/domain/entities"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type FFMPEG struct {
}

func (f *FFMPEG) ConvertToImages(file entities.File) (entities.Folder, error) {
	inputFileData := file.Content
	folderName := time.Now().Format("20060102150405")
	pr, pw, _ := os.Pipe()

	outputDir := fmt.Sprintf("./output_frames_%s/", folderName)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return entities.Folder{}, nil

	}
	outputPattern := filepath.Join(outputDir, "output%d.jpg")

	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-vf", "fps=1", outputPattern)

	cmd.Stdin = pr
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting ffmpeg: %v\n", err)
		return entities.Folder{}, nil
	}

	go func() {
		defer pw.Close()
		pw.Write(inputFileData)
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for ffmpeg: %v\n", err)
		return entities.Folder{}, nil

	}

	return entities.NewFolder(outputDir, folderName), nil

}

func NewFFMPEG() adapters.IVideoProcessorConverter {
	return &FFMPEG{}
}

package FFMPEG

import (
	"fmt"
	"hackaton-video-processor-worker/internal/domain/entities"
	"os"
	"path/filepath"
	"testing"
)

func TestConvertToImages(t *testing.T) {

	ffmpeg := NewFFMPEG()

	videoContent := []byte("fake video content")
	file := entities.File{
		Content: videoContent,
	}

	folder, err := ffmpeg.ConvertToImages(file)
	if err != nil {
		t.Fatalf("ConvertToImages failed: %v", err)
	}

	outputDir := fmt.Sprintf("./output_frames_%s/", folder.Name)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Fatalf("Output directory was not created: %v", err)
	}

	files, err := filepath.Glob(filepath.Join(outputDir, "output*.jpg"))
	if err != nil {
		t.Fatalf("Error reading output directory: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No images were generated")
	}

	// Limpa o diretório de saída após o teste
	defer os.RemoveAll(outputDir)
}

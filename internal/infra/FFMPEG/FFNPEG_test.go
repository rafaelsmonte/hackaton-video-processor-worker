package FFMPEG

import (
	"bytes"
	"fmt"
	"hackaton-video-processor-worker/internal/domain/entities"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommandExecutor é o mock da interface CommandExecutor
type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) Command(name string, args ...string) *exec.Cmd {
	ret := m.Called(name, args)
	return ret.Get(0).(*exec.Cmd)
}

// func TestConvertToImages(t *testing.T) {
// 	// Prepare mock input
// 	file := entities.File{
// 		Content: []byte("dummy data"),
// 	}

// 	// Criando mock do executor de comandos
// 	mockCmd := new(MockCommandExecutor)

// 	// Preparar ffmpegProcessor
// 	ffmpegProcessor := FFMPEG.NewFFMPEG()

// 	// Criar o mock para o comando
// 	mockExecCmd := new(exec.Cmd)

// 	// Simular o retorno do exec.Command
// 	mockCmd.On("Command", "ffmpeg", "-i", "pipe:0", "-vf", "fps=1", mock.Anything).
// 		Return(mockExecCmd)

// 	// Mockar o comportamento de Start e Wait do comando
// 	mockExecCmd.On("Start").Return(nil)
// 	mockExecCmd.On("Wait").Return(nil)

// 	// Rodar o método real
// 	folder, err := ffmpegProcessor.ConvertToImages(file)

// 	// Verificar se o método não retorna erro e cria a pasta corretamente
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, folder.Path)

// 	// Verificar se o comando foi chamado corretamente
// 	mockCmd.AssertExpectations(t)
// }

// Override the execCommand function to simulate exec.Command in tests

func TestCreateOutputDirectory(t *testing.T) {

	videoName := "fakeVideo.mp4"

	err := createFakeVideo(videoName)
	if err != nil {
	}
	defer os.Remove(videoName)

	file := entities.File{
		Content: readFile(videoName),
	}

	ffmpegProcessor := &FFMPEG{}
	outputDir := fmt.Sprintf("./output_frames_%s/", time.Now().Format("20060102150405"))

	err = os.MkdirAll(outputDir, 0755)
	assert.NoError(t, err)

	_, err = ffmpegProcessor.ConvertToImages(file)
	assert.NoError(t, err)

	defer os.RemoveAll(outputDir)
}
func TestCreateError(t *testing.T) {

	videoName := "fakeVideo.mp4"
	ffmpegProcessor := &FFMPEG{}
	outputDir := fmt.Sprintf("./output_frames_%s/", time.Now().Format("20060102150405"))
	file := entities.File{
		Content: readFile(videoName),
	}

	defer os.RemoveAll(outputDir)

	_, err := ffmpegProcessor.ConvertToImages(file)
	assert.Error(t, err)

}
func readFile(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
	return content
}
func createFakeVideo(name string) error {

	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-t", "1", "-i", "color=c=black:s=1280x720", "-c:v", "libx264", "-pix_fmt", "yuv420p", "-r", "30", name)

	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

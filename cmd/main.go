package main

import (
	"hackaton-video-processor-worker/internal/infra/api"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	log.Println("Started Env: ", os.Getenv("ENV"))
	api.SetUpServer()
}

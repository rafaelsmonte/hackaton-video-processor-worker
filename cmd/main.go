package main

import (
	httpServer "hackaton-video-processor-worker/internal/infra/http"
	"hackaton-video-processor-worker/internal/infra/sqs"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	log.Println("Started Env: ", os.Getenv("ENV"))

	go httpServer.StartHTTPServer(nil)
	sqs.SetUpSQSService()
}

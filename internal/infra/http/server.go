package httpServer

import (
	"log"
	"net/http"
	"os"
)

func StartHTTPServer() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor HTTP iniciado na porta %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

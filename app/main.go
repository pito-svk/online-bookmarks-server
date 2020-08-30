package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadConfig()

	var port string = os.Getenv("ADDRESS")

	if port == "" {
		port = ":2999"
	}

	r := chi.NewRouter()

	server := &http.Server{Addr: port, Handler: r}

	log.Printf("Server started on %s", port)
	server.ListenAndServe()
}

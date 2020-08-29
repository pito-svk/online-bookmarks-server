package main

import (
	"log"
	"net/http"
	"os"

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

	server := &http.Server{Addr: port, Handler: nil}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Server started on %s", port)
}

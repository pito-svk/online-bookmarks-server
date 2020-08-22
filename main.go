package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"peterparada.com/online-bookmarks/pkg/common/cmd"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := cmd.CreateRouter()


	var port string = os.Getenv("PORT") // || ":2999"

	server := &http.Server{Addr: port, Handler: router}

	go func () {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Online bookmarks api server is listening on %s", server.Addr)

	if err := server.Close(); err != nil {
		panic(err)
	}

	log.Println("Starting server")
}

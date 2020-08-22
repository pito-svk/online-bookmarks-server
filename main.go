package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"peterparada.com/online-bookmarks/pkg/common/cmd"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := cmd.Context()

	var port string = os.Getenv("ADDRESS")

	if port == "" {
		port = ":2999"
	}

	router := createService()

	server := &http.Server{Addr: port, Handler: router}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Online bookmarks api server is listening on %s", server.Addr)

	<-ctx.Done()

	if err := server.Close(); err != nil {
		panic(err)
	}

	time.Sleep(1000)

	log.Println("Starting server")
}

func createService() (*chi.Mux) {
	router := cmd.CreateRouter()

	return router
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_authHttpDelivery "peterparada.com/online-bookmarks/auth/delivery/http"
	_authUsecase "peterparada.com/online-bookmarks/auth/usecase"
	_pingHttpDelivery "peterparada.com/online-bookmarks/ping/delivery/http"
	_pingUsecase "peterparada.com/online-bookmarks/ping/usecase"
	_userRepo "peterparada.com/online-bookmarks/user/repository/filedb"
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

	pingUsecase := _pingUsecase.NewPingUsecase()
	_pingHttpDelivery.NewPingHandler(r, pingUsecase)

	userFileDB, err := bolt.Open("user.db", 0600, nil)
	if err != nil {
		log.Fatal("Error loding user database")
	}

	userRepo := _userRepo.NewFileDBUserRepository(userFileDB)

	authUsecase := _authUsecase.NewAuthUsecase(userRepo)
	_authHttpDelivery.NewAuthHandler(r, authUsecase)

	server := &http.Server{Addr: port, Handler: r}

	log.Printf("Server started on %s", port)
	server.ListenAndServe()
}

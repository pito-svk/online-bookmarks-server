package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/nanobox-io/golang-scribble"
	"github.com/sirupsen/logrus"
	_authHttpDelivery "peterparada.com/online-bookmarks/auth/delivery/http"
	_authUsecase "peterparada.com/online-bookmarks/auth/usecase"
	"peterparada.com/online-bookmarks/domain"
	_httpMetricsHttpDelivery "peterparada.com/online-bookmarks/httpMetrics/delivery/http"
	_httpMetricsUsecase "peterparada.com/online-bookmarks/httpMetrics/usecase"
	"peterparada.com/online-bookmarks/logging/repository"
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

func initLogger() domain.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)

	loggerImpl := repository.LoggerImpl{Logger: logger}

	return domain.Logger(&loggerImpl)
}

func getPort() string {
	var port string = os.Getenv("ADDRESS")

	if port == "" {
		port = ":2999"
	}

	return port
}

func main() {
	loadConfig()

	logger := initLogger()
	jwtSecret := os.Getenv("JWT_SECRET")

	r := chi.NewRouter()

	fileDB, err := scribble.New(".", nil)
	if err != nil {
		log.Fatal("Error loding file database")
	}

	userRepo := _userRepo.NewFileDBUserRepository(fileDB)

	httpMetricsUsecase := _httpMetricsUsecase.NewHTTPMetricsUsecase()
	pingUsecase := _pingUsecase.NewPingUsecase()
	authUsecase := _authUsecase.NewAuthUsecase(userRepo)

	_httpMetricsHttpDelivery.NewHTTPMetricsHandler(r, httpMetricsUsecase, logger)
	_pingHttpDelivery.NewPingHandler(r, pingUsecase)
	_authHttpDelivery.NewAuthHandler(r, authUsecase, logger, jwtSecret)

	port := getPort()
	server := &http.Server{Addr: port, Handler: r}

	log.Printf("Server started on %s", port)
	server.ListenAndServe()
}

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
	_httpRequestLogger "peterparada.com/online-bookmarks/common/delivery/http"
	"peterparada.com/online-bookmarks/domain"
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

type LoggerImpl struct {
	*logrus.Logger
}

// TODO: Move to separate module outside of main similar to mysql package in examples
func (logger *LoggerImpl) Trace(args ...interface{}) {
	if len(args) > 0 {
		if mapData, ok := args[0].(map[string]interface{}); ok {
			entry := logger.WithFields(logrus.Fields(mapData))

			entry.Log(logrus.TraceLevel, args[1:]...)
			return
		}
	}

	logger.Log(logrus.TraceLevel, args)
}

func initLogger() domain.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)

	loggerImpl := LoggerImpl{Logger: logger}

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

	httpRequestLoggerMiddleware := _httpRequestLogger.HttpRequestLoggerMiddleware(logger)

	r.Use(httpRequestLoggerMiddleware)

	fileDB, err := scribble.New(".", nil)
	if err != nil {
		log.Fatal("Error loding file database")
	}

	userRepo := _userRepo.NewFileDBUserRepository(fileDB)

	pingUsecase := _pingUsecase.NewPingUsecase()
	authUsecase := _authUsecase.NewAuthUsecase(userRepo)

	_pingHttpDelivery.NewPingHandler(r, pingUsecase)
	_authHttpDelivery.NewAuthHandler(r, authUsecase, logger, jwtSecret)

	port := getPort()
	server := &http.Server{Addr: port, Handler: r}

	log.Printf("Server started on %s", port)
	server.ListenAndServe()
}

package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"peterparada.com/online-bookmarks/domain"
)

type PingHandler struct {
	PingUsecase domain.PingUsecase
}

func NewPingHandler(router *chi.Mux) {
	handler := &PingHandler{}

	router.Get("/ping", handler.GetPing)
}

func (*PingHandler) GetPing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
}

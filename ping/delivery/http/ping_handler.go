package http

import (
	"github.com/go-chi/chi"
	"peterparada.com/online-bookmarks/domain"
)

type PingHandler struct {
	PingUsecase domain.PingUsecase
}

func NewPingHandler(router *chi.Mux) {
	// handler := &PingHandler{}

	// router.Get("/ping", handler.GetPing)
}
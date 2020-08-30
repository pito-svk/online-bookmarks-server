package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"peterparada.com/online-bookmarks/domain"
)

type PingHandler struct {
	PingUsecase domain.PingUsecase
}

func NewPingHandler(router *chi.Mux, us domain.PingUsecase) {
	handler := &PingHandler{
		PingUsecase: us,
	}

	router.Get("/ping", handler.GetPing)
}

func (p *PingHandler) GetPing(w http.ResponseWriter, r *http.Request) {
	pingResponse := p.PingUsecase.Get()

	w.Write([]byte(pingResponse))
}

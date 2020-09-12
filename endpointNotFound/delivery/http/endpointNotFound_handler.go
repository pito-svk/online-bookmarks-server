package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"peterparada.com/online-bookmarks/domain/entity"
)

type EndpointNotFoundHandler struct{}

func NewEndpointNotFoundHandler(router *chi.Mux) {
	handler := &EndpointNotFoundHandler{}

	router.NotFound(handler.HandleEndpointNotFound)
}

func (endpointNotFoundH *EndpointNotFoundHandler) HandleEndpointNotFound(w http.ResponseWriter, r *http.Request) {
	entity.DeliverEndpointNotFoundError(w)
}

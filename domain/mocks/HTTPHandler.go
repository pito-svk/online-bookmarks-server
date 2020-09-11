package mocks

import (
	"net/http"

	"peterparada.com/online-bookmarks/domain/entity"
)

type HTTPHandlerSettingRequestDuration struct {
	http.Handler
	WriterWithMetrics *entity.ResponseWriterWithMetrics
	Duration          int
}

func (h *HTTPHandlerSettingRequestDuration) ServeHTTP(w *entity.ResponseWriterWithMetrics, r *http.Request) {
	w.Duration = h.Duration
	h.Handler.ServeHTTP(w, r)
}

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

func (h *HTTPHandlerSettingRequestDuration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.WriterWithMetrics.Duration = h.Duration
	h.Handler.ServeHTTP(w, r)
}

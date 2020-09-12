package domain

import (
	"net/http"

	"peterparada.com/online-bookmarks/domain/entity"
)

type HTTPHandlerSettingRequestDuration interface {
	ServeHTTP(w *entity.ResponseWriterWithMetrics, r *http.Request)
}

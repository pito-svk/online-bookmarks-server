package domain

import (
	"net/http"

	"peterparada.com/online-bookmarks/domain/entity"
)

type HTTPMetricsUsecase interface {
	GetHTTPRequestMetrics(r *http.Request) (*entity.HTTPRequestMetrics, error)
	GetHTTPResponseMetrics(w *entity.ResponseWriterWithStatusCode) *entity.HTTPResponseMetrics
}

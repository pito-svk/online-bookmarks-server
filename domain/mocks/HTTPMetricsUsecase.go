package mocks

import (
	"net/http"

	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type httpMetricsUsecase struct {
}

func NewHTTPMetricsUsecase() domain.HTTPMetricsUsecase {
	return &httpMetricsUsecase{}
}

func (httpMetricsU *httpMetricsUsecase) GetHTTPRequestMetrics(r *http.Request) (*entity.HTTPRequestMetrics, error) {
	return &entity.HTTPRequestMetrics{
		URI: "/",
		Method: "GET",
		Referer: "",
		UserAgent: "",
		IP: "127.0.0.1",
	}, nil
}

func (httpMetricsU *httpMetricsUsecase) GetHTTPResponseMetrics(w *entity.ResponseWriterWithMetrics) *entity.HTTPResponseMetrics {
	return &entity.HTTPResponseMetrics{
		Code: 200,
		Duration: 50,
	}
}

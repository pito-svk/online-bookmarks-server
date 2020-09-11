package usecase

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
	ip, err := entity.GetIPAddressFromHTTPRequest(r)
	if err != nil {
		return nil, err
	}

	return &entity.HTTPRequestMetrics{
		URI:       r.URL.String(),
		Method:    r.Method,
		Referer:   r.Header.Get("Referer"),
		UserAgent: r.Header.Get("User-Agent"),
		IP:        ip,
	}, nil
}

func (httpMetricsU *httpMetricsUsecase) GetHTTPResponseMetrics(w *entity.ResponseWriterWithMetrics) *entity.HTTPResponseMetrics {
	return &entity.HTTPResponseMetrics{
		Code:     w.StatusCode,
		Duration: w.Duration,
	}
}

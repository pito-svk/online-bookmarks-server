package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type HTTPMetricsHandler struct {
	HTTPMetricsUsecase domain.HTTPMetricsUsecase
	Logger             domain.Logger
}

func NewHTTPMetricsHandler(router *chi.Mux, usecase domain.HTTPMetricsUsecase, logger domain.Logger) {
	handler := &HTTPMetricsHandler{
		HTTPMetricsUsecase: usecase,
		Logger:             logger,
	}

	router.Use(handler.LogHTTPMetrics)
}

func (httpMetricsH *HTTPMetricsHandler) LogHTTPMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMetrics, err := httpMetricsH.HTTPMetricsUsecase.GetHTTPRequestMetrics(r)
		if err != nil {
			// TODO: send error 500 or just log the error ?
		}

		// usecase.logHttpRequestData(requestData)

		writerWithMetrics := entity.NewResponseWriterWithMetrics(w)

		handlerSettingRequestDuration := entity.HTTPHandlerSettingRequestDuration{
			Handler: next,
		}

		handlerSettingRequestDuration.ServeHTTP(writerWithMetrics, r)

		responseMetrics := httpMetricsH.HTTPMetricsUsecase.GetHTTPResponseMetrics(writerWithMetrics)
		httpMetrics := map[string]interface{}{
			"uri":       requestMetrics.URI,
			"method":    requestMetrics.Method,
			"referer":   requestMetrics.Referer,
			"userAgent": requestMetrics.UserAgent,
			"ip":        requestMetrics.IP,
			"code":      responseMetrics.Code,
			"duration":  responseMetrics.Duration,
		}

		httpMetricsH.Logger.Trace(httpMetrics, "HTTP request")
	})
}

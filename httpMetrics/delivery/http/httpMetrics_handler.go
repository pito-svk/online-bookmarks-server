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

func composeHTTPMetrics(requestMetrics *entity.HTTPRequestMetrics, responseMetrics *entity.HTTPResponseMetrics) map[string]interface{} {
	return map[string]interface{}{
		"uri":       requestMetrics.URI,
		"method":    requestMetrics.Method,
		"referer":   requestMetrics.Referer,
		"userAgent": requestMetrics.UserAgent,
		"ip":        requestMetrics.IP,
		"code":      responseMetrics.Code,
		"duration":  responseMetrics.Duration,
	}
}

func logHTTPMetrics(logger domain.Logger, httpMetrics map[string]interface{}) {
	logger.Trace(httpMetrics, "HTTP request")
}

func (httpMetricsH *HTTPMetricsHandler) LogHTTPMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMetrics, err := httpMetricsH.HTTPMetricsUsecase.GetHTTPRequestMetrics(r)
		if err != nil {
			// TODO: send error 500 or just log the error ?
		}

		writerWithMetrics := entity.NewResponseWriterWithMetrics(w)
		handlerSettingRequestDuration := entity.HTTPHandlerSettingRequestDuration{Handler: next}

		handlerSettingRequestDuration.ServeHTTP(writerWithMetrics, r)

		responseMetrics := httpMetricsH.HTTPMetricsUsecase.GetHTTPResponseMetrics(writerWithMetrics)
		httpMetrics := composeHTTPMetrics(requestMetrics, responseMetrics)

		logHTTPMetrics(httpMetricsH.Logger, httpMetrics)
	})
}

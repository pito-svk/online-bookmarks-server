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

func logInternalServerError(logger domain.Logger, err error) {
	logger.Error(map[string]interface{}{
		"method":  "httpMetrics_handler/LogHTTPMetrics",
		"pointer": "internal_server_error",
	}, err)
}

func logHTTPMetrics(logger domain.Logger, httpMetrics map[string]interface{}) {
	logger.Trace(httpMetrics, "HTTP request")
}

func (httpMetricsH *HTTPMetricsHandler) getHTTPMetrics(w *entity.ResponseWriterWithMetrics, r *http.Request) (*entity.HTTPMetrics, error) {
	requestMetrics, err := httpMetricsH.HTTPMetricsUsecase.GetHTTPRequestMetrics(r)
	if err != nil {
		return nil, err
	}
	responseMetrics := httpMetricsH.HTTPMetricsUsecase.GetHTTPResponseMetrics(w)

	httpMetrics := entity.HTTPMetrics{RequestMetrics: requestMetrics, ResponseMetrics: responseMetrics}

	return &httpMetrics, nil
}

func (httpMetricsH *HTTPMetricsHandler) LogHTTPMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writerWithMetrics := entity.NewResponseWriterWithMetrics(w)
		handlerSettingRequestDuration := entity.HTTPHandlerSettingRequestDuration{Handler: next}

		handlerSettingRequestDuration.ServeHTTP(writerWithMetrics, r)

		httpMetrics, err := httpMetricsH.getHTTPMetrics(writerWithMetrics, r)
		if err != nil {
			logInternalServerError(httpMetricsH.Logger, err)
			entity.DeliverInternalServerErrorHTTPError(w)
		}

		logHTTPMetrics(httpMetricsH.Logger, httpMetrics.ToMap())
	})
}

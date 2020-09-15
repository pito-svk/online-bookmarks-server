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
	HTTPHandler        domain.HTTPHandlerSettingRequestDuration
	HTTPMetrics        *entity.HTTPMetrics
}

func NewHTTPMetricsHandler(router *chi.Mux, logger domain.Logger) {
	handler := &HTTPMetricsHandler{
		Logger: logger,
	}

	router.Use(handler.LogHTTPMetrics)
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

func (httpMetricsH *HTTPMetricsHandler) getHTTPRequestMetrics(r *http.Request) (*entity.HTTPRequestMetrics, error) {
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

func (httpMetricsU *HTTPMetricsHandler) getHTTPResponseMetrics(w *entity.ResponseWriterWithMetrics) *entity.HTTPResponseMetrics {
	return &entity.HTTPResponseMetrics{
		Code:     w.StatusCode,
		Duration: w.Duration,
	}
}

func (httpMetricsH *HTTPMetricsHandler) getHTTPMetrics(w *entity.ResponseWriterWithMetrics, r *http.Request) (*entity.HTTPMetrics, error) {
	requestMetrics, err := httpMetricsH.getHTTPRequestMetrics(r)
	if err != nil {
		return nil, err
	}
	responseMetrics := httpMetricsH.getHTTPResponseMetrics(w)

	httpMetrics := entity.HTTPMetrics{RequestMetrics: requestMetrics, ResponseMetrics: responseMetrics}

	return &httpMetrics, nil
}

func (httpMetricsH *HTTPMetricsHandler) setHTTPMetrics(httpMetrics *entity.HTTPMetrics) {
	httpMetricsH.HTTPMetrics = httpMetrics
}

func (httpMetricsH *HTTPMetricsHandler) setHTTPHandler(httpHandler domain.HTTPHandlerSettingRequestDuration) {
	httpMetricsH.HTTPHandler = httpHandler
}

func (httpMetricsH *HTTPMetricsHandler) LogHTTPMetrics(next http.Handler) http.Handler {
	if httpMetricsH.HTTPHandler == nil {
		httpMetricsH.setHTTPHandler(&entity.HTTPHandlerSettingRequestDuration{Handler: next})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseWriterWithMetrics := entity.NewResponseWriterWithMetrics(w)

		httpMetricsH.HTTPHandler.ServeHTTP(responseWriterWithMetrics, r)

		httpMetrics, err := httpMetricsH.getHTTPMetrics(responseWriterWithMetrics, r)
		if err != nil {
			logInternalServerError(httpMetricsH.Logger, err)
			entity.DeliverInternalServerErrorHTTPError(w)
		}

		httpMetricsH.setHTTPMetrics(httpMetrics)
		logHTTPMetrics(httpMetricsH.Logger, httpMetrics.ToMap())
	})
}

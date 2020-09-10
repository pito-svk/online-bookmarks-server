package entity

import (
	"net/http"
	"time"
)

type ResponseWriterWithMetrics struct {
	http.ResponseWriter
	requestTimeStart time.Time
	StatusCode       int
	Duration         int
}

func NewResponseWriterWithMetrics(w http.ResponseWriter) *ResponseWriterWithMetrics {
	return &ResponseWriterWithMetrics{
		ResponseWriter:   w,
		requestTimeStart: time.Now(),
		StatusCode:       http.StatusOK,
		Duration:         0,
	}
}

func (w *ResponseWriterWithMetrics) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

type HTTPHandlerSettingRequestDuration struct {
	http.Handler
}

func (h *HTTPHandlerSettingRequestDuration) ServeHTTP(w *ResponseWriterWithMetrics, r *http.Request) {
	h.Handler.ServeHTTP(w, r)

	// TODO: Define a method for it
	w.Duration = int(time.Now().Sub(w.requestTimeStart).Milliseconds())
}

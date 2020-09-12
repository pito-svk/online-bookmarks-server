package entity

import (
	"encoding/json"
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

func calcRequestDuration(requestStart time.Time, now time.Time) int {
	duration := calcDurationInMs(now, requestStart)

	if duration == 0 {
		duration = 1
	}

	return duration
}

func (h *HTTPHandlerSettingRequestDuration) ServeHTTP(w *ResponseWriterWithMetrics, r *http.Request) {
	h.Handler.ServeHTTP(w, r)
	w.Duration = calcRequestDuration(w.requestTimeStart, time.Now())
}

type httpErrorMessage struct {
	Error string `json:"error"`
}

func DeliverErrorParsingJSONBodyHTTPError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: "Error parsing JSON body"})
}

func DeliverBadRequestHTTPError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: err.Error()})
}

func DeliverConflictHTTPError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: err.Error()})
}

func DeliverInternalServerErrorHTTPError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: "Internal Server Error"})
}

func DeliverEndpointNotFoundError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: "Endpoint not found"})
}

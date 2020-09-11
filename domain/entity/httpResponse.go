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

func (h *HTTPHandlerSettingRequestDuration) ServeHTTP(w *ResponseWriterWithMetrics, r *http.Request) {
	h.Handler.ServeHTTP(w, r)

	duration := int(time.Now().Sub(w.requestTimeStart).Milliseconds())

	if duration == 0 {
		duration = 1
	}

	// TODO: Define a method for it
	w.Duration = duration
}

type httpErrorMessage struct {
	Error string `json:"error"`
}

func DeliverErrorParsingJSONBodyHTTPError(w http.ResponseWriter) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: "Error parsing JSON body"})
}

func DeliverBadRequestHTTPError(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: err.Error()})
}

func DeliverConflictHTTPError(w http.ResponseWriter, err error) {
	w.WriteHeader(409)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: err.Error()})
}

func DeliverInternalServerErrorHTTPError(w http.ResponseWriter) {
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(httpErrorMessage{Error: "Internal Server Error"})
}

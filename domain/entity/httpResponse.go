package entity

import "net/http"

type ResponseWriterWithStatusCode struct {
	http.ResponseWriter
	StatusCode int
}

func NewResponseWriterWithStatusCode(w http.ResponseWriter) *ResponseWriterWithStatusCode {
	return &ResponseWriterWithStatusCode{w, http.StatusOK}
}

func (w *ResponseWriterWithStatusCode) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

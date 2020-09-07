package http

import (
	"net/http"

	"peterparada.com/online-bookmarks/domain"
)

func HttpRequestLoggerMiddleware(logger domain.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			uri := r.URL.String()
			httpMethod := r.Method

			requestData := map[string]interface{}{
				"uri":    uri,
				"method": httpMethod,
			}

			logger.Trace(requestData)
		}

		return http.HandlerFunc(handlerFn)
	}
}

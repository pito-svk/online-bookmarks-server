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
			referer := r.Header.Get("Referer")
			userAgent := r.Header.Get("User-Agent")

			requestData := map[string]interface{}{
				"uri":       uri,
				"method":    httpMethod,
				"referer":   referer,
				"userAgent": userAgent,
			}

			logger.Trace(requestData, "HTTP request")

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(handlerFn)
	}
}

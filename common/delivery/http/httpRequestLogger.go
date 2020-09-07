package http

import (
	"net/http"

	"peterparada.com/online-bookmarks/domain"
)

type httpRequestData struct {
	URI        string
	HttpMethod string
}

func HttpRequestLoggerMiddleware(logger domain.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			uri := r.URL.String()
			httpMethod := r.Method

			requestData := httpRequestData{
				URI:        uri,
				HttpMethod: httpMethod,
			}

			loggingMap := map[string]interface{}{
				"uri":    requestData.URI,
				"method": requestData.HttpMethod,
			}

			logger.Trace(loggingMap)
		}

		return http.HandlerFunc(handlerFn)
	}
}

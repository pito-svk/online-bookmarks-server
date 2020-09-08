package http

import (
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"peterparada.com/online-bookmarks/domain"
)

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIP
}

func HttpRequestLoggerMiddleware(logger domain.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			uri := r.URL.String()
			httpMethod := r.Method
			referer := r.Header.Get("Referer")
			userAgent := r.Header.Get("User-Agent")
			// Add request size also here
// and test first !!!!!
			ip := requestGetRemoteAddress(r)

			httpMetrics := httpsnoop.CaptureMetrics(next, w, r)

			responseCode := httpMetrics.Code
			responseDuration := httpMetrics.Duration.Milliseconds()

			requestData := map[string]interface{}{
				"uri":       uri,
				"method":    httpMethod,
				"referer":   referer,
				"userAgent": userAgent,
				"ip":        ip,
				"code":      responseCode,
				"duration":  responseDuration,
			}

			logger.Trace(requestData, "HTTP request")
		}

		return http.HandlerFunc(handlerFn)
	}
}

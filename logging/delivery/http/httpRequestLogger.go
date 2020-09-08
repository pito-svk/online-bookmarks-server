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

func getIpAddress(r *http.Request) string {
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

type httpRequestData struct {
	URI             string
	HTTPMethod      string
	Referer         string
	UserAgent       string
	IP              string
	ResponseCode    int
	RequestDuration int
}

func getHttpRequestData(r *http.Request, httpMetrics httpsnoop.Metrics) httpRequestData {
	uri := r.URL.String()
	httpMethod := r.Method
	referer := r.Header.Get("Referer")
	userAgent := r.Header.Get("User-Agent")
	ip := getIpAddress(r)

	responseCode := httpMetrics.Code
	responseDuration := int(httpMetrics.Duration.Milliseconds())

	return httpRequestData{
		URI:             uri,
		HTTPMethod:      httpMethod,
		Referer:         referer,
		UserAgent:       userAgent,
		IP:              ip,
		ResponseCode:    responseCode,
		RequestDuration: responseDuration,
	}
}

func HttpRequestLoggerMiddleware(logger domain.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			httpMetrics := httpsnoop.CaptureMetrics(next, w, r)
			httpRequestData := getHttpRequestData(r, httpMetrics)

			requestData := map[string]interface{}{
				"uri":       httpRequestData.URI,
				"method":    httpRequestData.HTTPMethod,
				"referer":   httpRequestData.Referer,
				"userAgent": httpRequestData.UserAgent,
				"ip":        httpRequestData.IP,
				"code":      string(httpRequestData.ResponseCode),
				"duration":  string(httpRequestData.RequestDuration),
			}

			logger.Trace(requestData, "HTTP request")
		}

		return http.HandlerFunc(handlerFn)
	}
}

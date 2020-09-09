package http

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"peterparada.com/online-bookmarks/domain"
)

var privateCIDRs = []string{
	"127.0.0.0/8",
	"192.168.0.0/16",
	"172.16.0.0/12",
	"10.0.0.0/8",
	"169.254.0.0/16",
}

func isPrivateIpAddress(ip string) (bool, error) {
	ipAddress, _, err := net.ParseCIDR(fmt.Sprintf("%s/32", ip))
	if err != nil {
		return false, err
	}

	for _, privateCIDR := range privateCIDRs {
		_, privateNetwork, err := net.ParseCIDR(privateCIDR)
		if err != nil {
			return false, err
		}

		if privateNetwork.Contains(ipAddress) {
			return true, nil
		}
	}

	return false, nil
}

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func getIpAddressFromRequest(r *http.Request) string {
	hdrRealIP := r.Header.Get("X-Real-Ip")
	hdrForwardedFor := r.Header.Get("X-Forwarded-For")

	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}

	if hdrForwardedFor != "" {
		var publicParts []string

		parts := strings.Split(hdrForwardedFor, ",")
		for _, p := range parts {
			privateIp, _ := isPrivateIpAddress(strings.TrimSpace(p))

			if !privateIp {
				publicParts = append(publicParts, strings.TrimSpace(p))
			}
		}

		return publicParts[0]
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
	return httpRequestData{
		URI:             r.URL.String(),
		HTTPMethod:      r.Method,
		Referer:         r.Header.Get("Referer"),
		UserAgent:       r.Header.Get("User-Agent"),
		IP:              getIpAddressFromRequest(r),
		ResponseCode:    httpMetrics.Code,
		RequestDuration: int(httpMetrics.Duration.Milliseconds()),
	}
}

// TODO: Test middleware also
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
				"code":      httpRequestData.ResponseCode,
				"duration":  httpRequestData.RequestDuration,
			}

			logger.Trace(requestData, "HTTP request")
		}

		return http.HandlerFunc(handlerFn)
	}
}

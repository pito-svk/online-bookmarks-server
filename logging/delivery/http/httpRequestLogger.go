package http

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"peterparada.com/online-bookmarks/domain"
)

func isPrivateIpAddress(ip string) (bool, error) {
	ipNet, _, err := net.ParseCIDR(fmt.Sprintf("%s/32", ip))
	if err != nil {
		return false, err
	}

	_, cidr_192_168, err := net.ParseCIDR("192.168.0.0/16")
	if err != nil {
		return false, err
	}

	_, cidr_172_16, err := net.ParseCIDR("172.16.0.0/12")
	if err != nil {
		return false, err
	}

	_, cidr_10_0, err := net.ParseCIDR("10.0.0.0/8")
	if err != nil {
		return false, err
	}

	_, cidr_x, err := net.ParseCIDR("169.254.0.0/16")
	if err != nil {
		return false, err
	}

	isPrivateIp := ipNet.IsLoopback() || cidr_192_168.Contains(ipNet) || cidr_172_16.Contains(ipNet) || cidr_10_0.Contains(ipNet) || cidr_x.Contains(ipNet)

	return isPrivateIp, nil
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
	uri := r.URL.String()
	httpMethod := r.Method
	referer := r.Header.Get("Referer")
	userAgent := r.Header.Get("User-Agent")
	ip := getIpAddressFromRequest(r)

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

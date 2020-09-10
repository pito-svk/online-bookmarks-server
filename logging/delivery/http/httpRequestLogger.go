package http

import (
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

func isPrivateIPAddress(ipAddress string) (bool, error) {
	ip := net.ParseIP(ipAddress)

	for _, privateCIDR := range privateCIDRs {
		_, privateNetwork, err := net.ParseCIDR(privateCIDR)
		if err != nil {
			return false, err
		}

		if privateNetwork.Contains(ip) {
			return true, nil
		}
	}

	return false, nil
}

func parseIPFromXForwardedForHeader(xForwardedFor string) (string, error) {
	xForwardedForIps := strings.Split(xForwardedFor, ",")

	for _, ip := range xForwardedForIps {
		ipWithoutSpaces := strings.TrimSpace(ip)

		privateIP, err := isPrivateIPAddress(ipWithoutSpaces)
		if err != nil {
			return "", err
		}

		if !privateIP {
			return ipWithoutSpaces, nil
		}
	}

	return "", nil
}

func ipAddrFromRemoteAddr(ipAddress string) string {
	ipAddressPortIndex := strings.LastIndex(ipAddress, ":")
	ipAddressContainsPort := ipAddressPortIndex != -1

	if ipAddressContainsPort {
		ipAddressWithoutPort := ipAddress[:ipAddressPortIndex]
		return ipAddressWithoutPort
	}

	return ipAddress
}

func getIPAddressFromHttpRequest(r *http.Request) (string, error) {
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	if xForwardedFor != "" {
		ip, err := parseIPFromXForwardedForHeader(xForwardedFor)
		if err != nil {
			return "", err
		}

		if ip != "" {
			return ip, nil
		}
	}

	xRealIP := r.Header.Get("X-Real-Ip")

	if xRealIP != "" {
		return xRealIP, nil
	}

	return ipAddrFromRemoteAddr(r.RemoteAddr), nil
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

func getHttpRequestData(r *http.Request, httpMetrics httpsnoop.Metrics) (*httpRequestData, error) {
	ip, err := getIPAddressFromHttpRequest(r)
	if err != nil {
		return nil, err
	}

	return &httpRequestData{
		URI:             r.URL.String(),
		HTTPMethod:      r.Method,
		Referer:         r.Header.Get("Referer"),
		UserAgent:       r.Header.Get("User-Agent"),
		IP:              ip,
		ResponseCode:    httpMetrics.Code,
		RequestDuration: int(httpMetrics.Duration.Milliseconds()),
	}, nil
}

// TODO: Test middleware also
func HttpRequestLoggerMiddleware(logger domain.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handlerFn := func(w http.ResponseWriter, r *http.Request) {
			httpMetrics := httpsnoop.CaptureMetrics(next, w, r)
			httpRequestData, err := getHttpRequestData(r, httpMetrics)
			if err != nil {
				// TODO: Implement - return internal server error probably
			}

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

package usecase

import (
	"net"
	"net/http"
	"strings"

	"peterparada.com/online-bookmarks/domain"
	"peterparada.com/online-bookmarks/domain/entity"
)

type httpMetricsUsecase struct {
}

func NewHTTPMetricsUsecase() domain.HTTPMetricsUsecase {
	return &httpMetricsUsecase{}
}

func (httpMetricsU *httpMetricsUsecase) GetHTTPRequestMetrics(r *http.Request) (*entity.HTTPRequestMetrics, error) {
	ip, err := getIPAddressFromHTTPRequest(r)
	if err != nil {
		return nil, err
	}

	return &entity.HTTPRequestMetrics{
		URI:       r.URL.String(),
		Method:    r.Method,
		Referer:   r.Header.Get("Referer"),
		UserAgent: r.Header.Get("User-Agent"),
		IP:        ip,
	}, nil
}

func (httpMetricsU *httpMetricsUsecase) GetHTTPResponseMetrics(w *entity.ResponseWriterWithStatusCode) *entity.HTTPResponseMetrics {
	return &entity.HTTPResponseMetrics{
		Code: w.StatusCode,
	}
}

var privateCIDRs = []string{
	"127.0.0.0/8",
	"192.168.0.0/16",
	"172.16.0.0/12",
	"10.0.0.0/8",
	"169.254.0.0/16",
}

func getIPAddressFromHTTPRequest(r *http.Request) (string, error) {
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
		ip := parseIPFromXRealIPHeader(xRealIP)
		return ip, nil
	}

	return parseRemoteIPAddress(r.RemoteAddr), nil
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

func parseIPFromXRealIPHeader(xRealIP string) string {
	return xRealIP
}

func parseRemoteIPAddress(ipAddress string) string {
	ipAddressPortIndex := strings.LastIndex(ipAddress, ":")
	ipAddressContainsPort := ipAddressPortIndex != -1

	if ipAddressContainsPort {
		ipAddressWithoutPort := ipAddress[:ipAddressPortIndex]
		return ipAddressWithoutPort
	}

	return ipAddress
}

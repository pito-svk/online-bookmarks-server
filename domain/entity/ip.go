package entity

import (
	"net"
	"net/http"
	"strings"
)

// TODO: Shouldn't be part of entity
type IPAddress struct {
	Address string
}

var privateCIDRs = []string{
	"127.0.0.0/8",
	"192.168.0.0/16",
	"172.16.0.0/12",
	"10.0.0.0/8",
	"169.254.0.0/16",
}

func (ipAddress *IPAddress) IsPrivate() (bool, error) {
	ip := net.ParseIP(ipAddress.Address)

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

	for _, ipFromHeader := range xForwardedForIps {
		ipWithoutSpaces := strings.TrimSpace(ipFromHeader)

		ip := IPAddress{Address: ipWithoutSpaces}

		privateIP, err := ip.IsPrivate()
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

func GetIPAddressFromHTTPRequest(r *http.Request) (string, error) {
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

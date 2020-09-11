package entity

import "net"

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

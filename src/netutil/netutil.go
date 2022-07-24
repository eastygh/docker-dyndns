package netutil

import "net"

func ValidateIpV4(ipV4 string) bool {
	v4addr := net.ParseIP(ipV4)
	if v4addr == nil {
		return false
	}
	return v4addr.To4() != nil
}

func ValidateIpV6(ipV6 string) bool {
	v6addr := net.ParseIP(ipV6)
	if v6addr == nil {
		return false
	}
	return v6addr.To16() != nil
}

func IsDomainValid(domain string, domains []string) bool {
	for _, cur := range domains {
		if cur == domain {
			return true
		}
	}
	return false
}

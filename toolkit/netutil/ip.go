package netutil

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	envPodIp = "POD_IP"
)

func PodIP() string {
	return os.Getenv(envPodIp)
}

func InternalIP() (string, error) {
	ip, err := interfaceIP(func(ip net.IP) bool {
		return !ip.IsLoopback()
	})
	if err != nil {
		return "", err
	}
	return ip.String(), nil
}

func InternalIPV4() (string, error) {
	ip, err := interfaceIP(func(ip net.IP) bool {
		return !ip.IsLoopback() && ip.To4() != nil
	})
	if err != nil {
		return "", err
	}
	return ip.String(), nil
}

func interfaceIP(fn func(ip net.IP) bool) (net.IP, error) {
	infs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, inf := range infs {
		if inf.Flags&net.FlagUp != net.FlagUp || inf.Flags&net.FlagLoopback == net.FlagLoopback {
			continue
		}
		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && fn(ipnet.IP) {
				return ipnet.IP, nil
			}
		}
	}
	return nil, errors.New("no ip found")
}

func FingerOut(address string, v4only ...bool) (string, int) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
	if !addr.IP.IsUnspecified() && !addr.IP.IsLoopback() {
		return addr.IP.String(), addr.Port
	}
	if ip := PodIP(); ip != "" {
		return ip, addr.Port
	}

	ip, err := interfaceIP(func(ip net.IP) bool {
		if len(v4only) > 0 && v4only[0] {
			return !ip.IsLoopback() && ip.To4() != nil
		}
		return !ip.IsLoopback()
	})
	if err != nil {
		return addr.IP.String(), addr.Port
	}
	return ip.String(), addr.Port
}

// RemoteIP parses the IP from Request.RemoteAddr, normalizes and returns the IP (without the port).
func RemoteIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return ""
	}
	return ip
}

// ClientIP
// see also: IPValidator.ClientIP
func ClientIP(r *http.Request) string {
	return defaultIPValidator.ClientIP(r)
}

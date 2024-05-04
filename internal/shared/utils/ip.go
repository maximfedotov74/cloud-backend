package utils

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

func GetIP(r *http.Request) (net.IP, error) {
	netIp := net.ParseIP(r.Header.Get("X-Real-IP"))

	if netIp != nil {
		return netIp, nil
	}

	forwardFor := r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(forwardFor, ",") {
		netIp := net.ParseIP(i)

		if netIp != nil {
			return netIp, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}

	netIp = net.ParseIP(ip)

	if netIp != nil {
		return netIp, nil
	}

	netIp = net.ParseIP(ip)

	if netIp != nil {
		return netIp, nil
	}

	return nil, errors.New("no valid ip found")
}

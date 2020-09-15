package utils

import (
	"context"
	"go.elastic.co/apm"
	"net"
	"os"
	"runtime"
)

func CopyContext(ctx context.Context) context.Context {
	return apm.DetachedContext(ctx)
}

func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}

	return hostname
}

func LocalIP() string {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addr {
		// check the address type and if it is not a loopback the display it
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}
	return ""
}

func OSInfo() (string, string) {
	return runtime.GOOS, runtime.GOARCH
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

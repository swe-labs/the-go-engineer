package main

import (
	"fmt"
	"strconv"
	"strings"
)

func SplitHostPort(addr string) (host string, port int, err error) {
	host, portStr, cut := strings.Cut(addr, ":")
	if !cut {
		return "", 0, fmt.Errorf("missing port in address %q", addr)
	}
	port, err = strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port %q in address %q", portStr, addr)
	}
	if port < 1 || port > 65535 {
		return "", 0, fmt.Errorf("port %d out of range in address %q", port, addr)
	}
	return host, port, nil
}

func main() {
	fmt.Println(SplitHostPort("localhost:8080"))
}

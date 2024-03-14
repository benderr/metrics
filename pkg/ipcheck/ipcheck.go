// Package trustedip contains functionality for creating a middleware for check IP by trusted subnet (CIDR)
package ipcheck

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	_ "github.com/go-chi/chi/middleware"
)

// Middleware return function for check ip by trusted subnet rules
func Middleware(trustedSubnet string) func(next http.Handler) http.Handler {
	if len(trustedSubnet) == 0 {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	_, trustedIPs, err := net.ParseCIDR(trustedSubnet)

	if err != nil {
		panic(errors.New("invalid "))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			realIP := net.ParseIP(r.Header.Get("X-Real-IP"))

			fmt.Println("REAL IP", realIP, r.Header.Get("X-Real-IP"))
			if realIP == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if !trustedIPs.Contains(realIP) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

var ErrNotFoundIP = errors.New("host ip not found")

func GetHostIP() (net.IP, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return net.IP{}, err
	}

	for _, addr := range addr {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}
	return net.IP{}, ErrNotFoundIP
}

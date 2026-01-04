package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/hamza-s47/api-guard/internal/store"
)

func getClientIP(r *http.Request) string {
	// 1. X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, "")[0]
	}
	// 2. X-Real-IP
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		return rip
	}
	// 3. RemoteAddr fallback
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func RateLimit(store *store.MemoryStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			bucket := store.GetBucket(ip)
			w.Header().Set("X-RateLimit-Limit", "5")
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(bucket.Remaining()))
			if !bucket.Allow() {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

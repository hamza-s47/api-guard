package middleware

import (
	"net"
	"net/http"

	"github.com/hamza-s47/api-guard/internal/store"
)

func RateLimit(store *store.MemoryStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Invalid IP", http.StatusInternalServerError)
				return
			}

			bucket := store.GetBucket(ip)
			if !bucket.Allow() {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

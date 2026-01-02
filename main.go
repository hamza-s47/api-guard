package main

import (
	"log"
	"net/http"

	"github.com/hamza-s47/api-guard/handler"
	"github.com/hamza-s47/api-guard/internal/proxy"
	"github.com/hamza-s47/api-guard/internal/store"
	"github.com/hamza-s47/api-guard/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)

	// Reverse Proxy
	backendProxy, err := proxy.NewReverseProxy("http://localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/api/", backendProxy)

	store := store.NewMemoryStore()
	handlerWithMiddleware := middleware.Logging(middleware.RateLimit(store)(mux))

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithMiddleware))
}

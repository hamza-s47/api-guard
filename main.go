package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store := store.NewMemoryStore()
	finalHandler := middleware.Logging(middleware.RateLimit(store)(mux))
	server := &http.Server{Addr: ":8080", Handler: finalHandler}

	go func() {
		log.Println("API Gateway running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("\nShutting down...")

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(shutDownCtx)
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/hamza-s47/api-guard/handler"
	"github.com/hamza-s47/api-guard/internal/config"
	"github.com/hamza-s47/api-guard/internal/proxy"
	"github.com/hamza-s47/api-guard/internal/store"
	"github.com/hamza-s47/api-guard/middleware"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	config.Watch("config.yaml")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)

	// Reverse Proxy
	for _, route := range cfg.Routes {
		proxy, _ := proxy.NewReverseProxy(route.Backend)
		mux.Handle(route.Path, proxy)
	}
	// backendProxy, err := proxy.NewReverseProxy("http://localhost:9000")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// mux.Handle("/api/", backendProxy)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store := store.NewMemoryStore()
	finalHandler := middleware.Logging(middleware.RateLimit(store)(mux))
	server := &http.Server{Addr: ":" + strconv.Itoa(cfg.Server.Port), Handler: finalHandler}

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

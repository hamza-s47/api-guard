package main

import (
	"log"
	"net/http"

	"github.com/hamza-s47/api-guard/handler"
	"github.com/hamza-s47/api-guard/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	handlerWithMiddleware := middleware.Logging(mux)

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithMiddleware))
}

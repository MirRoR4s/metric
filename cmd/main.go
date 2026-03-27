package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	metric "github.com/MirRoR4s/metric/pkg"
)

func main() {
	ctx := context.Background()
	registry := metric.NewRegistry()
	counter, middleware := metric.HttpRequestsTotal()
	memoryMetric := metric.Memory(ctx)
	registry.Register(counter, memoryMetric)
	mux := http.NewServeMux()
	mux.Handle("/metrics", registry.Handler())
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, World!"))
	})
	srv := &http.Server{
		Addr:    ":8080",
		Handler: middleware(mux),
	}
	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Server closed")
		} else {
			log.Fatalf("Server error: %v", err)
		}
	}
}

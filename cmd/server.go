package main

import (
	"log"
	"net/http"
	"time"

	"github.com/victor-tsykanov/delivery/internal/common/config"
)

func main() {
	config.MustLoadEnv(".env")

	httpConfig, err := config.LoadHTTPConfig()
	if err != nil {
		log.Fatalf("failed load HTTP config: %v", err)
	}

	server := &http.Server{
		Addr:              httpConfig.Address(),
		ReadHeaderTimeout: 3 * time.Second,
	}

	http.HandleFunc("/health", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

package main

import (
	"github.com/PaytonWebber/go-webhook-listener/config"
	"github.com/PaytonWebber/go-webhook-listener/internal/handlers"
	"log"
	"net/http"
	"strconv"
)

func main() {
	cfg := config.LoadConfig()
	handler := handlers.NewWebhookHandler(&cfg)

	http.Handle("/restart", handler)

	portStr := strconv.Itoa(cfg.Port)
	addr := ":" + portStr

	log.Printf("Listening on port %s", portStr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

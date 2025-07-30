package main

import (
	"fmt"
	"log"
	"net/http"

	"proxy-stream/client"
	"proxy-stream/config"
	"proxy-stream/stream"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Error config: %v", err)
	}

	client := client.New(cfg.GetCookieSecure())
	stream := stream.New(cfg.GetStreamUrl(), false)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		client.CreateCookieUniqueClient(w, r)
		stream.StreamHandler(w, r)
	})

	address := fmt.Sprintf("0.0.0.0:%d", cfg.GetPort())
	log.Printf("Listening on http://localhost:%d", cfg.GetPort())

	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Eroare server:", err)
	}
}

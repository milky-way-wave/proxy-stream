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
		log.Fatal(err)
	}

	client := client.New(cfg.COOKIE_SECURE)
	stream := stream.New(cfg.STREAM_URL, false)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		client.CreateCookieUniqueClient(w, r)
		stream.StreamHandler(w, r)
	})

	address := fmt.Sprintf("0.0.0.0:%d", cfg.PORT)
	log.Printf("Listening on http://localhost:%d", cfg.PORT)

	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Eroare server:", err)
	}
}

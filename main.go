package main

import (
	"fmt"
	"log"
	"net/http"

	"proxy-stream/client"
	"proxy-stream/config"
	"proxy-stream/cookie"
	"proxy-stream/hmac"
	"proxy-stream/stream"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Error config: %v", err)
	}

	hmac := hmac.New(cfg.GetAppSecret())
	cookie := cookie.New(cfg.GetCookieSecure(), hmac)
	stream := stream.New(cfg.GetStreamUrl(), true)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookieValue, err := cookie.Has(r)
		if err != nil {
			cookie.Delete(w)
			cookieValue = cookie.Create(w, r)
		}

		client := client.New(cookieValue)
		stream.StreamHandler(w, r, client)
	})

	address := fmt.Sprintf("0.0.0.0:%d", cfg.GetPort())
	log.Printf("Listening on http://localhost:%d", cfg.GetPort())

	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}

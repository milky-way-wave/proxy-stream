package client

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net"
	"net/http"
)

type client struct {
	cookies_keu_name string
	secure           bool
}

func New(secure bool) client {
	return client{
		cookies_keu_name: "c_uid",
		secure:           secure,
	}
}

func (c *client) genUid(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	ip := r.Header.Get("X-Real-IP")

	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	raw := userAgent + ip
	hash := md5.Sum([]byte(raw))

	return hex.EncodeToString(hash[:])
}

func (c *client) CreateCookieUniqueClient(w http.ResponseWriter, r *http.Request) {
	clientID, err := r.Cookie(c.cookies_keu_name)
	if err != nil {
		clientID := c.genUid(r)

		http.SetCookie(w, &http.Cookie{
			Name:     c.cookies_keu_name,
			Value:    clientID,
			Path:     "/",
			HttpOnly: true,
			Secure:   c.secure,
		})
	}

	userAgent := r.Header.Get("User-Agent")
	ip := r.Header.Get("X-Real-IP")

	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	log.Printf("[%s][%s] => %s", clientID, ip, userAgent)

}

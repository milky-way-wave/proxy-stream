package cookie

import (
	"crypto/md5"
	"fmt"
	"net"
	"net/http"
	"strings"

	"proxy-stream/helper/generator"
	"proxy-stream/hmac"
)

const (
	COOKIE_MAX_AGE_100_YEARS = 100 * 365 * 24 * 3600 // 3 years in seconds
)

type Cookie struct {
	name   string
	secure bool
	hmac   hmac.Hmac
}

func New(secureCookie bool, hmac hmac.Hmac) Cookie {
	return Cookie{
		name:   "cid",
		secure: secureCookie,
		hmac:   hmac,
	}
}

func (c *Cookie) Create(w http.ResponseWriter, r *http.Request) string {

	value := c.prepare(r)

	http.SetCookie(w, &http.Cookie{
		Name:   c.name,
		Value:  value,
		Path:   "/",
		Secure: c.secure,
		MaxAge: COOKIE_MAX_AGE_100_YEARS,
	})

	return value
}

func (c *Cookie) Delete(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   c.name,
		Value:  "",
		Path:   "/",
		Secure: c.secure,
		MaxAge: -1,
	})
}

func (c *Cookie) Has(r *http.Request) (string, error) {
	httpCookie, err := r.Cookie(c.name)

	if err != nil {
		return "", err
	}

	cookieValue := httpCookie.Value
	key := strings.Split(cookieValue, ":")

	if len(key) != 2 {
		return "", fmt.Errorf("not found cookies[cid]")
	}

	if !c.hmac.Verify(key[0], key[1]) {
		return "", fmt.Errorf("invalid cookie")
	}

	return cookieValue, nil
}

func (c *Cookie) prepare(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	ip := r.Header.Get("X-Real-IP")

	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	uuid := generator.UUIDv7()
	hash := fmt.Sprintf("%X", md5.Sum([]byte(fmt.Sprint("%s:%s:%s", userAgent, ip, uuid))))

	return fmt.Sprintf("%s:%s", hash, c.hmac.Sign(hash))
}

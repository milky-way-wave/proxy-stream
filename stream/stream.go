package stream

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type stream struct {
	url      string
	showMeta bool
}

func New(url string, showMeta bool) stream {
	return stream{
		url:      url,
		showMeta: showMeta,
	}
}

func (s *stream) isBrokenPipeErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.EOF) {
		return false
	}
	if opErr, ok := err.(*net.OpError); ok && opErr.Err != nil {
		if strings.Contains(opErr.Err.Error(), "broken pipe") {
			return true
		}
	}
	return strings.Contains(err.Error(), "write: broken pipe")
}

func (s *stream) extractStreamTitle(meta string) string {
	prefix := "StreamTitle='"
	start := strings.Index(meta, prefix)
	if start == -1 {
		return ""
	}
	start += len(prefix)
	end := strings.Index(meta[start:], "';")
	if end == -1 {
		return ""
	}
	return meta[start : start+end]
}

func (s *stream) StreamHandler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", s.url, nil)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Icy-MetaData", "1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Nu pot accesa streamul Icecast", http.StatusBadGateway)
		log.Println("Eroare la conectare cu upstream:", err)
		return
	}
	defer resp.Body.Close()

	metaintStr := resp.Header.Get("Icy-Metaint")
	metaint, err := strconv.Atoi(metaintStr)

	if err != nil || metaint == 0 {
		log.Println("Upstream nu trimite metadata, fac proxy simplu")
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		_, err := io.Copy(w, resp.Body)
		if err != nil && !s.isBrokenPipeErr(err) {
			log.Printf("Eroare la copiere stream: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusOK)

	bufAudio := make([]byte, metaint)

	for {
		_, err := io.ReadFull(resp.Body, bufAudio)
		if err != nil {
			if err == io.EOF {
				return
			}
			if !s.isBrokenPipeErr(err) {
				log.Printf("Eroare citire audio upstream: %v", err)
			}
			return
		}

		_, err = w.Write(bufAudio)
		if err != nil {
			if !s.isBrokenPipeErr(err) {
				log.Printf("Eroare scriere audio client: %v", err)
			}
			return
		}

		lenByte := make([]byte, 1)
		_, err = io.ReadFull(resp.Body, lenByte)
		if err != nil {
			if err == io.EOF {
				return
			}
			if !s.isBrokenPipeErr(err) {
				log.Printf("Eroare citire lungime metadata upstream: %v", err)
			}
			return
		}

		metaLen := int(lenByte[0]) * 16
		if metaLen == 0 {
			continue
		}

		metaBuf := make([]byte, metaLen)
		_, err = io.ReadFull(resp.Body, metaBuf)
		if err != nil {
			if err == io.EOF {
				return
			}
			if !s.isBrokenPipeErr(err) {
				log.Printf("Eroare citire metadata upstream: %v", err)
			}
			return
		}

		if s.showMeta {
			metaStr := string(metaBuf)
			streamTitle := s.extractStreamTitle(metaStr)
			if streamTitle != "" {
				log.Println("🎵 ", streamTitle)
			}
		}

		_, err = w.Write(metaBuf)
		if err != nil {
			if !s.isBrokenPipeErr(err) {
				log.Printf("Eroare scriere metadata client: %v", err)
			}
			return
		}
	}
}

package stream

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"proxy-stream/client"
	"strconv"
	"strings"
)

type UserCounter interface {
	AddUser(id string)
	RemoveUser(id string)
}

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

func (s *stream) handleError(err error, client client.Client, context string) bool {
	if err == nil {
		return false
	}
	if s.isBrokenPipeErr(err) {
		return true
	}
	log.Printf("❌ Error %s: %v", context, err)
	return true
}

func (s *stream) extractStreamUrl(meta string) string {
	prefix := "StreamUrl='"
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

func (s *stream) StreamHandler(w http.ResponseWriter, r *http.Request, client client.Client) {
	client.Connect()

	defer func() {
		client.Disconnect()
	}()

	req, err := http.NewRequest("GET", s.url, nil)
	if err != nil {
		log.Printf("❌ Failed to create request: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Icy-MetaData", "1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("❌ Cannot connect to upstream %s: %v", s.url, err)
		http.Error(w, "Cannot access Icecast stream", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	metaintStr := resp.Header.Get("Icy-Metaint")
	metaint, err := strconv.Atoi(metaintStr)
	if err != nil || metaint == 0 {
		log.Println("ℹ️  Upstream doesn't send metadata, doing simple proxy")
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		_, err := io.Copy(w, resp.Body)
		if s.handleError(err, client, "copying stream") {
			return
		}
		return
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(http.StatusOK)

	bufAudio := make([]byte, metaint)
	for {
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		_, err := io.ReadFull(resp.Body, bufAudio)
		if err != nil {
			if err == io.EOF {
				log.Println("📡 Upstream stream ended")
				return
			}
			if s.handleError(err, client, "reading audio from upstream") {
				return
			}
		}

		_, err = w.Write(bufAudio)
		if s.handleError(err, client, "writing audio to client") {
			return
		}

		lenByte := make([]byte, 1)
		_, err = io.ReadFull(resp.Body, lenByte)
		if err != nil {
			if err == io.EOF {
				log.Println("📡 Upstream metadata stream ended")
				return
			}
			if s.handleError(err, client, "reading metadata length from upstream") {
				return
			}
		}

		metaLen := int(lenByte[0]) * 16
		if metaLen == 0 {
			continue
		}

		metaBuf := make([]byte, metaLen)
		_, err = io.ReadFull(resp.Body, metaBuf)
		if err != nil {
			if err == io.EOF {
				log.Println("📡 Upstream metadata ended")
				return
			}
			if s.handleError(err, client, "reading metadata from upstream") {
				return
			}
		}

		metaStr := string(metaBuf)
		streamTitle := s.extractStreamTitle(metaStr)
		streamUrl := s.extractStreamUrl(metaStr)

		if s.showMeta {
			logMetadata := ""

			if streamTitle != "" {
				logMetadata += fmt.Sprintf("🎵 %s", streamTitle)
			}
			if streamUrl != "" {
				logMetadata += fmt.Sprintf(" <=> 🔗 %s", streamUrl)
			}

			log.Printf("metadata: %q", logMetadata)
		}

		_, err = w.Write(metaBuf)
		if s.handleError(err, client, "writing metadata to client") {
			return
		}
	}
}

package cacheproxy

import (
	"crypto/sha256"
	"io"
	"net/http"
)

func fileMIME(respBody []byte, header http.Header) string {
	if contentType := header.Get("Content-Type"); contentType != "" &&
		contentType != "application/octet-stream" {
		return contentType
	}
	return http.DetectContentType(respBody)
}

func checksum(body []byte) []byte {
	h := sha256.Sum256(body)
	return h[:]
}

// bodyReader decodes the request body based on the Accept-Encoding header.
func bodyReader(respBody io.Reader) ([]byte, error) {
	return io.ReadAll(respBody)
}

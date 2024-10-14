package cacheproxy

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

func (proxy *CacheableProxy) cacheKey(req *http.Request) string {
	query, err := url.QueryUnescape(req.URL.RawQuery)
	if err != nil {
		query = req.URL.RawQuery
	}
	cacheKey := fmt.Sprintf(
		"file://%s@%s#%s",
		req.Method, proxy.targetURL.String(),
		req.URL.Path+query,
	)
	return strings.TrimSpace(cacheKey)
}

func (proxy *CacheableProxy) InterceptFile(resp *http.Response) error {
	// Get the requested file URL from the request
	fileURL := resp.Request.RequestURI
	cacheKey := proxy.cacheKey(resp.Request)
	now := time.Now()

	cachedFile, err := proxy.storage.Get(cacheKey)
	if err == nil && len(cachedFile.Checksum) > 0 &&
		cachedFile.ModifiedAt.Sub(now) < proxy.cacheTTL {
		return nil
	}

	// Not cached, make a request to the target site and store the result in the cache
	const size = 11 << 7
	respBody := make([]byte, 0, size)
	if respBody, err = io.ReadAll(resp.Body); err != nil {
		return err
	}

	fileInfo := FileInformation{
		FileMIME: FileMIME{
			Name:      fileURL,
			Extension: filepath.Ext(fileURL),
			MimeType:  fileMIME(respBody, resp.Header),
		},
		Envelope: FileEnvelope{
			Headers: resp.Header,
			Status:  uint16(resp.StatusCode),
		},
		Content:       respBody,
		Checksum:      checksum(respBody),
		CreatedAt:     now,
		ModifiedAt:    now,
		ExtraMetadata: make(map[string]string),
	}

	// Reassign the body so that it can be sent to the client
	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	if proxy.isFileTracked(fileInfo) {
		return proxy.storage.Set(cacheKey, fileInfo)
	}

	return nil
}

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

func (proxy *CacheableProxy) isFileTracked(info FileInformation) bool {
	for _, extension := range proxy.trackedExtensions {
		mimeList := strings.Split(info.MimeType, ";")
		if strings.EqualFold(info.Extension, extension) {
			return true
		}
		if slices.ContainsFunc(
			mimeList, func(s string) bool { return strings.EqualFold(s, extension) },
		) {
			return true
		}
	}
	return false
}

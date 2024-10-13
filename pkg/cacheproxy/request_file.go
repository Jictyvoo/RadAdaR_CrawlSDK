package cacheproxy

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func (proxy CacheableProxy) cacheKey(req *http.Request) string {
	cacheKey := fmt.Sprintf("file://%s@%s", req.Method, req.URL.String())
	return cacheKey
}

func (proxy CacheableProxy) InterceptFile(resp *http.Response) error {
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
	var respBody []byte
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fileInfo := FileInformation{
		FileMIME: FileMIME{
			Name:      fileURL,
			Extension: filepath.Ext(fileURL),
			MimeType:  http.DetectContentType(respBody),
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

	if proxy.isFileTracked(fileInfo) {
		return proxy.storage.Set(cacheKey, fileInfo)
	}

	return nil
}

func checksum(body []byte) []byte {
	h := sha256.Sum256(body)
	return h[:]
}

func (proxy CacheableProxy) isFileTracked(info FileInformation) bool {
	for _, extension := range proxy.trackedExtensions {
		if strings.EqualFold(info.Extension, extension) {
			return true
		}
	}
	return false
}

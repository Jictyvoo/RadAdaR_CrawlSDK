package cacheproxy

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

type (
	CacheStorage   = KVStorage[FileInformation, string]
	CacheableProxy struct {
		storage           CacheStorage
		cacheTTL          time.Duration
		targetURL         *url.URL
		port              uint16
		trackedExtensions []string
		proxy             *httputil.ReverseProxy
	}
)

func New(storage CacheStorage, targetURL string, port uint16) (*CacheableProxy, error) {
	target, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return nil, err
	}

	cacheableProxy := &CacheableProxy{
		storage:   storage,
		targetURL: target,
		port:      port,
		cacheTTL:  36 * time.Hour,
		proxy:     httputil.NewSingleHostReverseProxy(target),
	}
	cacheableProxy.proxy.ModifyResponse = cacheableProxy.InterceptFile
	return cacheableProxy, nil
}

func (proxy CacheableProxy) Handler(w http.ResponseWriter, r *http.Request) {
	slog.Info(
		"[ PROXY SERVER ] Request received",
		slog.String("URL", r.URL.String()), slog.Time("time", time.Now()),
	)

	cacheKey := proxy.cacheKey(r)
	fileInfo, err := proxy.storage.Get(cacheKey)
	if err != nil {
		// Finally return proxy
		proxy.proxy.ServeHTTP(w, r)
		return
	}

	// Restore file response
	_, err = w.Write(fileInfo.Content)
	w.WriteHeader(int(fileInfo.Envelope.Status))
	if err != nil {
		slog.Error("[ PROXY SERVER ] Error writing response", slog.String("error", err.Error()))
	}
}

func (proxy CacheableProxy) Listen(ctx context.Context) error {
	http.HandleFunc("/", proxy.Handler)

	return http.ListenAndServe(":"+strconv.FormatUint(uint64(proxy.port), 10), nil)
}

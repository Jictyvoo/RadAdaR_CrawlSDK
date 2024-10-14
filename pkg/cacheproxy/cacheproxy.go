package cacheproxy

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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
		reverse           *httputil.ReverseProxy
	}
)

func New(storage CacheStorage, targetURL string, port uint16) (*CacheableProxy, error) {
	target, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return nil, err
	}

	cacheableProxy := &CacheableProxy{
		storage:           storage,
		targetURL:         target,
		port:              port,
		cacheTTL:          36 * time.Hour,
		reverse:           httputil.NewSingleHostReverseProxy(target),
		trackedExtensions: []string{"text/html", "image/jpeg"},
	}
	cacheableProxy.reverse.ModifyResponse = cacheableProxy.InterceptFile
	cacheableProxy.reverse.Director = cacheableProxy.Director
	return cacheableProxy, nil
}

func (proxy *CacheableProxy) Handler(w http.ResponseWriter, r *http.Request) {
	slog.Info(
		"[ PROXY SERVER ] Request received",
		slog.String("URL", r.URL.String()), slog.Time("time", time.Now()),
	)

	cacheKey := proxy.cacheKey(r)
	fileInfo, err := proxy.storage.Get(cacheKey)
	if err != nil || len(fileInfo.Checksum) <= 0 {
		// Finally return reverse
		proxy.reverse.ServeHTTP(w, r)
		return
	}

	// Restore file response
	for key, values := range fileInfo.Envelope.Headers {
		w.Header().Set(key, strings.Join(values, ","))
	}
	w.WriteHeader(int(fileInfo.Envelope.Status))
	_, err = w.Write(fileInfo.Content)
	if err != nil {
		slog.Error("[ PROXY SERVER ] Error writing response", slog.String("error", err.Error()))
	}
}

func (proxy *CacheableProxy) Listen(_ context.Context, startFeedback chan string) error {
	http.HandleFunc("/", proxy.Handler)
	if proxy.port == 0 {
		listener, err := proxy.prepareListener()
		if err != nil {
			return err
		}

		startFeedback <- proxy.ServeHost()

		defer listener.Close()
		return http.Serve(listener, nil)
	}

	startFeedback <- proxy.ServeHost()
	return http.ListenAndServe(proxy.ServeHost(), nil)
}

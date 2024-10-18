package crawler

import (
	"net/url"
	"slices"
)

type (
	InterceptorCallback = func(htmlStr string, target *url.URL)
	HTMLPipe            chan string
	htmlObserverEntry   struct {
		identifierKey string
		pipe          HTMLPipe
	}
)

type Interceptor struct {
	urlMultiplexer map[url.URL][]htmlObserverEntry
}

func NewInterceptor() *Interceptor {
	return &Interceptor{urlMultiplexer: make(map[url.URL][]htmlObserverEntry, 11)}
}

func (i Interceptor) observers(target *url.URL) (observersList []htmlObserverEntry) {
	// Try to get observers with base path
	lookTo := url.URL{Host: target.Host, RawPath: target.RawPath, Path: target.Path}
	foundObservers := i.urlMultiplexer[lookTo]
	observersList = slices.Clone(foundObservers)

	// Try to get observers with only host
	lookTo = url.URL{Host: target.Host}
	foundObservers = i.urlMultiplexer[lookTo]
	observersList = append(observersList, foundObservers...)
	return
}

func (i Interceptor) CreateObserver(
	on *url.URL, useBasePath bool,
	optWatcherKey ...string,
) (pipe HTMLPipe) {
	key := url.URL{Host: on.Host}
	if useBasePath {
		key.Path = on.Path
		key.RawPath = on.RawPath
	}

	// Check if it has a watcher key, and search for it if does
	var watcherKey string
	if len(optWatcherKey) > 0 {
		watcherKey = optWatcherKey[0]
	}

	if watcherKey != "" {
		for _, entries := range i.urlMultiplexer[key] {
			if watcherKey == entries.identifierKey {
				pipe = entries.pipe
				return
			}
		}
	}

	pipe = make(HTMLPipe, 3)
	i.urlMultiplexer[key] = append(
		i.urlMultiplexer[key],
		htmlObserverEntry{identifierKey: watcherKey, pipe: pipe},
	)
	return
}

func (i Interceptor) HandleResponse(htmlStr string, target *url.URL) {
	observersList := i.observers(target)
	for _, entry := range observersList {
		entry.pipe <- htmlStr
	}
}

func (i Interceptor) Dispose() {
	for _, observerEntries := range i.urlMultiplexer {
		for _, entry := range observerEntries {
			close(entry.pipe)
		}
	}
}

package cacheproxy

import (
	"net/http"

	"github.com/jictyvoo/tcg_deck-resolver/pkg/httptransport"
)

func (proxy CacheableProxy) Director(req *http.Request) {
	req.URL.Scheme = proxy.targetURL.Scheme
	req.URL.Host = proxy.targetURL.Host
	req.Host = proxy.targetURL.Host
	return
}

func (proxy CacheableProxy) RedirectRoundTripper() http.RoundTripper {
	return httptransport.NewTransportRewrite(
		proxy.targetURL, "localhost"+proxy.serveHost(),
	)
}

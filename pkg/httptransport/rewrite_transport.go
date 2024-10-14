package httptransport

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type TransportRewrite struct {
	originRoute   *url.URL
	redirectRoute string
	Transport     http.RoundTripper
}

func NewTransportRewrite(originRoute *url.URL, redirectRoute string) *TransportRewrite {
	return &TransportRewrite{
		originRoute:   originRoute,
		redirectRoute: redirectRoute,
		Transport:     defaultTransport,
	}
}

func (t *TransportRewrite) RoundTrip(req *http.Request) (*http.Response, error) {
	// Check if the request URL matches the domain
	if strings.Contains(req.URL.Host, t.originRoute.Host) {
		// Rewrite the request URL to localhost
		slog.Info(
			"Redirecting host URL",
			slog.String("host", t.redirectRoute),
			slog.String("origin", req.URL.Host),
		)
		req.URL.Host = t.redirectRoute
		req.URL.Scheme = "http"
	}

	// Call the next transport (which sends the request)
	return t.Transport.RoundTrip(req)
}

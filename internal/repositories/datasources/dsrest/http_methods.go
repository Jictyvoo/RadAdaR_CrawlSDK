package dsrest

import (
	"io"
	"net/http"
)

type HTTPMethod string

// Common HTTP methods.
//
// Unless otherwise noted, these are defined in RFC 7231 section 4.3.
const (
	MethodGet     HTTPMethod = http.MethodGet
	MethodHead    HTTPMethod = http.MethodHead
	MethodPost    HTTPMethod = http.MethodPost
	MethodPut     HTTPMethod = http.MethodPut
	MethodPatch   HTTPMethod = http.MethodPatch
	MethodDelete  HTTPMethod = http.MethodDelete
	MethodConnect HTTPMethod = http.MethodConnect
	MethodOptions HTTPMethod = http.MethodOptions
	MethodTrace   HTTPMethod = http.MethodTrace
)

type HTTPResponse struct {
	StatusCode int
	Body       io.Reader
	Headers    map[string][]string
}

type RESTDataSource interface {
	// Get performs an HTTP GET request to the specified URL.
	Get(url string, headers map[string][]string) (HTTPResponse, error)

	// Head performs an HTTP HEAD request to the specified URL.
	Head(url string, headers map[string][]string) (HTTPResponse, error)

	// Delete performs an HTTP DELETE request to the specified URL.
	Delete(url string, headers map[string][]string) (HTTPResponse, error)

	// Post performs an HTTP POST request to the specified URL with the given body.
	Post(url string, headers map[string][]string, data []byte) (HTTPResponse, error)

	// Put performs an HTTP PUT request to the specified URL with the given body.
	Put(url string, headers map[string][]string, data []byte) (HTTPResponse, error)

	// Patch performs an HTTP PATCH request to the specified URL with the given body.
	Patch(url string, headers map[string][]string, data []byte) (HTTPResponse, error)
}

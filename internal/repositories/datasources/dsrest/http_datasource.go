package dsrest

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/jictyvoo/tcg_deck-resolver/pkg/httptransport"
)

type HTTPDatasource struct {
	client *http.Client
}

func NewHTTPDatasource(roundTripper http.RoundTripper) *HTTPDatasource {
	if roundTripper == nil {
		roundTripper = httptransport.DefaultTransport
	}
	client := &http.Client{
		Transport: roundTripper,
		Timeout:   10 * time.Second, // Set timeout
	}
	return &HTTPDatasource{client: client}
}

func (d HTTPDatasource) attemptRequest(
	method HTTPMethod, url string, headers http.Header, data []byte,
) (HTTPResponse, error) {
	var bodyBuffer io.Reader
	if len(data) > 0 {
		bodyBuffer = bytes.NewReader(data)
	}

	// Creating request
	req, err := http.NewRequest(string(method), url, bodyBuffer)
	req.Header = headers

	if err != nil {
		return HTTPResponse{}, err
	}

	var resp *http.Response
	if resp, err = d.client.Do(req); err != nil {
		return HTTPResponse{}, err
	}

	respStruct := HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       resp.Body,
		Headers:    resp.Header,
	}

	return respStruct, nil
}

// Get performs an HTTP GET request to the specified URL.
func (d HTTPDatasource) Get(url string, headers map[string][]string) (HTTPResponse, error) {
	return d.attemptRequest(MethodGet, url, headers, nil)
}

// Head performs an HTTP HEAD request to the specified URL.
func (d HTTPDatasource) Head(url string, headers map[string][]string) (HTTPResponse, error) {
	return d.attemptRequest(MethodGet, url, headers, nil)
}

// Delete performs an HTTP DELETE request to the specified URL.
func (d HTTPDatasource) Delete(url string, headers map[string][]string) (HTTPResponse, error) {
	return d.attemptRequest(MethodGet, url, headers, nil)
}

// Post performs an HTTP POST request to the specified URL with the given body.
func (d HTTPDatasource) Post(
	url string,
	headers map[string][]string,
	data []byte,
) (HTTPResponse, error) {
	return d.attemptRequest(MethodGet, url, headers, data)
}

// Put performs an HTTP PUT request to the specified URL with the given body.
func (d HTTPDatasource) Put(
	url string,
	headers map[string][]string,
	data []byte,
) (HTTPResponse, error) {
	return d.attemptRequest(MethodGet, url, headers, data)
}

// Patch performs an HTTP PATCH request to the specified URL with the given body.
func (d HTTPDatasource) Patch(
	url string,
	headers map[string][]string,
	data []byte,
) (HTTPResponse, error) {
	return d.attemptRequest(MethodGet, url, headers, data)
}

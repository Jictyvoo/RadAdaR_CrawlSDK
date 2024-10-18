package httptransport

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestRoundTripper struct {
	name                         string
	redirectHost                 string
	originHost                   string
	originAssert, redirectAssert func(originHost, redirectHost string, w http.ResponseWriter, req *http.Request)
	expectedStatus               int
	assertThirdServer            func(w http.ResponseWriter, req *http.Request)
}

func (tester TestRoundTripper) Test(t *testing.T) {
	originServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tester.originAssert(tester.originHost, tester.redirectHost, w, r)
		}),
	)
	defer originServer.Close()

	redirectServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tester.redirectAssert(tester.originHost, tester.redirectHost, w, r)
		}),
	)
	defer redirectServer.Close()

	// Parse origin and redirect URLs
	originURL, err := url.Parse(originServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	var redirectURL *url.URL
	if redirectURL, err = url.Parse(redirectServer.URL); err != nil {
		t.Fatal(err)
	}

	// Create the TransportRewrite with the original and redirect route
	tester.redirectHost = redirectURL.Host
	tester.originHost = originURL.Host
	tr := NewTransportRewrite(originURL, tester.redirectHost)

	reqURL := originServer.URL + "/resource"
	// Create a request to the original host
	if tester.assertThirdServer != nil {
		thirdServer := httptest.NewServer(http.HandlerFunc(tester.assertThirdServer))
		defer thirdServer.Close()

		reqURL = thirdServer.URL + "/resource"
	}
	req := httptest.NewRequest(http.MethodGet, reqURL, nil)

	// Rewrite the request and send it using RoundTrip
	var resp *http.Response
	if resp, err = tr.RoundTrip(req); err != nil {
		t.Fatalf("Unexpected error in RoundTrip: %v", err)
	}

	// Verify that the response status code is OK
	if resp.StatusCode != tester.expectedStatus {
		t.Errorf("Expected status code `%d`, but got %s", tester.expectedStatus, resp.Status)
	}
}

func TestTransportRewrite_RoundTrip(t *testing.T) {
	// Test RoundTrip rewriting the URL when originRoute matches
	testList := []TestRoundTripper{
		{
			name: "Request Origin Matcher Rewrite",
			originAssert: func(_, _ string, w http.ResponseWriter, r *http.Request) {
				t.Errorf("Request should not have reached the original server")
			},
			redirectAssert: func(originHOST, redirectHOST string, w http.ResponseWriter, r *http.Request) {
				// Check if the request was rewritten correctly to the redirect host
				if r.URL.Host != "" && r.URL.Host != redirectHOST {
					t.Errorf(
						"Expected request to have Host '%s', but got '%s'",
						redirectHOST, r.Host,
					)
					return
				}

				if r.Host != originHOST {
					t.Errorf(
						"Expected request to have original Host '%s', but got '%s'",
						originHOST, r.Host,
					)
				}
				w.WriteHeader(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Request Unmatched Origin",
			originAssert: func(_, _ string, w http.ResponseWriter, r *http.Request) {
				t.Errorf("Request should not have reached the original server")
			},
			redirectAssert: func(_, _ string, w http.ResponseWriter, r *http.Request) {
				t.Errorf("Request should not have reached the redirect server")
			},
			assertThirdServer: func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			},
			expectedStatus: http.StatusTeapot,
		},
	}

	for _, tCase := range testList {
		t.Run(tCase.name, tCase.Test)
	}
}

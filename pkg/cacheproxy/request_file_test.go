package cacheproxy

import (
	"net/http"
	"net/url"
	"testing"
)

func TestCacheKey(t *testing.T) {
	// Define the base target URL for the proxy
	baseURL, _ := url.Parse("http://example.com")

	// Define test cases
	tests := []struct {
		name        string
		request     *http.Request
		expectedKey string
	}{
		{
			name: "Simple GET request without query",
			request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/file.txt", RawQuery: ""},
			},
			expectedKey: "file://GET@http://example.com#/file.txt",
		},
		{
			name: "GET request with query parameters",
			request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/file.txt", RawQuery: "foo=bar"},
			},
			expectedKey: "file://GET@http://example.com#/file.txtfoo=bar",
		},
		{
			name: "POST request without query",
			request: &http.Request{
				Method: "POST",
				URL:    &url.URL{Path: "/upload", RawQuery: ""},
			},
			expectedKey: "file://POST@http://example.com#/upload",
		},
		{
			name: "GET request with unescaped query",
			request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/file.txt", RawQuery: "name=John%20Doe"},
			},
			expectedKey: "file://GET@http://example.com#/file.txtname=John Doe",
		},
		{
			name: "GET request with invalid unescape query",
			request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/file.txt", RawQuery: "%ZZ"},
			},
			expectedKey: "file://GET@http://example.com#/file.txt%ZZ", // Raw query used if unescape fails
		},
		{
			name: "GET request with whitespace in path",
			request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/file%20with%20space.txt", RawQuery: ""},
			},
			expectedKey: "file://GET@http://example.com#/file%20with%20space.txt",
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the proxy with the base URL
			proxy := CacheableProxy{
				targetURL: baseURL,
			}

			// Call the cacheKey method
			actualKey := proxy.cacheKey(tt.request)

			// Compare the result with the expected key
			if actualKey != tt.expectedKey {
				t.Errorf("Test %s failed. Expected %q, got %q", tt.name, tt.expectedKey, actualKey)
			}
		})
	}
}

func TestIsFileTracked(t *testing.T) {
	// Define test cases in a slice
	tests := []struct {
		name              string
		trackedExtensions []string
		fileInfo          FileInformation
		expected          bool
	}{
		{
			name:              "ExtensionMatch",
			trackedExtensions: []string{".jpg", ".png"},
			fileInfo: FileInformation{
				FileMIME: FileMIME{Extension: ".jpg", MimeType: "image/jpeg"},
			},
			expected: true,
		},
		{
			name:              "MimeTypeMatch",
			trackedExtensions: []string{"image/jpeg", "image/png"},
			fileInfo: FileInformation{
				FileMIME: FileMIME{Extension: ".jpeg", MimeType: "image/jpeg"},
			},
			expected: true,
		},
		{
			name:              "NoMatch",
			trackedExtensions: []string{".gif", "image/gif"},
			fileInfo: FileInformation{
				FileMIME: FileMIME{Extension: ".jpg", MimeType: "image/jpeg"},
			},
			expected: false,
		},
		{
			name:              "CaseInsensitiveExtensionMatch",
			trackedExtensions: []string{".JPG", ".PNG"},
			fileInfo: FileInformation{
				FileMIME: FileMIME{Extension: ".jpg", MimeType: "image/jpeg"},
			},
			expected: true,
		},
		{
			name:              "CaseInsensitiveMimeTypeMatch",
			trackedExtensions: []string{".gif", "IMAGE/JPEG"},
			fileInfo: FileInformation{
				FileMIME: FileMIME{Extension: ".jpeg", MimeType: "image/jpeg"},
			},
			expected: true,
		},
	}

	// Iterate through the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the CacheableProxy
			proxy := CacheableProxy{
				trackedExtensions: tt.trackedExtensions,
			}

			// Call the function and check the result
			result := proxy.isFileTracked(tt.fileInfo)
			if result != tt.expected {
				t.Errorf("For %s, expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
}

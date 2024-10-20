package cacheproxy

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"
	"testing"
)

func TestFileMIME(t *testing.T) {
	tests := []struct {
		name         string
		respBody     []byte
		header       http.Header
		expectedMIME string
	}{
		{
			name:         "Content-Type set in header",
			respBody:     []byte{},
			header:       http.Header{"Content-Type": []string{"image/png"}},
			expectedMIME: "image/png",
		},
		{
			name:         "Content-Type is application/octet-stream (fallback to DetectContentType)",
			respBody:     []byte("This is some text content"),
			header:       http.Header{"Content-Type": []string{"application/octet-stream"}},
			expectedMIME: "text/plain; charset=utf-8",
		},
		{
			name:         "No Content-Type in header",
			respBody:     []byte("<html><body></body></html>"),
			header:       http.Header{},
			expectedMIME: "text/html; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mimeType := fileMIME(tt.respBody, tt.header)
			if mimeType != tt.expectedMIME {
				t.Errorf("Expected %s, but got %s", tt.expectedMIME, mimeType)
			}
		})
	}
}

func TestChecksum(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		expected string
	}{
		{
			name:     "Empty body",
			body:     []byte{},
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "Short string body",
			body:     []byte("hello"),
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:     "Longer body",
			body:     []byte("this is a longer body to test the checksum function"),
			expected: "cd65a2e8b563ee4daf85c042cc2ad1545ff1e0d061774eec40e8befdd8c8c8b4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checksum(tt.body)
			resultHex := hex.EncodeToString(result)
			if resultHex != tt.expected {
				t.Errorf("Expected %x, but got %x", tt.expected, result)
			}
		})
	}
}

func TestBodyReader(t *testing.T) {
	tests := []struct {
		name        string
		respBody    io.Reader
		expected    []byte
		expectError bool
	}{
		{
			name:        "Valid body",
			respBody:    bytes.NewBuffer([]byte("this is the response body")),
			expected:    []byte("this is the response body"),
			expectError: false,
		},
		{
			name:        "Empty body",
			respBody:    bytes.NewBuffer([]byte{}),
			expected:    []byte{},
			expectError: false,
		},
		{
			name:        "Error during read (simulated by broken reader)",
			respBody:    errReader{}, // This simulates a read error
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bodyReader(tt.respBody)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v", tt.expectError, err != nil)
			}
			if string(result) != string(tt.expected) {
				t.Errorf("Expected %s, but got %s", tt.expected, result)
			}
		})
	}
}

// errReader is an io.Reader that always returns an error.
type errReader struct{}

func (errReader) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

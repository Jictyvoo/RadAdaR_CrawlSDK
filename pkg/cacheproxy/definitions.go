package cacheproxy

import "time"

type KVStorage[T any, K comparable] interface {
	Set(key K, value T) error
	Get(key K) (T, error)
}

type (
	FileMIME struct {
		Name      string
		Extension string
		MimeType  string
	}
	FileEnvelope struct {
		Headers map[string][]string
		Status  uint16
	}
	FileInformation struct {
		FileMIME
		Envelope      FileEnvelope
		Content       []byte
		Checksum      []byte
		CreatedAt     time.Time
		ModifiedAt    time.Time
		ExtraMetadata map[string]string
	}
)

package badgerepo

import (
	"os"
	"reflect"
	"testing"

	"github.com/jictyvoo/radadar_crawlsdk/pkg/cacheproxy"
)

// Helper function to create a temporary directory for the database path
func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "TMP@badger__test.")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(dir) // Cleanup after the test
	})
	return dir
}

// Test NewRemoteFileCache to ensure database is created successfully
func TestNewRemoteFileCache(t *testing.T) {
	dbPath := createTempDir(t)

	cache, err := NewRemoteFileCache(dbPath)
	if err != nil {
		t.Fatalf("Failed to create RemoteFileCache: %v", err)
	}
	defer cache.Close()

	// Ensure the database directory exists
	if _, err = os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Database path does not exist: %v", err)
	}
}

// Test Set and Get operations in the RemoteFileCache
func TestRemoteFileCache_SetAndGet(t *testing.T) {
	dbPath := createTempDir(t)

	// Create a new cache
	cache, err := NewRemoteFileCache(dbPath)
	if err != nil {
		t.Fatalf("Failed to create RemoteFileCache: %v", err)
	}
	defer cache.Close()

	// Define a key and value
	const key = "setKEY"
	fileInfo := fixtureFileInfo()

	// Test Set operation
	err = cache.Set(key, fileInfo)
	if err != nil {
		t.Fatalf("Failed to set key in cache: %v", err)
	}

	// Test Get operation
	var retrievedFileInfo cacheproxy.FileInformation
	if retrievedFileInfo, err = cache.Get(key); err != nil {
		t.Fatalf("Failed to get key from cache: %v", err)
	}

	if !reflect.DeepEqual(fileInfo, retrievedFileInfo) {
		t.Errorf("retrieved value is not valid: expected %v, got %v", fileInfo, retrievedFileInfo)
	}
}

// Test Keys operation to ensure all stored keys are retrieved
func TestRemoteFileCache_Keys(t *testing.T) {
	dbPath := createTempDir(t)

	// Create a new cache
	cache, err := NewRemoteFileCache(dbPath)
	if err != nil {
		t.Fatalf("Failed to create RemoteFileCache: %v", err)
	}
	defer cache.Close()

	// Define multiple key-value pairs
	fileInfo := fixtureFileInfo()
	keys := []string{"fetched@01_A", "fetched@02_B", "fetched@03_C"}

	// Set key-value pairs in the cache
	for _, key := range keys {
		if err = cache.Set(key, fileInfo); err != nil {
			t.Fatalf("Failed to set key %v in cache: %v", key, err)
		}
	}

	// Retrieve all keys
	var retrievedKeys []string
	if retrievedKeys, err = cache.Keys(); err != nil {
		t.Fatalf("Failed to retrieve keys: %v", err)
	}

	// Ensure all keys were retrieved
	if len(retrievedKeys) != len(keys) {
		t.Errorf("Expected %v keys, got %v", len(keys), len(retrievedKeys))
	}

	// Verify each key is present in the retrieved keys
	keyMap := make(map[string]bool, len(retrievedKeys))
	for _, key := range retrievedKeys {
		keyMap[key] = true
	}
	for _, key := range keys {
		if !keyMap[key] {
			t.Errorf("Expected key %v was not found in retrieved keys", key)
		}
	}
}

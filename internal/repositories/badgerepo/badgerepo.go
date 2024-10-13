package badgerepo

import (
	"log/slog"
	"time"

	"github.com/dgraph-io/badger/v4"

	"github.com/jictyvoo/tcg_deck-resolver/pkg/cacheproxy"
)

type RemoteFileCache struct {
	db       *badger.DB
	entryTTL time.Duration
}

// NewRemoteFileCache initializes a new Badger database instance for the RemoteFileCache
func NewRemoteFileCache(dbPath string) (*RemoteFileCache, error) {
	// Set up Badger options and open the database
	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &RemoteFileCache{db: db, entryTTL: 36 * time.Hour}, nil
}

// Set stores a key-value pair in the Badger database
func (r *RemoteFileCache) Set(key string, information cacheproxy.FileInformation) error {
	err := r.db.Update(func(txn *badger.Txn) error {
		valBytes, encodErr := EncodeFileInfo(information)
		if encodErr != nil {
			return encodErr
		}
		badgerEntry := badger.NewEntry([]byte(key), valBytes).WithTTL(r.entryTTL)
		return txn.SetEntry(badgerEntry)
	})

	return err
}

// Get retrieves a value by key from the Badger database
func (r *RemoteFileCache) Get(key string) (cacheproxy.FileInformation, error) {
	var valCopy []byte
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		// Retrieve the value and copy it
		valCopy, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		slog.Error("Failed to get value", slog.String("key", key), slog.String("err", err.Error()))
	}

	return DecodeFileInfo(valCopy)
}

// Close closes the Badger database
func (r *RemoteFileCache) Close() error {
	return r.db.Close()
}

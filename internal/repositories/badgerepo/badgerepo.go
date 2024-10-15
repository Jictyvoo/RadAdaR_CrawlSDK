package badgerepo

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"

	"github.com/jictyvoo/tcg_deck-resolver/pkg/cacheproxy"
)

type RemoteFileCache struct {
	db            *badger.DB
	entryTTL      time.Duration
	mutex         sync.Mutex
	finishThreads context.CancelFunc
}

// NewRemoteFileCache initializes a new Badger database instance for the RemoteFileCache
func NewRemoteFileCache(dbPath string) (*RemoteFileCache, error) {
	// Set up Badger options and open the database
	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	go gcThread(ctx, db)
	return &RemoteFileCache{db: db, entryTTL: 36 * time.Hour, finishThreads: cancel}, nil
}

// Set stores a key-value pair in the Badger database
func (r *RemoteFileCache) Set(key string, information cacheproxy.FileInformation) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	err := r.db.Update(func(txn *badger.Txn) error {
		valBytes, encodErr := EncodeFileInfo(information)
		if encodErr != nil {
			return encodErr
		}
		badgerEntry := badger.NewEntry([]byte(key), valBytes).WithTTL(r.entryTTL).WithDiscard()
		return txn.SetEntry(badgerEntry)
	})

	return err
}

// Get retrieves a value by key from the Badger database
func (r *RemoteFileCache) Get(key string) (cacheproxy.FileInformation, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

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
		return cacheproxy.FileInformation{}, err
	}

	slog.Info("Successfully loaded data from cache", slog.String("key", key))
	return DecodeFileInfo(valCopy)
}

// Close closes the Badger database
func (r *RemoteFileCache) Close() error {
	if r.finishThreads != nil {
		r.finishThreads()
	}
	return r.db.Close()
}

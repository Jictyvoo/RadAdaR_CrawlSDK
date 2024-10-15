package badgerepo

import (
	"context"
	"time"

	"github.com/dgraph-io/badger/v4"
)

func gcThread(ctx context.Context, db *badger.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = runGC(db)
		}
	}
}

func runGC(db *badger.DB) (err error) {
	for {
		if err = db.RunValueLogGC(0.7); err != nil {
			return
		}
	}
}

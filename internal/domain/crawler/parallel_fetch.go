package crawler

import (
	"context"
	"iter"
	"log/slog"
	"sync"

	"github.com/jictyvoo/radadar_crawlsdk/pkg/datatypes"
)

type (
	FetchDatasource interface {
		DownloadPage(url string) (string, error)
		Close() error
	}
	ParallelFetch struct {
		factory       datatypes.Factory[FetchDatasource]
		urlQueue      chan string
		outputQueue   chan FetchResult
		ctx           context.Context
		cancel        context.CancelFunc
		wg            sync.WaitGroup
		instanceMutex sync.Mutex
	}
)

func NewParallelFetch(factory datatypes.Factory[FetchDatasource]) *ParallelFetch {
	const bufferLen = 3
	p := &ParallelFetch{
		factory:     factory,
		urlQueue:    make(chan string, bufferLen),
		outputQueue: make(chan FetchResult, bufferLen),
	}
	return p
}

func (pf *ParallelFetch) Start(ctx context.Context, numWorkers int) error {
	pf.instanceMutex.Lock()
	defer pf.instanceMutex.Unlock()

	// Finish an older start run
	if pf.cancel != nil {
		pf.cancel()
	}

	pf.ctx, pf.cancel = context.WithCancel(ctx)
	pf.wg.Add(numWorkers)

	for range numWorkers {
		datasource, err := pf.factory.New()
		if err != nil {
			slog.Error("Failed to create new datasource instance", slog.String("err", err.Error()))
			pf.cancel() // Finish already created workers
			return err
		}
		newWorker := fetchWorker{
			inputCh:    pf.urlQueue,
			outputCh:   pf.outputQueue,
			datasource: datasource,
		}
		go newWorker.Run(pf.ctx, pf.wg.Done)
	}

	return nil
}

func (pf *ParallelFetch) Fetch(urls ...string) {
	for _, url := range urls {
		pf.urlQueue <- url
	}

	return
}

// Responses is an iterator over the elements received on output channel.
func (pf *ParallelFetch) Responses() iter.Seq2[string, FetchResult] {
	return func(yield func(string, FetchResult) bool) {
		for v := range pf.outputQueue {
			if !yield(v.url, v) {
				return
			}
		}
	}
}

func (pf *ParallelFetch) Stop() {
	pf.instanceMutex.Lock()
	defer pf.instanceMutex.Unlock()

	pf.cancel()  // cancel all workers
	pf.wg.Wait() // wait until when all workers are done
}

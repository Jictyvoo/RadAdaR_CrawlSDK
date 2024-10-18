package crawler

import (
	"context"
	"errors"
)

type (
	FetchResult struct {
		Err      error
		url      string
		RespBody string
	}
	fetchWorker struct {
		inputCh    chan string
		outputCh   chan FetchResult
		datasource FetchDatasource
	}
)

var ErrInputChannelClosed = errors.New("input channel closed")

func (w fetchWorker) Run(ctx context.Context, onFinishCallback func()) (err error) {
	defer onFinishCallback()
	defer func(datasource FetchDatasource) {
		closErr := datasource.Close()
		if closErr != nil {
			err = errors.Join(err, closErr)
		}
	}(w.datasource)

	for {
		select {
		// Stop worker if context is done
		case <-ctx.Done():
			return ctx.Err()

		// Main loop use case while receiving from the channel
		case url, ok := <-w.inputCh:
			if !ok {
				return ErrInputChannelClosed
			}
			respBody, downErr := w.datasource.DownloadPage(url)
			w.outputCh <- FetchResult{
				Err:      downErr,
				url:      url,
				RespBody: respBody,
			}
		}
	}
}

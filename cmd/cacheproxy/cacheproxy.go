package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jictyvoo/tcg_deck-resolver/internal/repositories/badgerepo"
	"github.com/jictyvoo/tcg_deck-resolver/pkg/cacheproxy"
)

func main() {
	var (
		port      uint
		targetURL string
	)
	flag.UintVar(&port, "port", 0, "port to listen on")
	flag.StringVar(&targetURL, "target-url", "", "target URL")
	flag.Parse()

	if targetURL == "" {
		flag.Usage()
		return
	}

	repo, err := badgerepo.NewRemoteFileCache("http_cache.badger")
	if err != nil {
		slog.Error("failed to init badger repo", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer repo.Close()

	var proxy *cacheproxy.CacheableProxy
	if proxy, err = cacheproxy.New(
		repo, targetURL, uint16(port),
	); err != nil {
		slog.Error(
			"failed to initialize cacheproxy",
			slog.String("url", targetURL),
			slog.String("port", strconv.Itoa(int(port))),
			slog.String("err", err.Error()),
		)
		return
	}

	var (
		startFeedback = make(chan string, 1)
		wg            sync.WaitGroup
		ctx           = gracefulShutdown()
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if listErr := proxy.Listen(ctx, startFeedback); listErr != nil &&
			!errors.Is(listErr, context.Canceled) {
			log.Fatal(listErr)
		}
	}()
	<-startFeedback
	close(startFeedback)

	// Wait for the server to start (feedback can be expanded later if needed)
	slog.Info(fmt.Sprintf("Listening on address %s", proxy.ServeHost()))

	// Wait for all goroutines to finish
	wg.Wait()

	// Ensure a graceful exit with a small delay
	slog.Info("Shutdown complete, exiting...")
	time.Sleep(1 * time.Second)
}

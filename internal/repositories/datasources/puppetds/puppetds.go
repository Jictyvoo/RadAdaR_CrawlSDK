package puppetds

import (
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type BrowserPuppetDatasource struct {
	roundTripper http.RoundTripper
}

func NewPuppetDatasource(roundTripper http.RoundTripper) *BrowserPuppetDatasource {
	return &BrowserPuppetDatasource{roundTripper: roundTripper}
}

func (wpr BrowserPuppetDatasource) DownloadPage(baseURL string) (string, error) {
	// Headless runs the browser on foreground, you can also use flag "-rod=show"
	// Devtools opens the tab in each new tab opened automatically
	l := launcher.New()
	/*l = l.Headless(false).
	Devtools(true)*/

	defer l.Cleanup()

	url := l.MustLaunch()

	// Trace shows verbose debug information for each action executed
	// SlowMotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(url).
		// Trace(true).
		// SlowMotion(2 * time.Second).
		MustConnect()

	// Even you forget to close, rod will close it after main process ends.
	defer browser.MustClose()

	router := browser.HijackRequests().MustAdd("*", func(ctx *rod.Hijack) {
		client := &http.Client{Transport: wpr.roundTripper, Timeout: time.Hour}
		_ = ctx.LoadResponse(client, true)
	})
	go router.Run()
	defer router.Stop()

	// Create a new page
	page, err := browser.Page(proto.TargetCreateTarget{URL: baseURL})
	if err != nil {
		return "", err
	}
	if err = page.WaitStable(3 * time.Second); err != nil {
		return "", err
	}
	pageHeight := page.MustEval("() => document.body.scrollHeight").Int()
	windowHeight := page.MustEval(`() => window.innerHeight`).Int()

	totalChunks := pageHeight / (windowHeight * 2)
	for chunk := 0; chunk < totalChunks; chunk += 1 {
		page.Mouse.MustScroll(0, float64(windowHeight*chunk))
		time.Sleep(time.Millisecond << 8)
	}
	return page.HTML()
}

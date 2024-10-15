package puppetds

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type rodBrowser struct {
	launcher *launcher.Launcher
	browser  *rod.Browser
	router   *rod.HijackRouter
}

func newBrowser(roundTripper http.RoundTripper) (rodBrowser, error) {
	rb := rodBrowser{}
	// Headless runs the browser on foreground, you can also use flag "-rod=show"
	// Devtools opens the tab in each new tab opened automatically
	rb.launcher = launcher.New()
	/*l = l.Headless(false).
	Devtools(true)*/

	url, err := rb.launcher.Launch()
	if err != nil {
		return rodBrowser{}, err
	}

	// Trace shows verbose debug information for each action executed
	// SlowMotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	rb.browser = rod.New().
		ControlURL(url).
		// Trace(true).
		// SlowMotion(2 * time.Second).
		MustConnect()

	rb.router = rb.browser.HijackRequests()
	err = rb.router.Add("*", "", func(ctx *rod.Hijack) {
		client := &http.Client{Transport: roundTripper, Timeout: time.Second << 8}
		_ = ctx.LoadResponse(client, true)
	})
	if err != nil {
		return rodBrowser{}, err
	}

	go rb.router.Run()
	return rb, nil
}

func (rb rodBrowser) Page(url string) (*rod.Page, error) {
	return rb.browser.Page(proto.TargetCreateTarget{URL: url})
}

func (rb rodBrowser) Close() error {
	defer rb.launcher.Cleanup()
	return errors.Join(
		rb.router.Stop(),
		// Even you forget to close, rod will close it after main process ends.
		rb.browser.Close(),
	)
}

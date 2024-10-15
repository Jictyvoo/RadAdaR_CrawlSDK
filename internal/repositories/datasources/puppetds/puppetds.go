package puppetds

import (
	"net/http"
	"time"
)

type BrowserPuppetDatasource struct {
	roundTripper http.RoundTripper
	browser      rodBrowser
}

func NewBrowserDatasource(roundTripper http.RoundTripper) (*BrowserPuppetDatasource, error) {
	browser, err := newBrowser(roundTripper)
	if err != nil {
		return nil, err
	}
	return &BrowserPuppetDatasource{roundTripper: roundTripper, browser: browser}, nil
}

func (wpr BrowserPuppetDatasource) DownloadPage(baseURL string) (string, error) {
	// Create a new page
	page, err := wpr.browser.Page(baseURL)
	if err != nil {
		return "", err
	}
	if err = page.WaitStable(3 * time.Second); err != nil {
		return "", err
	}

	err = scrollPageDown(page)
	return page.HTML()
}

func (wpr BrowserPuppetDatasource) Close() error {
	return wpr.browser.Close()
}

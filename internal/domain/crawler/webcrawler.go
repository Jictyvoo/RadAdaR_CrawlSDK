package crawler

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/temoto/robotstxt"
)

type WebCrawler struct {
	baseURL   *url.URL
	userAgent string
}

func NewWebCrawler(baseURL, userAgent string) *WebCrawler {
	reqURL, _ := url.Parse(baseURL)
	return &WebCrawler{baseURL: reqURL, userAgent: userAgent}
}

func (wc WebCrawler) LoadRobotsTXT() (*robotstxt.RobotsData, error) {
	reqURL := url.URL{
		Scheme: wc.baseURL.Scheme,
		Host:   wc.baseURL.Host,
		Path:   "/robots.txt",
	}
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}
	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var body []byte
	if body, err = io.ReadAll(res.Body); err != nil {
		return nil, err
	}
	return robotstxt.FromBytes(body)
}

func (wc WebCrawler) Crawl() error {
	robots, err := wc.LoadRobotsTXT()
	if err != nil {
		return err
	}
	robotGroup := robots.FindGroup(wc.userAgent)
	for {
		time.Sleep(robotGroup.CrawlDelay)
	}
}

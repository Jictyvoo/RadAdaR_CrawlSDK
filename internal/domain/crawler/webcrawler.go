package crawler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/temoto/robotstxt"
)

type WebCrawler struct {
	baseURL *url.URL
}

func NewWebCrawler(baseURL string) *WebCrawler {
	reqURL, _ := url.Parse(baseURL)
	return &WebCrawler{baseURL: reqURL}
}

func (wc WebCrawler) LoadRobotsTXT() (*robotstxt.RobotsData, error) {
	reqURL := url.URL{
		Scheme: wc.baseURL.Scheme,
		Host:   wc.baseURL.Host,
		Path:   "/robots.txt",
	}
	req, _ := http.NewRequest("GET", reqURL.String(), nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return robotstxt.FromBytes(body)
}

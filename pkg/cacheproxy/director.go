package cacheproxy

import (
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/jictyvoo/radadar_crawlsdk/pkg/httptransport"
)

func (proxy *CacheableProxy) Director(req *http.Request) {
	req.URL.Scheme = proxy.targetURL.Scheme
	req.URL.Host = proxy.targetURL.Host
	req.Host = proxy.targetURL.Host
	return
}

func (proxy *CacheableProxy) RedirectRoundTripper() http.RoundTripper {
	return httptransport.NewTransportRewrite(
		proxy.targetURL, "localhost"+proxy.ServeHost(),
	)
}

func (proxy *CacheableProxy) ServeHost() string {
	return ":" + strconv.FormatUint(uint64(proxy.port), 10)
}

func (proxy *CacheableProxy) prepareListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", proxy.ServeHost())
	if err != nil {
		return nil, err
	}

	if proxy.port != 0 {
		return listener, nil
	}

	address := listener.Addr().String()
	var index int
	for index = len(address) - 1; index >= 0; index-- {
		if address[index] == ':' {
			break
		}
	}

	var port uint64
	if port, err = strconv.ParseUint(address[index+1:], 10, 16); err != nil {
		err = errors.Join(err, listener.Close())
		return nil, err
	}

	// Update port with bound on listener
	proxy.port = uint16(port)
	return listener, nil
}

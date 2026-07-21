package httpclient

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

func newTransport(proxyURL string) (*http.Transport, error) {
	if proxyURL == "" {
		return &http.Transport{}, nil
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("httpclient: invalid proxy url: %w", err)
	}
	switch u.Scheme {
	case "http", "https":
		return &http.Transport{Proxy: http.ProxyURL(u)}, nil
	case "socks5", "socks5h":
		dialer, err := proxy.FromURL(u, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("httpclient: invalid socks5 proxy: %w", err)
		}
		return &http.Transport{
			DialContext: func(_ context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}, nil
	default:
		return nil, fmt.Errorf("httpclient: unsupported proxy scheme %q", u.Scheme)
	}
}

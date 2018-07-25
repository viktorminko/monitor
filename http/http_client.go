package http

import (
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	Proxy *url.URL
}

func (h *Client) Init() (*http.Client, error) {
	client := &http.Client{}

	if h.Proxy == nil {
		return client, nil
	}

	var httpTransport *http.Transport

	switch strings.ToLower(h.Proxy.Scheme) {
	case "socks5":
		dialer, err := proxy.FromURL(h.Proxy, proxy.Direct)

		if err != nil {
			return nil, err
		}

		httpTransport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}}
	case "http":
		httpTransport = &http.Transport{Proxy: http.ProxyURL(h.Proxy)}

	default:
		return nil, fmt.Errorf("requested proxy scheme is not supported: %v", h.Proxy.Scheme)
	}

	client.Transport = httpTransport
	return client, nil
}

func (h *Client) Call(req *http.Request, timeout time.Duration) (*http.Response, error) {

	client, err := h.Init()
	if err != nil {
		return nil, err
	}

	if timeout != 0 {
		client.Timeout = timeout
	}

	log.Println("Sending request: ", req.Method, " ", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Options struct {
	Retry      int
	RetryDelay time.Duration
	UserAgent  string
	Proxy      string
	Timeout    time.Duration
}

type Request struct {
	Method  string
	URL     string
	Query   map[string]string
	Headers map[string]string
	Body    string
}

func Do(req Request, opts Options) (string, error) {
	return DoWithContext(req, opts, context.Background())
}

func DoWithContext(req Request, opts Options, ctx context.Context) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	transport, err := newTransport(opts.Proxy)
	if err != nil {
		return "", err
	}
	client := &http.Client{Transport: transport, Timeout: opts.Timeout}

	attempts := opts.Retry
	if attempts < 1 {
		attempts = 1
	}

	var lastErr error
	for i := 0; i < attempts; i++ {
		if i > 0 && opts.RetryDelay > 0 {
			time.Sleep(opts.RetryDelay)
		}
		body, err := do(client, req, opts.UserAgent, ctx)
		if err == nil {
			return body, nil
		}
		lastErr = err
	}
	return "", lastErr
}

func do(client *http.Client, r Request, userAgent string, ctx context.Context) (string, error) {
	method := r.Method
	if method == "" {
		method = http.MethodGet
	}

	reqURL := r.URL
	if len(r.Query) > 0 {
		q := url.Values{}
		for k, v := range r.Query {
			q.Set(k, v)
		}
		sep := "?"
		if strings.Contains(reqURL, "?") {
			sep = "&"
		}
		reqURL += sep + q.Encode()
	}

	var body io.Reader
	if r.Body != "" {
		body = strings.NewReader(r.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return "", fmt.Errorf("httpclient: %w", err)
	}
	for k, v := range r.Headers {
		httpReq.Header.Set(k, v)
	}
	if userAgent != "" {
		httpReq.Header.Set("User-Agent", userAgent)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("httpclient: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("httpclient: %w", err)
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("httpclient: unexpected status %d from %s", resp.StatusCode, reqURL)
	}
	return string(data), nil
}

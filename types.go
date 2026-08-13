// Package snapsave downloads Instagram, Facebook, TikTok and Twitter (X) media
package snapsave

import "time"

type MediaType string

const (
	MediaImage MediaType = "image"
	MediaVideo MediaType = "video"
)

type Media struct {
	Resolution   string    `json:"resolution,omitempty"`
	ShouldRender bool      `json:"shouldRender,omitempty"`
	Thumbnail    string    `json:"thumbnail,omitempty"`
	Type         MediaType `json:"type,omitempty"`
	URL          string    `json:"url,omitempty"`
}

type Data struct {
	Description string  `json:"description,omitempty"`
	Preview     string  `json:"preview,omitempty"`
	Media       []Media `json:"media,omitempty"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    *Data  `json:"data,omitempty"`
}

type Options struct {
	// Retry is the number of retry attempts on request failure. Default 1.
	Retry int
	// RetryDelay between retries. Default 500ms.
	RetryDelay time.Duration
	// UserAgent overrides the default request User-Agent.
	UserAgent string
	// Proxy is a proxy URL (http, https, or socks5), e.g. "http://user:pass@host:port".
	Proxy string
	// Timeout is the maximum time allowed for a request before it is aborted.
	// When set, the library stops waiting even if the upstream service hangs.
	Timeout time.Duration
}

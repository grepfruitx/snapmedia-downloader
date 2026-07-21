package parser

import (
	"net/url"
	"strings"
)

type Media struct {
	Resolution   string
	ShouldRender bool
	Thumbnail    string
	Type         string
	URL          string
}

type Result struct {
	Description string
	Preview     string
	Media       []Media
}

func FixThumbnail(u string) string {
	const prefix = "https://snapinsta.app/photo.php?photo="
	if !strings.Contains(u, prefix) {
		return u
	}
	decoded, err := url.QueryUnescape(strings.Replace(u, prefix, "", 1))
	if err != nil {
		return u
	}
	return decoded
}

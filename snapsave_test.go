package snapsave

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/grepfruitx/snapmedia-downloader/internal/parser"
)

func TestDownloadInvalidURL(t *testing.T) {
	resp := Download("https://www.example.com/not-supported", nil)
	if resp.Success || resp.Message != "Invalid URL" {
		t.Fatalf("got %+v", resp)
	}
}

func TestDownloadSnapSaveHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	orig := snapsaveActionURL
	snapsaveActionURL = srv.URL
	defer func() { snapsaveActionURL = orig }()

	resp := Download("https://www.facebook.com/watch?v=1234567890123456", &Options{Retry: 1})
	if resp.Success {
		t.Fatalf("expected failure, got %+v", resp)
	}
	if resp.Message == "" {
		t.Fatal("expected error message")
	}
}

func TestDownloadTikTokMissingToken(t *testing.T) {
	home := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html><body>no token here</body></html>`))
	}))
	defer home.Close()

	origHome := snaptikHomeURL
	snaptikHomeURL = home.URL
	defer func() { snaptikHomeURL = origHome }()

	resp := Download("https://www.tiktok.com/@user_name/video/1234567890123456789", &Options{Retry: 1})
	if resp.Success {
		t.Fatalf("expected failure, got %+v", resp)
	}
	if resp.Message == "" {
		t.Fatal("expected error message")
	}
}

func TestDownloadTwitterBadJSON(t *testing.T) {
	home := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<input type="hidden" name="token" value="tok123">`))
	}))
	defer home.Close()
	action := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer action.Close()

	origHome, origAction := twitterHomeURL, twitterActionURL
	twitterHomeURL, twitterActionURL = home.URL, action.URL
	defer func() { twitterHomeURL, twitterActionURL = origHome, origAction }()

	resp := Download("https://x.com/user_name/status/1234567890123456789", &Options{Retry: 1})
	if resp.Success {
		t.Fatalf("expected failure, got %+v", resp)
	}
	if resp.Message == "" {
		t.Fatal("expected error message")
	}
}

func TestResolveOptionsDefaults(t *testing.T) {
	o := resolveOptions(nil)
	if o.Retry != 1 || o.RetryDelay != 500*time.Millisecond || o.UserAgent != defaultUserAgent || o.Proxy != "" {
		t.Fatalf("got %+v", o)
	}

	o = resolveOptions(&Options{Retry: 5, RetryDelay: 2 * time.Second, UserAgent: "custom", Proxy: "http://proxy:8080"})
	if o.Retry != 5 || o.RetryDelay != 2*time.Second || o.UserAgent != "custom" || o.Proxy != "http://proxy:8080" {
		t.Fatalf("got %+v", o)
	}
}

func TestToData(t *testing.T) {
	d := toData(&parser.Result{
		Description: "desc",
		Preview:     "https://example.com/p.jpg",
		Media: []parser.Media{
			{URL: "https://example.com/v.mp4", Type: "video", Resolution: "720p", ShouldRender: true},
			{URL: "https://example.com/i.jpg", Type: "image", Thumbnail: "https://example.com/t.jpg"},
		},
	})
	if d.Description != "desc" || d.Preview != "https://example.com/p.jpg" || len(d.Media) != 2 {
		t.Fatalf("got %+v", d)
	}
	if d.Media[0].Type != MediaVideo || !d.Media[0].ShouldRender || d.Media[0].Resolution != "720p" {
		t.Errorf("media[0] = %+v", d.Media[0])
	}
	if d.Media[1].Type != MediaImage || d.Media[1].Thumbnail != "https://example.com/t.jpg" {
		t.Errorf("media[1] = %+v", d.Media[1])
	}
}

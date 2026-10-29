package snapsave

import (
	"context"
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

func TestDownloadFacebookHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	orig := savefSearchURL
	savefSearchURL = srv.URL
	defer func() { savefSearchURL = orig }()

	resp := Download("https://www.facebook.com/watch?v=1234567890123456", &Options{Retry: 1})
	if resp.Success {
		t.Fatalf("expected failure, got %+v", resp)
	}
	if resp.Message == "" {
		t.Fatal("expected error message")
	}
}

func TestDownloadTikTokAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":false,"message":"Video not found"}`))
	}))
	defer srv.Close()

	orig := tiktokAPIURL
	tiktokAPIURL = srv.URL
	defer func() { tiktokAPIURL = orig }()

	resp := Download("https://www.tiktok.com/@user_name/video/1234567890123456789", &Options{Retry: 1})
	if resp.Success || resp.Message != "Video not found" {
		t.Fatalf("got %+v", resp)
	}
}

func TestDownloadTikTokSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil || r.PostForm.Get("url") == "" {
			t.Errorf("expected url form field, got %v", r.PostForm)
		}
		_, _ = w.Write([]byte(`{"success":true,"title":"a video","cover":"https://example.com/c.jpg","video_nowm":"https://example.com/v.mp4"}`))
	}))
	defer srv.Close()

	orig := tiktokAPIURL
	tiktokAPIURL = srv.URL
	defer func() { tiktokAPIURL = orig }()

	resp := Download("https://www.tiktok.com/@user_name/video/1234567890123456789", &Options{Retry: 1})
	if !resp.Success {
		t.Fatalf("expected success, got %+v", resp)
	}
	if resp.Data.Description != "a video" || resp.Data.Preview != "https://example.com/c.jpg" {
		t.Fatalf("got %+v", resp.Data)
	}
	if len(resp.Data.Media) != 1 || resp.Data.Media[0].URL != "https://example.com/v.mp4" || resp.Data.Media[0].Type != MediaVideo {
		t.Fatalf("media = %+v", resp.Data.Media)
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
	if o.Retry != 1 || o.RetryDelay != 500*time.Millisecond || o.UserAgent != defaultUserAgent || o.Proxy != "" || o.Timeout != 15*time.Second {
		t.Fatalf("got %+v", o)
	}

	o = resolveOptions(&Options{Retry: 5, RetryDelay: 2 * time.Second, UserAgent: "custom", Proxy: "http://proxy:8080", Timeout: 3 * time.Second})
	if o.Retry != 5 || o.RetryDelay != 2*time.Second || o.UserAgent != "custom" || o.Proxy != "http://proxy:8080" || o.Timeout != 3*time.Second {
		t.Fatalf("got %+v", o)
	}

	o = resolveOptions(&Options{Timeout: 0})
	if o.Timeout != 0 {
		t.Fatalf("expected explicit zero timeout override, got %+v", o)
	}
}

func TestDownloadWithContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	orig := savefSearchURL
	savefSearchURL = srv.URL
	defer func() { savefSearchURL = orig }()

	resp := DownloadWithContext(ctx, "https://www.facebook.com/watch?v=1234567890123456", &Options{Retry: 1, Timeout: 200 * time.Millisecond})
	if resp.Success {
		t.Fatalf("expected timeout failure, got %+v", resp)
	}
	if resp.Message == "" {
		t.Fatal("expected error message")
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

func TestDownloadInstagramFastDL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil || r.PostForm.Get("q") != "https://www.instagram.com/p/Defd2kYAYK6/" || r.PostForm.Get("t") != "media" {
			t.Errorf("unexpected form %v", r.PostForm)
		}
		_, _ = w.Write([]byte(`{"status":"ok","p":"instagram","v":"v1","data":"<ul><li><div class=\"download-items\"><div class=\"download-items__btn\"><a href=\"https://dl.snapcdn.app/saveinsta?token=a\" title=\"Download Image\">Download Image</a></div></div></li></ul>"}`))
	}))
	defer srv.Close()

	orig := fastdlSearchURL
	fastdlSearchURL = srv.URL
	defer func() { fastdlSearchURL = orig }()

	resp := Download("https://www.instagram.com/p/Defd2kYAYK6/?img_index=2&exln=Znp3", &Options{Retry: 1})
	if !resp.Success || len(resp.Data.Media) != 1 || resp.Data.Media[0].Type != MediaImage {
		t.Fatalf("got %+v", resp)
	}
}

func TestDownloadInstagramFastDLMess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"ok","p":"instagram","mess":"Error: Video is private."}`))
	}))
	defer srv.Close()

	orig := fastdlSearchURL
	fastdlSearchURL = srv.URL
	defer func() { fastdlSearchURL = orig }()

	resp := Download("https://www.instagram.com/reel/C51YHfWJwHK/", &Options{Retry: 1})
	if resp.Success || resp.Message != "Error: Video is private." {
		t.Fatalf("got %+v", resp)
	}
}

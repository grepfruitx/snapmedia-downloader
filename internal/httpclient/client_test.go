package httpclient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestDoSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "test-agent" {
			t.Errorf("User-Agent = %q", r.Header.Get("User-Agent"))
		}
		if r.URL.Query().Get("lang") != "en" {
			t.Errorf("query lang = %q", r.URL.Query().Get("lang"))
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != "url=abc" {
			t.Errorf("body = %q", string(body))
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	got, err := Do(Request{
		Method: "POST",
		URL:    srv.URL,
		Query:  map[string]string{"lang": "en"},
		Body:   "url=abc",
	}, Options{Retry: 1, UserAgent: "test-agent"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ok" {
		t.Errorf("got %q, want %q", got, "ok")
	}
}

func TestDoRetriesUntilSuccess(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte("done"))
	}))
	defer srv.Close()

	got, err := Do(Request{URL: srv.URL}, Options{Retry: 3, RetryDelay: time.Millisecond})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "done" {
		t.Errorf("got %q, want %q", got, "done")
	}
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
}

func TestDoExhaustsRetries(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := Do(Request{URL: srv.URL}, Options{Retry: 2, RetryDelay: time.Millisecond})
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
}

func TestUnsupportedProxyScheme(t *testing.T) {
	_, err := Do(Request{URL: "http://example.com"}, Options{Retry: 1, Proxy: "ftp://proxy.example.com"})
	if err == nil {
		t.Fatal("expected error for unsupported proxy scheme")
	}
}

func TestDoTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = w.Write([]byte("late"))
	}))
	defer srv.Close()

	_, err := Do(Request{URL: srv.URL}, Options{Retry: 1, Timeout: 50 * time.Millisecond})
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(err.Error(), "Client.Timeout") && !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("expected timeout in error, got %v", err)
	}
}

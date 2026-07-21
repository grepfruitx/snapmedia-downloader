package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	snapsave "github.com/grepfruitx/snapmedia-downloader"
	"github.com/grepfruitx/snapmedia-downloader/internal/httpclient"
)

func main() {
	retry := flag.Int("retry", 1, "number of retry attempts")
	retryDelay := flag.Duration("retry-delay", 500*time.Millisecond, "delay between retries")
	userAgent := flag.String("user-agent", "", "custom User-Agent header")
	proxy := flag.String("proxy", "", "proxy URL (http, https, or socks5)")
	output := flag.String("o", "", "save media to this file or directory instead of printing JSON")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: snapmedia [flags] <url>")
		os.Exit(1)
	}

	opts := &snapsave.Options{Retry: *retry, RetryDelay: *retryDelay, UserAgent: *userAgent, Proxy: *proxy}
	resp := snapsave.Download(flag.Arg(0), opts)
	if !resp.Success {
		fmt.Fprintln(os.Stderr, "error:", resp.Message)
		os.Exit(1)
	}

	if *output == "" {
		b, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println(string(b))
		return
	}

	clientOpts := httpclient.Options{Retry: *retry, RetryDelay: *retryDelay, UserAgent: *userAgent, Proxy: *proxy}
	if err := saveMedia(resp.Data.Media, *output, clientOpts); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func saveMedia(media []snapsave.Media, output string, opts httpclient.Options) error {
	if len(media) == 0 {
		return fmt.Errorf("no media to save")
	}

	if len(media) == 1 {
		path := output
		if info, err := os.Stat(output); err == nil && info.IsDir() {
			path = filepath.Join(output, mediaFilename(media[0], 0))
		}
		if err := downloadTo(media[0].URL, path, opts); err != nil {
			return err
		}
		fmt.Println("saved", path)
		return nil
	}

	if err := os.MkdirAll(output, 0o755); err != nil {
		return err
	}
	for i, m := range media {
		path := filepath.Join(output, mediaFilename(m, i))
		if err := downloadTo(m.URL, path, opts); err != nil {
			return err
		}
		fmt.Println("saved", path)
	}
	return nil
}

func mediaFilename(m snapsave.Media, i int) string {
	ext := ""
	if u, err := url.Parse(m.URL); err == nil {
		ext = filepath.Ext(u.Path)
	}
	if ext == "" {
		ext = ".mp4"
		if m.Type == snapsave.MediaImage {
			ext = ".jpg"
		}
	}
	return fmt.Sprintf("media-%d%s", i, ext)
}

func downloadTo(mediaURL, path string, opts httpclient.Options) error {
	body, err := httpclient.Do(httpclient.Request{Method: "GET", URL: mediaURL}, opts)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(body), 0o644)
}

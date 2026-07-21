# snapmedia-downloader

Download Instagram, Facebook, TikTok and Twitter (X) media via the snapsave.app / snaptik.app services. Go port of [snapsave-media-downloader](https://github.com/ahmedrangel/snapsave-media-downloader).

## Install

```bash
go get github.com/grepfruitx/snapmedia-downloader
```

## Library usage

```go
package main

import (
	"fmt"

	snapsave "github.com/grepfruitx/snapmedia-downloader"
)

func main() {
	resp := snapsave.Download("https://www.instagram.com/reel/C51YHfWJwHK/", nil)
	if !resp.Success {
		fmt.Println("error:", resp.Message)
		return
	}
	for _, m := range resp.Data.Media {
		fmt.Println(m.Type, m.URL)
	}
}
```

`Download` never returns a Go `error` — failures are reported via `Response.Success` / `Response.Message`, matching the original library's response contract.

### Options

```go
resp := snapsave.Download(url, &snapsave.Options{
	Retry:      3,                      // default 1
	RetryDelay: 500 * time.Millisecond, // default 500ms
	UserAgent:  "custom-ua",            // default Chrome UA
	Proxy:      "http://user:pass@proxy.example.com:8080", // http, https, or socks5
})
```

## CLI usage

```bash
go install github.com/grepfruitx/snapmedia-downloader/cmd/snapmedia@latest
```

```bash
# print the JSON response
snapmedia "https://www.instagram.com/reel/C51YHfWJwHK/"

# save the media to disk instead (file, or directory when there are multiple items)
snapmedia -o ./downloads "https://www.facebook.com/watch?v=1234567890123456"
```

Flags:

| Flag | Default | Description |
|---|---|---|
| `--retry` | `1` | number of retry attempts |
| `--retry-delay` | `500ms` | delay between retries |
| `--user-agent` | Chrome UA | custom `User-Agent` header |
| `--proxy` | _(none)_ | proxy URL: `http://`, `https://`, or `socks5://` |
| `-o` | _(none)_ | save media here instead of printing JSON (file for a single item, directory for multiple) |

## Project layout

```
snapmedia-downloader/
├── snapsave.go          # public API: Download(url, opts)
├── types.go             # Media, Data, Response, Options
├── cmd/snapmedia/       # CLI
└── internal/
    ├── urlmatch/        # platform URL regexes + normalization
    ├── decrypt/         # snapsave.app / snaptik.app response decoding
    ├── parser/          # HTML extraction (goquery)
    └── httpclient/       # HTTP client with retry + proxy support
```

## Known limitations

- **TikTok** (`snaptik.app`) is currently broken: the site dropped the `token` input its scraping flow relies on. This affects the original JS library the same way — it's an upstream site change, not a bug in this port.
- These are third-party scraping services with no stability guarantees; selectors may break whenever the target sites change their markup.

## Development

```bash
go test ./...
golangci-lint run ./...
```

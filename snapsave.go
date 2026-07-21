package snapsave

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/grepfruitx/snapmedia-downloader/internal/decrypt"
	"github.com/grepfruitx/snapmedia-downloader/internal/httpclient"
	"github.com/grepfruitx/snapmedia-downloader/internal/parser"
	"github.com/grepfruitx/snapmedia-downloader/internal/urlmatch"
)

const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"

var (
	snapsaveActionURL = "https://snapsave.app/action.php"
	snaptikHomeURL    = "https://snaptik.app/"
	snaptikActionURL  = "https://snaptik.app/abc2.php"
	twitterHomeURL    = "https://twitterdownloader.snapsave.app/"
	twitterActionURL  = "https://twitterdownloader.snapsave.app/action.php"
)

func resolveOptions(opts *Options) Options {
	resolved := Options{Retry: 1, RetryDelay: 500 * time.Millisecond, UserAgent: defaultUserAgent}
	if opts == nil {
		return resolved
	}
	if opts.Retry > 0 {
		resolved.Retry = opts.Retry
	}
	if opts.RetryDelay > 0 {
		resolved.RetryDelay = opts.RetryDelay
	}
	if opts.UserAgent != "" {
		resolved.UserAgent = opts.UserAgent
	}
	resolved.Proxy = opts.Proxy
	return resolved
}

func toClientOpts(o Options) httpclient.Options {
	return httpclient.Options{Retry: o.Retry, RetryDelay: o.RetryDelay, UserAgent: o.UserAgent, Proxy: o.Proxy}
}

func toData(r *parser.Result) *Data {
	media := make([]Media, len(r.Media))
	for i, m := range r.Media {
		media[i] = Media{
			Resolution:   m.Resolution,
			ShouldRender: m.ShouldRender,
			Thumbnail:    m.Thumbnail,
			Type:         MediaType(m.Type),
			URL:          m.URL,
		}
	}
	return &Data{Description: r.Description, Preview: r.Preview, Media: media}
}

// Download fetches and extracts media for a Facebook, Instagram, TikTok or
// Twitter (X) URL. It never returns a Go error; failures are reported via
// Response.Success/Message, mirroring the original API's response contract.
func Download(rawURL string, opts *Options) *Response {
	o := resolveOptions(opts)

	if !urlmatch.IsValid(rawURL) {
		return &Response{Success: false, Message: "Invalid URL"}
	}

	form := url.Values{"url": {urlmatch.Normalize(rawURL)}}

	var data *Data
	var err error
	switch {
	case urlmatch.TikTok.MatchString(rawURL):
		data, err = downloadTikTok(form, o)
	case urlmatch.Twitter.MatchString(rawURL):
		data, err = downloadTwitter(form, o)
	default:
		data, err = downloadSnapSave(form, o)
	}
	if err != nil {
		return &Response{Success: false, Message: err.Error()}
	}
	return &Response{Success: true, Data: data}
}

func downloadSnapSave(form url.Values, o Options) (*Data, error) {
	body, err := httpclient.Do(httpclient.Request{
		Method: "POST",
		URL:    snapsaveActionURL,
		Query:  map[string]string{"lang": "en"},
		Headers: map[string]string{
			"accept":       "*/*",
			"content-type": "application/x-www-form-urlencoded",
			"origin":       "https://snapsave.app",
			"referer":      "https://snapsave.app/",
		},
		Body: form.Encode(),
	}, toClientOpts(o))
	if err != nil {
		return nil, err
	}

	decoded, err := decrypt.DecryptSnapSave(body)
	if err != nil {
		return nil, err
	}
	result, err := parser.ParseSnapSave(decoded)
	if err != nil {
		return nil, err
	}
	return toData(result), nil
}

func downloadTikTok(form url.Values, o Options) (*Data, error) {
	home, err := httpclient.Do(httpclient.Request{
		Method: "GET",
		URL:    snaptikHomeURL,
	}, toClientOpts(o))
	if err != nil {
		return nil, err
	}
	token, err := parser.ExtractInputValue(home, "token")
	if err != nil {
		return nil, err
	}
	form.Set("token", token)

	body, err := httpclient.Do(httpclient.Request{
		Method: "POST",
		URL:    snaptikActionURL,
		Headers: map[string]string{
			"accept":       "*/*",
			"content-type": "application/x-www-form-urlencoded",
			"origin":       "https://snaptik.app",
			"referer":      "https://snaptik.app/",
		},
		Body: form.Encode(),
	}, toClientOpts(o))
	if err != nil {
		return nil, err
	}

	decoded, err := decrypt.DecryptSnaptik(body)
	if err != nil {
		return nil, err
	}
	result, err := parser.ParseSnaptik(decoded)
	if err != nil {
		return nil, err
	}
	return toData(result), nil
}

type twitterActionResponse struct {
	Data string `json:"data"`
}

func downloadTwitter(form url.Values, o Options) (*Data, error) {
	home, err := httpclient.Do(httpclient.Request{
		Method: "GET",
		URL:    twitterHomeURL,
	}, toClientOpts(o))
	if err != nil {
		return nil, err
	}
	token, err := parser.ExtractInputValue(home, "token")
	if err != nil {
		return nil, err
	}
	form.Set("token", token)

	body, err := httpclient.Do(httpclient.Request{
		Method: "POST",
		URL:    twitterActionURL,
		Headers: map[string]string{
			"accept":       "*/*",
			"content-type": "application/x-www-form-urlencoded",
			"origin":       "https://twitterdownloader.snapsave.app",
			"referer":      "https://twitterdownloader.snapsave.app/",
		},
		Body: form.Encode(),
	}, toClientOpts(o))
	if err != nil {
		return nil, err
	}

	var action twitterActionResponse
	if err := json.Unmarshal([]byte(body), &action); err != nil {
		return nil, err
	}
	result, err := parser.ParseTwitter(action.Data)
	if err != nil {
		return nil, err
	}
	return toData(result), nil
}

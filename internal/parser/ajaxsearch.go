package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ajaxSearchResponse is the JSON envelope returned by the /api/ajaxSearch
// endpoint shared by fastdl.to (Instagram) and savef.app (Facebook).
type ajaxSearchResponse struct {
	Status string `json:"status"`
	Mess   string `json:"mess"`
	Data   string `json:"data"`
}

// ParseAjaxSearch decodes the ajaxSearch envelope and returns the embedded
// HTML. Service-side failures (deleted/private posts etc.) arrive with
// status "ok" and a non-empty "mess"; they are turned into errors.
func ParseAjaxSearch(body string) (string, error) {
	var r ajaxSearchResponse
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		return "", fmt.Errorf("parser: ajaxSearch: %w", err)
	}
	if r.Mess != "" {
		return "", errors.New(htmlText(r.Mess))
	}
	if r.Status != "ok" {
		return "", fmt.Errorf("parser: ajaxSearch: unexpected status %q", r.Status)
	}
	if r.Data == "" {
		return "", errors.New("parser: blank data")
	}
	return r.Data, nil
}

// ParseFastDL parses fastdl.to result HTML. Each div.download-items is one
// post/carousel/story item; video items also carry a "Download Thumbnail"
// button (div.dl-thumb) that must be skipped.
func ParseFastDL(html string) (*Result, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	var media []Media
	doc.Find("div.download-items").Each(func(_ int, item *goquery.Selection) {
		a := item.Find("div.download-items__btn").Not(".dl-thumb").Find("a").First()
		href, _ := a.Attr("href")
		if href == "" {
			return
		}
		title, _ := a.Attr("title")
		if title == "" {
			title = strings.TrimSpace(a.Text())
		}
		img := item.Find("div.download-items__thumb > img")
		thumbnail, ok := img.Attr("data-src")
		if !ok {
			thumbnail, _ = img.Attr("src")
		}
		m := Media{URL: href, Type: "video"}
		if strings.Contains(title, "Image") || strings.Contains(title, "Photo") {
			m.Type = "image"
		} else {
			m.Thumbnail = thumbnail
		}
		media = append(media, m)
	})
	if len(media) == 0 {
		return nil, errors.New("parser: blank data")
	}
	return &Result{Media: media}, nil
}

// ParseSaveF parses savef.app result HTML. Only the first table (MP4) is
// used; rows that require server-side rendering have no direct link and are
// skipped. Rows are ordered best quality first, so only the first usable row
// is returned — callers send every returned video to the user.
func ParseSaveF(html string) (*Result, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	result := &Result{}
	result.Preview, _ = doc.Find("div.thumbnail img").First().Attr("src")
	doc.Find("table").First().Find("tbody > tr").EachWithBreak(func(_ int, row *goquery.Selection) bool {
		href, _ := row.Find("a[href]").Attr("href")
		if !strings.HasPrefix(href, "http") {
			return true
		}
		result.Media = append(result.Media, Media{
			Resolution: strings.TrimSpace(row.Find("td.video-quality").Text()),
			Thumbnail:  result.Preview,
			Type:       "video",
			URL:        href,
		})
		return false
	})
	if len(result.Media) == 0 {
		return nil, errors.New("parser: blank data")
	}
	return result, nil
}

func htmlText(s string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return s
	}
	return strings.TrimSpace(doc.Text())
}

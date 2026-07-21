package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseTwitter(html string) (*Result, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	mediaURL, _ := doc.Find("#download-block > .abuttons > a").Attr("href")
	description := strings.TrimSpace(doc.Find(".videotikmate-middle > p > span").Text())
	preview, _ := doc.Find(".videotikmate-left > img").Attr("src")
	spanText := strings.TrimSpace(doc.Find("#download-block > .abuttons > a > span > span").Text())

	mediaType := "video"
	if spanText == "Download photo" {
		mediaType = "image"
	}

	return &Result{
		Description: description,
		Preview:     preview,
		Media:       []Media{{URL: mediaURL, Type: mediaType}},
	}, nil
}

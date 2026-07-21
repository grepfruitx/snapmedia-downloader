package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseSnaptik(html string) (*Result, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	mediaURL, _ := doc.Find(".download-box > .video-links > a").Attr("href")
	description := strings.TrimSpace(doc.Find(".video-title").Text())
	preview, _ := doc.Find("#thumbnail").Attr("src")
	spanText := strings.TrimSpace(doc.Find(".video-links > a.button.download-file").Text())

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

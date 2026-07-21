package parser

import (
	"errors"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	isProgressAPI  = regexp.MustCompile(`(?i)get_progressApi`)
	progressAPIURL = regexp.MustCompile(`get_progressApi\('(.*?)'\)`)
)

func hrefOrOnclick(s *goquery.Selection) string {
	if v, ok := s.Find("a").Attr("href"); ok {
		return v
	}
	v, _ := s.Find("button").Attr("onclick")
	return v
}

func ParseSnapSave(html string) (*Result, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	result := &Result{}
	var media []Media

	switch {
	case doc.Find("table.table").Length() > 0 || doc.Find("article.media > figure").Length() > 0:
		result.Description = strings.TrimSpace(doc.Find("span.video-des").Text())
		result.Preview, _ = doc.Find("article.media > figure").Find("img").Attr("src")

		switch {
		case doc.Find("table.table").Length() > 0:
			doc.Find("tbody > tr").Each(func(_ int, row *goquery.Selection) {
				td := row.Find("td")
				resolution := td.Eq(0).Text()
				mediaURL := hrefOrOnclick(td.Eq(2))
				shouldRender := isProgressAPI.MatchString(mediaURL)
				if shouldRender {
					if m := progressAPIURL.FindStringSubmatch(mediaURL); m != nil {
						mediaURL = "https://snapsave.app" + m[1]
					}
				}
				mediaType := "video"
				if resolution == "" {
					mediaType = "image"
				}
				media = append(media, Media{
					Resolution:   resolution,
					ShouldRender: shouldRender,
					URL:          mediaURL,
					Type:         mediaType,
				})
			})
		case doc.Find("div.card").Length() > 0:
			doc.Find("div.card").Each(func(_ int, card *goquery.Selection) {
				body := card.Find("div.card-body")
				aText := strings.TrimSpace(body.Find("a").Text())
				mediaURL, _ := body.Find("a").Attr("href")
				mediaType := "video"
				if aText == "Download Photo" {
					mediaType = "image"
				}
				media = append(media, Media{URL: mediaURL, Type: mediaType})
			})
		default:
			mediaURL := hrefOrOnclick(doc.Selection)
			aText := strings.TrimSpace(doc.Find("a").Text())
			mediaType := "video"
			if aText == "Download Photo" {
				mediaType = "image"
			}
			media = append(media, Media{URL: mediaURL, Type: mediaType})
		}

	case doc.Find("div.download-items").Length() > 0:
		doc.Find("div.download-items").Each(func(_ int, item *goquery.Selection) {
			thumbnail, _ := item.Find("div.download-items__thumb > img").Attr("src")
			btn := item.Find("div.download-items__btn")
			downloadURL, _ := btn.Find("a").Attr("href")
			spanText := strings.TrimSpace(btn.Find("span").Text())
			mediaType := "video"
			if spanText == "Download Photo" {
				mediaType = "image"
			}
			mediaURL := downloadURL
			m := Media{Type: mediaType}
			if mediaType == "image" {
				mediaURL = thumbnail
			} else if thumbnail != "" {
				m.Thumbnail = FixThumbnail(thumbnail)
			}
			m.URL = mediaURL
			media = append(media, m)
		})
	}

	if len(media) == 0 {
		return nil, errors.New("parser: blank data")
	}
	result.Media = media
	return result, nil
}

package parser

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractInputValue(html, name string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	v, ok := doc.Find(fmt.Sprintf(`input[name='%s']`, name)).Attr("value")
	if !ok {
		return "", fmt.Errorf("parser: input[name=%s] not found", name)
	}
	return v, nil
}

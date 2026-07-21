package parser

import "testing"

func TestParseTwitter(t *testing.T) {
	html := `
	<div class="videotikmate-left"><img src="https://example.com/preview.jpg"></div>
	<div class="videotikmate-middle"><p><span>X description example</span></p></div>
	<div id="download-block"><div class="abuttons">
		<a href="https://example.com/video.mp4"><span><span>Download video</span></span></a>
	</div></div>`

	result, err := ParseTwitter(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Description != "X description example" {
		t.Errorf("description = %q", result.Description)
	}
	if result.Preview != "https://example.com/preview.jpg" {
		t.Errorf("preview = %q", result.Preview)
	}
	if len(result.Media) != 1 || result.Media[0].Type != "video" || result.Media[0].URL != "https://example.com/video.mp4" {
		t.Errorf("media = %+v", result.Media)
	}
}

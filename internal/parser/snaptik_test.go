package parser

import "testing"

func TestParseSnaptik(t *testing.T) {
	html := `
	<div class="video-title">Cat video</div>
	<img id="thumbnail" src="https://example.com/thumb.jpg">
	<div class="download-box"><div class="video-links">
		<a href="https://example.com/video.mp4"></a>
		<a class="button download-file">Download video</a>
	</div></div>`

	result, err := ParseSnaptik(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Description != "Cat video" {
		t.Errorf("description = %q", result.Description)
	}
	if result.Preview != "https://example.com/thumb.jpg" {
		t.Errorf("preview = %q", result.Preview)
	}
	if len(result.Media) != 1 || result.Media[0].Type != "video" || result.Media[0].URL != "https://example.com/video.mp4" {
		t.Errorf("media = %+v", result.Media)
	}
}

func TestParseSnaptikPhoto(t *testing.T) {
	html := `
	<div class="download-box"><div class="video-links">
		<a href="https://example.com/photo.jpg"></a>
		<a class="button download-file">Download photo</a>
	</div></div>`

	result, err := ParseSnaptik(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Media) != 1 || result.Media[0].Type != "image" {
		t.Errorf("media = %+v", result.Media)
	}
}

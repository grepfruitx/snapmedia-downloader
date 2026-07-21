package parser

import "testing"

func TestParseSnapSaveTable(t *testing.T) {
	html := `
	<article class="media"><figure><img src="https://example.com/preview.jpg"></figure></article>
	<span class="video-des">This is a description</span>
	<table class="table"><tbody>
		<tr><td>720p (HD)</td><td></td><td><a href="https://example.com/video-720.mp4"></a></td></tr>
		<tr><td>360p (SD)</td><td></td><td><button onclick="get_progressApi('/render.php?token=abc')"></button></td></tr>
	</tbody></table>`

	result, err := ParseSnapSave(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Description != "This is a description" {
		t.Errorf("description = %q", result.Description)
	}
	if result.Preview != "https://example.com/preview.jpg" {
		t.Errorf("preview = %q", result.Preview)
	}
	if len(result.Media) != 2 {
		t.Fatalf("media count = %d, want 2", len(result.Media))
	}
	if result.Media[0].Resolution != "720p (HD)" || result.Media[0].Type != "video" || result.Media[0].URL != "https://example.com/video-720.mp4" {
		t.Errorf("media[0] = %+v", result.Media[0])
	}
	if !result.Media[1].ShouldRender || result.Media[1].URL != "https://snapsave.app/render.php?token=abc" {
		t.Errorf("media[1] = %+v", result.Media[1])
	}
}

func TestParseSnapSaveCard(t *testing.T) {
	html := `
	<article class="media"><figure><img src="https://example.com/preview.jpg"></figure></article>
	<span class="video-des">desc</span>
	<div class="card"><div class="card-body"><a href="https://example.com/photo.jpg">Download Photo</a></div></div>`

	result, err := ParseSnapSave(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Media) != 1 || result.Media[0].Type != "image" || result.Media[0].URL != "https://example.com/photo.jpg" {
		t.Errorf("media = %+v", result.Media)
	}
}

func TestParseSnapSaveSingleLink(t *testing.T) {
	html := `
	<article class="media"><figure><img src="https://example.com/preview.jpg"></figure></article>
	<span class="video-des">desc</span>
	<a href="https://example.com/video.mp4">Download Video</a>`

	result, err := ParseSnapSave(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Media) != 1 || result.Media[0].Type != "video" || result.Media[0].URL != "https://example.com/video.mp4" {
		t.Errorf("media = %+v", result.Media)
	}
}

func TestParseSnapSaveDownloadItems(t *testing.T) {
	html := `
	<div class="download-items">
		<div class="download-items__thumb"><img src="https://snapinsta.app/photo.php?photo=https%3A%2F%2Fexample.com%2Fthumb.jpg"></div>
		<div class="download-items__btn"><a href="https://example.com/video.mp4"></a><span>Download Video</span></div>
	</div>`

	result, err := ParseSnapSave(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Media) != 1 {
		t.Fatalf("media count = %d, want 1", len(result.Media))
	}
	m := result.Media[0]
	if m.Type != "video" || m.URL != "https://example.com/video.mp4" || m.Thumbnail != "https://example.com/thumb.jpg" {
		t.Errorf("media[0] = %+v", m)
	}
}

func TestParseSnapSaveBlank(t *testing.T) {
	if _, err := ParseSnapSave(`<html></html>`); err == nil {
		t.Fatal("expected error for blank data")
	}
}

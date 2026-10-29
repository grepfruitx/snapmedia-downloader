package parser

import "testing"

func TestParseAjaxSearchOK(t *testing.T) {
	html, err := ParseAjaxSearch(`{"status":"ok","p":"instagram","v":"v1","data":"<ul></ul>"}`)
	if err != nil || html != "<ul></ul>" {
		t.Fatalf("got %q, %v", html, err)
	}
}

func TestParseAjaxSearchMess(t *testing.T) {
	_, err := ParseAjaxSearch(`{"status":"ok","p":"instagram","mess":"Error: Video is private. Please use the tool <a href=\"/en/private\"><b>Instagram Private Downloader</b></a>"}`)
	if err == nil || err.Error() != "Error: Video is private. Please use the tool Instagram Private Downloader" {
		t.Fatalf("got %v", err)
	}
}

func TestParseAjaxSearchBadJSON(t *testing.T) {
	if _, err := ParseAjaxSearch(`<html>`); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseFastDL(t *testing.T) {
	html := `<ul class="download-box">
<li><div class="download-items">
  <div class="download-items__thumb"><img src="https://i.snapcdn.app/photo?token=img1"></div>
  <div class="photo-option"><select><option value="https://dl.snapcdn.app/saveinsta?token=full1">3024x4032</option></select></div>
  <div class="download-items__btn"><a href="https://dl.snapcdn.app/saveinsta?token=full1" title="Download Image"><span><span>Download Image</span></span></a></div>
</div></li>
<li><div class="download-items">
  <div class="download-items__thumb"><img class="lazy" src="/imgs/loader.gif" data-src="https://i.snapcdn.app/photo?token=thumb2"></div>
  <div class="download-items__btn dl-thumb"><a href="https://dl.snapcdn.app/saveinsta?token=thumbdl2" title="Download Thumbnail"><span>Download Thumbnail</span></a></div>
  <div class="download-items__btn"><a href="https://dl.snapcdn.app/saveinsta?token=video2" title="Download Video"><span>Download Video</span></a></div>
</div></li>
</ul>`
	r, err := ParseFastDL(html)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Media) != 2 {
		t.Fatalf("expected 2 media, got %+v", r.Media)
	}
	if m := r.Media[0]; m.Type != "image" || m.URL != "https://dl.snapcdn.app/saveinsta?token=full1" {
		t.Errorf("image: %+v", m)
	}
	if m := r.Media[1]; m.Type != "video" || m.URL != "https://dl.snapcdn.app/saveinsta?token=video2" || m.Thumbnail != "https://i.snapcdn.app/photo?token=thumb2" {
		t.Errorf("video: %+v", m)
	}
}

func TestParseFastDLBlank(t *testing.T) {
	if _, err := ParseFastDL(`<ul class="download-box"></ul>`); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseSaveF(t *testing.T) {
	html := `<div class="detail"><div class="thumbnail"><div class="image-fb"><img src="https://scontent.fbcdn.net/thumb.jpg"></div></div></div>
<table class="table is-fullwidth"><tbody>
<tr><td class="video-quality">720p (HD)</td><td>No</td><td><a href="https://dl.snapcdn.app/download?token=hd" class="download-link-fb">Download</a></td></tr>
<tr><td class="video-quality">360p (SD)</td><td>No</td><td><a href="https://dl.snapcdn.app/download?token=sd" class="download-link-fb">Download</a></td></tr>
<tr><td class="video-quality">540p</td><td>Yes</td><td><button onclick="convertFile(this, 'mp4')" data-videourl="https://video.fbcdn.net/x.mp4">Render</button></td></tr>
</tbody></table>
<table class="table is-fullwidth"><tbody>
<tr><td>320kbps</td><td>Yes</td><td><a onclick="convertFile(this, 'mp3')" class="button">Render</a></td></tr>
</tbody></table>`
	r, err := ParseSaveF(html)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Media) != 1 {
		t.Fatalf("expected only the best quality row, got %+v", r.Media)
	}
	m := r.Media[0]
	if m.Type != "video" || m.Resolution != "720p (HD)" || m.URL != "https://dl.snapcdn.app/download?token=hd" || m.Thumbnail != "https://scontent.fbcdn.net/thumb.jpg" {
		t.Errorf("got %+v", m)
	}
	if r.Preview != "https://scontent.fbcdn.net/thumb.jpg" {
		t.Errorf("preview %q", r.Preview)
	}
}

func TestParseSaveFRenderOnly(t *testing.T) {
	html := `<table class="table"><tbody><tr><td class="video-quality">540p</td><td>Yes</td><td><button onclick="convertFile(this, 'mp4')">Render</button></td></tr></tbody></table>`
	if _, err := ParseSaveF(html); err == nil {
		t.Fatal("expected error when no direct link")
	}
}

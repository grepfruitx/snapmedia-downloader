package urlmatch

import "testing"

func TestFacebook(t *testing.T) {
	urls := []string{
		"https://facebook.com/watch?v=1234567890123456",
		"https://www.facebook.com/watch?v=1234567890123456",
		"https://www.facebook.com/watch/?v=1234567890123456",
		"https://www.facebook.com/reel/1234567890123456",
		"https://www.facebook.com/reel/1234567890123456?locale=es_LA",
		"https://web.facebook.com/reel/1234567890123456",
		"https://www.facebook.com/user_name/videos/1234567890123456",
		"https://www.facebook.com/user_name/posts/1234567890123456",
		"https://www.facebook.com/1234567890123456/videos/1234567890123456",
		"https://www.facebook.com/1234567890123456/posts/1234567890123456",
		"https://www.facebook.com/Page_Name/videos/1234567890123456",
		"https://www.facebook.com/Page_Name/posts/1234567890123456",
		"https://fb.watch/abCDefghI9",
		"https://www.facebook.com/share/v/ab0cdEfGHIJKMnlO/",
		"https://www.facebook.com/share/r/ab0cdEfGHIJKMnlO/",
	}
	for _, u := range urls {
		if !Facebook.MatchString(u) {
			t.Errorf("expected match: %s", u)
		}
	}
}

func TestInstagram(t *testing.T) {
	urls := []string{
		"https://instagram.com/p/CcDeFg9hiJK",
		"https://www.instagram.com/p/CcDeFg9hiJK",
		"https://www.instagram.com/p/CcDeFg9hiJK/",
		"https://www.instagram.com/reel/1234567890123456",
		"https://www.instagram.com/reel/1234567890123456/",
		"https://www.instagram.com/reels/1234567890123456",
		"https://www.instagram.com/reels/1234567890123456/",
		"https://www.instagram.com/tv/CcDeFg9hiJK",
		"https://www.instagram.com/tv/CcDeFg9hiJK/",
		"https://www.instagram.com/stories/user_name",
		"https://www.instagram.com/stories/user_name/",
		"https://www.instagram.com/stories/user_name/1234567890123456",
		"https://www.instagram.com/stories/user_name/1234567890123456/",
		"https://instagram.com/user_name/p/CcDeFg9hiJK",
		"https://www.instagram.com/user_name/reel/1234567890123456/",
	}
	for _, u := range urls {
		if !Instagram.MatchString(u) {
			t.Errorf("expected match: %s", u)
		}
	}
}

func TestTikTok(t *testing.T) {
	urls := []string{
		"https://tiktok.com/@user_name/video/1234567890123456789",
		"https://www.tiktok.com/@user_name/video/1234567890123456789/",
		"https://www.tiktok.com/@user_name/video/1234567890123456789?is_from_webapp=1&sender_device=pc",
		"https://www.tiktok.com/@user_name/photo/1234567890123456789/",
		"https://www.tiktok.com/@user_name/photo/1234567890123456789?is_from_webapp=1&sender_device=pc",
		"https://www.tiktok.com/t/ZTFErV3e5",
		"https://www.tiktok.com/t/ZTFErV3e5/",
		"https://vm.tiktok.com/ZAbc8d911",
		"https://vm.tiktok.com/ZAbc8d911/",
		"https://m.tiktok.com/v/1234567890123456789",
		"https://m.tiktok.com/v/1234567890123456789/",
		"https://vt.tiktok.com/ZAbc8d911",
		"https://vt.tiktok.com/ZAbc8d911/",
	}
	for _, u := range urls {
		if !TikTok.MatchString(u) {
			t.Errorf("expected match: %s", u)
		}
	}
}

func TestTwitter(t *testing.T) {
	urls := []string{
		"https://twitter.com/user_name/status/1234567890123456789",
		"https://x.com/user_name/status/1234567890123456789",
		"https://twitter.com/i/web/status/1234567890123456789",
		"https://x.com/i/web/status/1234567890123456789",
		"https://twitter.com/user_name/status/1234567890123456789/photo/1",
		"https://x.com/user_name/status/1234567890123456789/photo/1",
	}
	for _, u := range urls {
		if !Twitter.MatchString(u) {
			t.Errorf("expected match: %s", u)
		}
	}
}

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"https://tiktok.com":     "https://www.tiktok.com",
		"https://www.tiktok.com": "https://www.tiktok.com",
		"https://vm.tiktok.com":  "https://vm.tiktok.com",
		"https://x.com":          "https://x.com",
	}
	for in, want := range cases {
		if got := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

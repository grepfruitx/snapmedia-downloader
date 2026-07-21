package urlmatch

import "regexp"

var (
	Facebook  = regexp.MustCompile(`^https?://(?:www\.|web\.|m\.)?facebook\.com/(watch(\?v=|/\?v=)[0-9]+|reel/[0-9]+|[a-zA-Z0-9.\-_]+/(videos|posts)/[0-9]+|[0-9]+/(videos|posts)/[0-9]+|[a-zA-Z0-9]+/(videos|posts)/[0-9]+|share/(v|r)/[a-zA-Z0-9]+/?)([^/?#&]+).*$|^https://fb\.watch/[a-zA-Z0-9]+$`)
	Instagram = regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/(?:[^/]+/)?(?:p|reel|reels|tv|stories|share)/([^/?#&]+).*`)
	TikTok    = regexp.MustCompile(`^https?://(?:www\.|m\.|vm\.|vt\.)?tiktok\.com/(?:@[^/]+/(?:video|photo)/\d+|v/\d+|t/[\w]+|[\w]+)/?`)
	Twitter   = regexp.MustCompile(`^https://(?:x|twitter)\.com(?:/(?:i/web|[^/]+)/status/(\d+)(?:.*)?)?$`)
)

func IsValid(url string) bool {
	return Facebook.MatchString(url) || Instagram.MatchString(url) || Twitter.MatchString(url) || TikTok.MatchString(url)
}

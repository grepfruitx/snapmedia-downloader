package urlmatch

import (
	"regexp"
	"strings"
)

var (
	schemeAlnum = regexp.MustCompile(`(?i)^(https?://)[a-z0-9]`)
	twoLabel    = regexp.MustCompile(`(?i)^(https?://)([^./]+\.[^./]+)(/.*)?$`)
)

func Normalize(u string) string {
	if Twitter.MatchString(u) {
		return u
	}
	if Instagram.MatchString(u) {
		if i := strings.IndexByte(u, '?'); i != -1 {
			u = u[:i]
		}
	}
	m := schemeAlnum.FindStringSubmatch(u)
	if m == nil {
		return u
	}
	if strings.HasPrefix(strings.ToLower(u[len(m[1]):]), "www.") {
		return u
	}
	dm := twoLabel.FindStringSubmatch(u)
	if dm == nil {
		return u
	}
	return dm[1] + "www." + dm[2] + dm[3]
}

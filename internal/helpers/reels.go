package helpers

import (
	"fmt"
	"net/url"
	"regexp"
)

const (
	kkInstHost = "www.kkinstagram.com"
	reelPath   = "/reel/%s/"
)

var reelRegExp = regexp.MustCompile(`https?:\/\/(www\.)?instagram\.com\/reel\/([a-zA-Z0-9_-]+)`)

func ExtractReelId(text string) (string, bool) {
	match := reelRegExp.FindStringSubmatch(text)

	if len(match) < 3 {
		return "", false
	}

	return match[2], true
}

func BuildKkInstagramUrl(reelId string) (string, error) {
	u := url.URL{
		Scheme: "https",
		Host:   kkInstHost,
		Path:   fmt.Sprintf(reelPath, reelId),
	}

	return u.String(), nil
}

package helpers

import (
	"net/url"
	"regexp"
)

const (
	kkInstHost = "www.kkclip.com"
	reelPath   = "/reel/%s/"
)

var instagramRegExp = regexp.MustCompile(`https?:\/\/(www\.)?instagram\.com\/[^\s]+`)

func ExtractReelId(text string) (string, bool) {
	match := instagramRegExp.FindString(text)

	if match == "" {
		return "", false
	}

	return match, true
}

func BuildKkInstagramUrl(instagramUrl string) (string, error) {
	u, err := url.Parse(instagramUrl)
	if err != nil {
		return "", err
	}

	u.Host = kkInstHost

	return u.String(), nil
}

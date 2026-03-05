package helpers

import (
	"net/url"
	"regexp"
)

const (
	tfTiktokHost      = "tfxktok.com"
	tiktokVmShortLink = "vm.tiktok.com"
	tiktokVtShortLink = "vt.tiktok.com"
)

var tiktokRegex = regexp.MustCompile(`https?:\/\/(vm|vt|www|m)?\.?tiktok\.com\/[^\s]+`)

func IsShortUrl(raw string) bool {
	u, err := url.Parse(raw)

	if err != nil {
		return false
	}

	switch u.Host {
	case tiktokVmShortLink, tiktokVtShortLink:
		return true
	}

	return false
}

func ConverTikTokToVxUrl(originalUrl string) (string, error) {
	u, err := url.Parse(originalUrl)
	if err != nil {
		return "", err
	}

	u.Host = tfTiktokHost

	return u.String(), nil
}

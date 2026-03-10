package utilities

import (
	"fmt"
	"net/url"
	"strings"
)

// Platform represents a supported media platform.
type Platform string

const (
	PlatformYouTube   Platform = "youtube"
	PlatformTwitter   Platform = "twitter"
	PlatformInstagram Platform = "instagram"
	PlatformTikTok    Platform = "tiktok"
	PlatformUnknown   Platform = "unknown"
)

// platformHosts maps known hostnames (without "www.") to their Platform.
var platformHosts = map[string]Platform{
	"youtube.com":    PlatformYouTube,
	"youtu.be":       PlatformYouTube,
	"m.youtube.com":  PlatformYouTube,
	"twitter.com":    PlatformTwitter,
	"x.com":          PlatformTwitter,
	"mobile.x.com":   PlatformTwitter,
	"instagram.com":  PlatformInstagram,
	"m.instagram.com": PlatformInstagram,
	"tiktok.com":     PlatformTikTok,
	"m.tiktok.com":   PlatformTikTok,
	"vm.tiktok.com":  PlatformTikTok,
	"vt.tiktok.com":  PlatformTikTok,
}

// DetectPlatform takes a raw URL string and returns which Platform it belongs to.
// Returns PlatformUnknown if the URL is invalid or the host is not recognised.
func DetectPlatform(rawURL string) Platform {
	parsed, err := parseRawURL(rawURL)
	if err != nil {
		return PlatformUnknown
	}

	host := normaliseHost(parsed.Hostname())

	if p, ok := platformHosts[host]; ok {
		return p
	}

	return PlatformUnknown
}

// normaliseHost strips the "www." prefix and lowercases the hostname.
func normaliseHost(host string) string {
	host = strings.ToLower(host)
	host = strings.TrimPrefix(host, "www.")
	return host
}

// parseRawURL adds a scheme when missing, then parses the URL.
func parseRawURL(rawURL string) (*url.URL, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil, fmt.Errorf("empty URL")
	}

	// Add scheme if the user didn't provide one.
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}

	return url.Parse(rawURL)
}

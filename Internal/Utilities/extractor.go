package utilities

import (
	"fmt"
	"regexp"
	"strings"
)

// ---------- YouTube ----------

// youtubeIDRegex matches an 11-character YouTube video ID (letters, digits, - and _).
var youtubeIDRegex = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

// ExtractYouTubeID returns the video ID from a YouTube URL.
//
// Supported formats:
//   - https://www.youtube.com/watch?v=VIDEO_ID
//   - https://youtu.be/VIDEO_ID
//   - https://www.youtube.com/embed/VIDEO_ID
//   - https://www.youtube.com/shorts/VIDEO_ID
//   - https://www.youtube.com/live/VIDEO_ID
//   - https://m.youtube.com/watch?v=VIDEO_ID
func ExtractYouTubeID(rawURL string) (string, error) {
	parsed, err := parseRawURL(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	host := normaliseHost(parsed.Hostname())

	switch host {
	case "youtu.be":
		// https://youtu.be/VIDEO_ID
		id := strings.TrimPrefix(parsed.Path, "/")
		id = strings.SplitN(id, "/", 2)[0] // strip anything after the id
		if !youtubeIDRegex.MatchString(id) {
			return "", fmt.Errorf("invalid YouTube video ID in short URL: %q", id)
		}
		return id, nil

	case "youtube.com", "m.youtube.com":
		// /watch?v=VIDEO_ID
		if v := parsed.Query().Get("v"); v != "" {
			if !youtubeIDRegex.MatchString(v) {
				return "", fmt.Errorf("invalid YouTube video ID: %q", v)
			}
			return v, nil
		}

		// /embed/VIDEO_ID, /shorts/VIDEO_ID, /live/VIDEO_ID
		for _, prefix := range []string{"/embed/", "/shorts/", "/live/", "/v/"} {
			if strings.HasPrefix(parsed.Path, prefix) {
				id := strings.TrimPrefix(parsed.Path, prefix)
				id = strings.SplitN(id, "/", 2)[0]
				if !youtubeIDRegex.MatchString(id) {
					return "", fmt.Errorf("invalid YouTube video ID in path: %q", id)
				}
				return id, nil
			}
		}

		return "", fmt.Errorf("could not find video ID in YouTube URL: %s", rawURL)

	default:
		return "", fmt.Errorf("not a YouTube URL: %s", rawURL)
	}
}

// ---------- Twitter / X ----------

// twitterStatusRegex captures the numeric status (tweet) ID from the path.
var twitterStatusRegex = regexp.MustCompile(`^/[^/]+/status/(\d+)`)

// ExtractTwitterID returns the tweet/post ID from a Twitter/X URL.
//
// Supported formats:
//   - https://twitter.com/user/status/1234567890
//   - https://x.com/user/status/1234567890
//   - https://mobile.x.com/user/status/1234567890
func ExtractTwitterID(rawURL string) (string, error) {
	parsed, err := parseRawURL(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	host := normaliseHost(parsed.Hostname())
	if host != "twitter.com" && host != "x.com" && host != "mobile.x.com" {
		return "", fmt.Errorf("not a Twitter/X URL: %s", rawURL)
	}

	matches := twitterStatusRegex.FindStringSubmatch(parsed.Path)
	if matches == nil || len(matches) < 2 {
		return "", fmt.Errorf("could not find tweet ID in URL: %s", rawURL)
	}

	return matches[1], nil
}

// ---------- Instagram ----------

// instagramPostRegex captures the shortcode from Instagram post/reel paths.
var instagramPostRegex = regexp.MustCompile(`^/(p|reel|reels|tv)/([A-Za-z0-9_-]+)`)

// ExtractInstagramID returns the shortcode from an Instagram URL.
//
// Supported formats:
//   - https://www.instagram.com/p/SHORTCODE/
//   - https://www.instagram.com/reel/SHORTCODE/
//   - https://www.instagram.com/reels/SHORTCODE/
//   - https://www.instagram.com/tv/SHORTCODE/
func ExtractInstagramID(rawURL string) (string, error) {
	parsed, err := parseRawURL(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	host := normaliseHost(parsed.Hostname())
	if host != "instagram.com" && host != "m.instagram.com" {
		return "", fmt.Errorf("not an Instagram URL: %s", rawURL)
	}

	matches := instagramPostRegex.FindStringSubmatch(parsed.Path)
	if matches == nil || len(matches) < 3 {
		return "", fmt.Errorf("could not find post shortcode in Instagram URL: %s", rawURL)
	}

	return matches[2], nil
}

// ---------- TikTok ----------

// tiktokVideoRegex captures the numeric video ID from a TikTok standard URL.
var tiktokVideoRegex = regexp.MustCompile(`^/@[^/]+/video/(\d+)`)

// tiktokShortRegex captures the alphanumeric short-code from a vm.tiktok.com URL.
var tiktokShortRegex = regexp.MustCompile(`^/([A-Za-z0-9]+)`)

// ExtractTikTokID returns the video ID (or short code) from a TikTok URL.
//
// Supported formats:
//   - https://www.tiktok.com/@user/video/1234567890123456789
//   - https://vm.tiktok.com/ZMxxxxxx/
//   - https://vt.tiktok.com/ZMxxxxxx/
//   - https://m.tiktok.com/@user/video/1234567890123456789
func ExtractTikTokID(rawURL string) (string, error) {
	parsed, err := parseRawURL(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	host := normaliseHost(parsed.Hostname())

	switch host {
	case "tiktok.com", "m.tiktok.com":
		matches := tiktokVideoRegex.FindStringSubmatch(parsed.Path)
		if matches == nil || len(matches) < 2 {
			return "", fmt.Errorf("could not find video ID in TikTok URL: %s", rawURL)
		}
		return matches[1], nil

	case "vm.tiktok.com", "vt.tiktok.com":
		// Short URLs like vm.tiktok.com/ZMxxxxxx/
		path := strings.Trim(parsed.Path, "/")
		if path == "" {
			return "", fmt.Errorf("empty TikTok short URL path: %s", rawURL)
		}
		matches := tiktokShortRegex.FindStringSubmatch("/" + path)
		if matches == nil || len(matches) < 2 {
			return "", fmt.Errorf("could not parse TikTok short URL: %s", rawURL)
		}
		return matches[1], nil

	default:
		return "", fmt.Errorf("not a TikTok URL: %s", rawURL)
	}
}

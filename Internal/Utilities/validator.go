package utilities

import (
	"fmt"
	"net/url"
	"strings"
)

// ValidationResult holds the outcome of validating and analysing a user-supplied URL.
type ValidationResult struct {
	OriginalURL string   `json:"original_url"`
	Platform    Platform `json:"platform"`
	ContentID   string   `json:"content_id"`
	ContentType string   `json:"content_type"` // "video", "post", "reel", "short", etc.
	IsValid     bool     `json:"is_valid"`
	Error       string   `json:"error,omitempty"`
}

// ValidateURL checks whether the raw URL is well-formed and belongs to a
// supported platform. It returns a ValidationResult with the platform,
// content ID, and content type already resolved.
func ValidateURL(rawURL string) ValidationResult {
	result := ValidationResult{OriginalURL: rawURL}

	// Basic parse check.
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		result.Error = "URL is empty"
		return result
	}

	// Ensure scheme
	normalized := rawURL
	if !strings.Contains(normalized, "://") {
		normalized = "https://" + normalized
	}

	parsed, err := url.Parse(normalized)
	if err != nil || parsed.Host == "" {
		result.Error = fmt.Sprintf("malformed URL: %s", rawURL)
		return result
	}

	// Detect platform.
	platform := DetectPlatform(rawURL)
	result.Platform = platform

	if platform == PlatformUnknown {
		result.Error = fmt.Sprintf("unsupported platform for URL: %s", rawURL)
		return result
	}

	// Extract the content ID according to the platform.
	switch platform {
	case PlatformYouTube:
		id, err := ExtractYouTubeID(rawURL)
		if err != nil {
			result.Error = err.Error()
			return result
		}
		result.ContentID = id
		result.ContentType = detectYouTubeContentType(parsed)

	case PlatformTwitter:
		id, err := ExtractTwitterID(rawURL)
		if err != nil {
			result.Error = err.Error()
			return result
		}
		result.ContentID = id
		result.ContentType = "tweet"

	case PlatformInstagram:
		id, err := ExtractInstagramID(rawURL)
		if err != nil {
			result.Error = err.Error()
			return result
		}
		result.ContentID = id
		result.ContentType = detectInstagramContentType(parsed)

	case PlatformTikTok:
		id, err := ExtractTikTokID(rawURL)
		if err != nil {
			result.Error = err.Error()
			return result
		}
		result.ContentID = id
		result.ContentType = "video"
	}

	result.IsValid = true
	return result
}

// detectYouTubeContentType infers whether the URL points to a regular video,
// a short, a live stream, etc.
func detectYouTubeContentType(u *url.URL) string {
	path := u.Path
	switch {
	case strings.HasPrefix(path, "/shorts/"):
		return "short"
	case strings.HasPrefix(path, "/live/"):
		return "live"
	default:
		return "video"
	}
}

// detectInstagramContentType infers whether the URL points to a post, reel, or
// IGTV video.
func detectInstagramContentType(u *url.URL) string {
	path := u.Path
	switch {
	case strings.HasPrefix(path, "/reel/"), strings.HasPrefix(path, "/reels/"):
		return "reel"
	case strings.HasPrefix(path, "/tv/"):
		return "igtv"
	default:
		return "post"
	}
}

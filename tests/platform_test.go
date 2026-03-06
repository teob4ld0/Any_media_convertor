package utilities_test

import (
	"testing"

	utilities "github.com/teoba/any-media-convertor/Utilities"
)

func TestDetectPlatform(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected utilities.Platform
	}{
		// YouTube
		{"YouTube standard", "https://www.youtube.com/watch?v=P7dLp2Eb7cg&list=RDP7dLp2Eb7cg&start_radio=1", utilities.PlatformYouTube},
		{"YouTube short URL", "https://youtu.be/P7dLp2Eb7cg", utilities.PlatformYouTube},
		{"YouTube mobile", "https://m.youtube.com/watch?v=P7dLp2Eb7cg", utilities.PlatformYouTube},
		{"YouTube no scheme", "youtube.com/watch?v=P7dLp2Eb7cg", utilities.PlatformYouTube},
		{"YouTube shorts", "https://www.youtube.com/shorts/P7dLp2Eb7cg", utilities.PlatformYouTube},

		// Twitter / X
		{"Twitter standard", "https://twitter.com/user/status/123456789", utilities.PlatformTwitter},
		{"X standard", "https://x.com/user/status/123456789", utilities.PlatformTwitter},
		{"X with www", "https://www.x.com/user/status/123456789", utilities.PlatformTwitter},
		{"X mobile", "https://mobile.x.com/user/status/123456789", utilities.PlatformTwitter},

		// Instagram
		{"Instagram post", "https://www.instagram.com/p/ABC123def/", utilities.PlatformInstagram},
		{"Instagram reel", "https://www.instagram.com/reel/ABC123def/", utilities.PlatformInstagram},
		{"Instagram mobile", "https://m.instagram.com/p/ABC123def/", utilities.PlatformInstagram},
		{"Instagram no scheme", "instagram.com/p/ABC123def/", utilities.PlatformInstagram},

		// TikTok
		{"TikTok standard", "https://www.tiktok.com/@user/video/1234567890123456789", utilities.PlatformTikTok},
		{"TikTok short URL", "https://vm.tiktok.com/ZMxxxxxx/", utilities.PlatformTikTok},
		{"TikTok vt short", "https://vt.tiktok.com/ZSxxxxxx/", utilities.PlatformTikTok},
		{"TikTok mobile", "https://m.tiktok.com/@user/video/1234567890123456789", utilities.PlatformTikTok},

		// Unknown
		{"Unknown platform", "https://example.com/video/123", utilities.PlatformUnknown},
		{"Empty string", "", utilities.PlatformUnknown},
		{"Garbage", "not a url at all !!!", utilities.PlatformUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilities.DetectPlatform(tt.url)
			if got != tt.expected {
				t.Errorf("DetectPlatform(%q) = %q, want %q", tt.url, got, tt.expected)
			}
		})
	}
}

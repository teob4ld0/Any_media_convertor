package utilities_test

import (
	"testing"

	utilities "github.com/teoba/any-media-convertor/Utilities"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		wantValid       bool
		wantPlatform    utilities.Platform
		wantID          string
		wantContentType string
	}{
		// YouTube — solo se pasa la URL, el ID se verifica como no-vacío.
		{
			"YouTube video",
			"https://www.youtube.com/watch?v=P7dLp2Eb7cg&list=RDP7dLp2Eb7cg&start_radio=1",
			true, utilities.PlatformYouTube, "", "video",
		},
		{
			"YouTube short",
			"https://www.youtube.com/shorts/P7dLp2Eb7cg",
			true, utilities.PlatformYouTube, "", "short",
		},
		{
			"YouTube live",
			"https://www.youtube.com/live/P7dLp2Eb7cg",
			true, utilities.PlatformYouTube, "", "live",
		},

		// Twitter
		{
			"Twitter tweet",
			"https://twitter.com/user/status/123456789",
			true, utilities.PlatformTwitter, "123456789", "tweet",
		},
		{
			"X tweet",
			"https://x.com/user/status/987654321",
			true, utilities.PlatformTwitter, "987654321", "tweet",
		},

		// Instagram
		{
			"Instagram post",
			"https://www.instagram.com/p/CxYzAbCdEf/",
			true, utilities.PlatformInstagram, "CxYzAbCdEf", "post",
		},
		{
			"Instagram reel",
			"https://www.instagram.com/reel/CxYzAbCdEf/",
			true, utilities.PlatformInstagram, "CxYzAbCdEf", "reel",
		},
		{
			"Instagram tv",
			"https://www.instagram.com/tv/CxYzAbCdEf/",
			true, utilities.PlatformInstagram, "CxYzAbCdEf", "igtv",
		},

		// TikTok
		{
			"TikTok video",
			"https://www.tiktok.com/@user/video/1234567890123456789",
			true, utilities.PlatformTikTok, "1234567890123456789", "video",
		},
		{
			"TikTok short",
			"https://vm.tiktok.com/ZMxxxxxx/",
			true, utilities.PlatformTikTok, "ZMxxxxxx", "video",
		},

		// Invalid
		{
			"empty URL",
			"",
			false, "", "", "",
		},
		{
			"unsupported site",
			"https://example.com/video/123",
			false, utilities.PlatformUnknown, "", "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utilities.ValidateURL(tt.url)
			if result.IsValid != tt.wantValid {
				t.Errorf("ValidateURL(%q).IsValid = %v, want %v (error: %s)",
					tt.url, result.IsValid, tt.wantValid, result.Error)
			}
			if result.IsValid {
				if result.Platform != tt.wantPlatform {
					t.Errorf("ValidateURL(%q).Platform = %q, want %q", tt.url, result.Platform, tt.wantPlatform)
				}
				if result.ContentID == "" {
					t.Errorf("ValidateURL(%q).ContentID is empty, expected a value", tt.url)
				}
				// Si wantID está especificado, verificar match exacto; sino solo que no esté vacío.
				if tt.wantID != "" && result.ContentID != tt.wantID {
					t.Errorf("ValidateURL(%q).ContentID = %q, want %q", tt.url, result.ContentID, tt.wantID)
				}
				if result.ContentType != tt.wantContentType {
					t.Errorf("ValidateURL(%q).ContentType = %q, want %q", tt.url, result.ContentType, tt.wantContentType)
				}
				t.Logf("Extracted ID: %s", result.ContentID)
			}
		})
	}
}

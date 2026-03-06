package utilities_test

import (
	"testing"

	utilities "github.com/teoba/any-media-convertor/Utilities"
)

// ---------- YouTube ----------

func TestExtractYouTubeID(t *testing.T) {
	// URLs válidas: solo se pasa la URL, la función debe extraer un ID de 11 caracteres.
	validURLs := []struct {
		name string
		url  string
	}{
		{"watch standard", "https://www.youtube.com/watch?v=P7dLp2Eb7cg&list=RDP7dLp2Eb7cg&start_radio=1"},
		{"watch with extra params", "https://www.youtube.com/watch?v=P7dLp2Eb7cg&t=42s"},
		{"short URL", "https://youtu.be/P7dLp2Eb7cg"},
		{"short URL with params", "https://youtu.be/P7dLp2Eb7cg?t=10"},
		{"embed", "https://www.youtube.com/embed/P7dLp2Eb7cg"},
		{"shorts", "https://www.youtube.com/shorts/P7dLp2Eb7cg"},
		{"live", "https://www.youtube.com/live/P7dLp2Eb7cg"},
		{"mobile", "https://m.youtube.com/watch?v=P7dLp2Eb7cg"},
		{"no scheme", "youtube.com/watch?v=P7dLp2Eb7cg"},
	}

	for _, tt := range validURLs {
		t.Run(tt.name, func(t *testing.T) {
			id, err := utilities.ExtractYouTubeID(tt.url)
			if err != nil {
				t.Errorf("ExtractYouTubeID(%q) unexpected error: %v", tt.url, err)
				return
			}
			if len(id) != 11 {
				t.Errorf("ExtractYouTubeID(%q) returned ID %q (len %d), want 11 chars", tt.url, id, len(id))
			}
			t.Logf("Extracted ID: %s", id)
		})
	}

	// URLs inválidas: la función debe devolver error.
	invalidURLs := []struct {
		name string
		url  string
	}{
		{"empty", ""},
		{"not youtube", "https://example.com/watch?v=abc"},
		{"missing v param", "https://www.youtube.com/watch"},
		{"invalid id too short", "https://www.youtube.com/watch?v=short"},
	}

	for _, tt := range invalidURLs {
		t.Run(tt.name, func(t *testing.T) {
			_, err := utilities.ExtractYouTubeID(tt.url)
			if err == nil {
				t.Errorf("ExtractYouTubeID(%q) expected error, got nil", tt.url)
			}
		})
	}
}

// ---------- Twitter / X ----------

func TestExtractTwitterID(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantID  string
		wantErr bool
	}{
		{"twitter standard", "https://twitter.com/elonmusk/status/1234567890", "1234567890", false},
		{"x.com standard", "https://x.com/elonmusk/status/9876543210", "9876543210", false},
		{"x.com with www", "https://www.x.com/someone/status/1111111111", "1111111111", false},
		{"mobile.x.com", "https://mobile.x.com/user/status/2222222222", "2222222222", false},
		{"with trailing stuff", "https://twitter.com/user/status/3333333333/photo/1", "3333333333", false},
		{"no scheme", "twitter.com/user/status/4444444444", "4444444444", false},

		// Error cases
		{"not twitter", "https://youtube.com/watch?v=abc", "", true},
		{"no status", "https://twitter.com/user", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utilities.ExtractTwitterID(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractTwitterID(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
				return
			}
			if got != tt.wantID {
				t.Errorf("ExtractTwitterID(%q) = %q, want %q", tt.url, got, tt.wantID)
			}
		})
	}
}

// ---------- Instagram ----------

func TestExtractInstagramID(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantID  string
		wantErr bool
	}{
		{"post", "https://www.instagram.com/p/CxYzAbCdEf/", "CxYzAbCdEf", false},
		{"reel", "https://www.instagram.com/reel/CxYzAbCdEf/", "CxYzAbCdEf", false},
		{"reels", "https://www.instagram.com/reels/CxYzAbCdEf/", "CxYzAbCdEf", false},
		{"tv", "https://www.instagram.com/tv/CxYzAbCdEf/", "CxYzAbCdEf", false},
		{"mobile", "https://m.instagram.com/p/CxYzAbCdEf/", "CxYzAbCdEf", false},
		{"no scheme", "instagram.com/p/CxYzAbCdEf/", "CxYzAbCdEf", false},

		// Error cases
		{"not instagram", "https://twitter.com/p/abc", "", true},
		{"profile URL", "https://www.instagram.com/someuser/", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utilities.ExtractInstagramID(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractInstagramID(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
				return
			}
			if got != tt.wantID {
				t.Errorf("ExtractInstagramID(%q) = %q, want %q", tt.url, got, tt.wantID)
			}
		})
	}
}

// ---------- TikTok ----------

func TestExtractTikTokID(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantID  string
		wantErr bool
	}{
		{"standard", "https://www.tiktok.com/@user/video/1234567890123456789", "1234567890123456789", false},
		{"mobile", "https://m.tiktok.com/@user/video/9876543210987654321", "9876543210987654321", false},
		{"vm short", "https://vm.tiktok.com/ZMxxxxxx/", "ZMxxxxxx", false},
		{"vt short", "https://vt.tiktok.com/ZSyyyyyy/", "ZSyyyyyy", false},
		{"no scheme", "tiktok.com/@user/video/1234567890123456789", "1234567890123456789", false},

		// Error cases
		{"not tiktok", "https://youtube.com/@user/video/123", "", true},
		{"no video path", "https://www.tiktok.com/@user", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utilities.ExtractTikTokID(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractTikTokID(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
				return
			}
			if got != tt.wantID {
				t.Errorf("ExtractTikTokID(%q) = %q, want %q", tt.url, got, tt.wantID)
			}
		})
	}
}

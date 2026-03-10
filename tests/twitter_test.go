package utilities_test

import (
	"errors"
	"testing"

	extractors "github.com/teoba/any-media-convertor/Internal/Extractors"
)

// mockHTTPClient implements extractors.HTTPClient returning canned responses.
type mockHTTPClient struct {
	postBody string
	postErr  error
	getBody  string
	getErr   error
}

func (m *mockHTTPClient) Get(_ string, _ map[string]string) ([]byte, error) {
	return []byte(m.getBody), m.getErr
}

func (m *mockHTTPClient) Post(_ string, _ map[string]string, _ string) ([]byte, error) {
	return []byte(m.postBody), m.postErr
}

// tweetJSON builds a minimal syndication API response containing the given variants.
func tweetJSON(variants ...extractors.VideoVariant) string {
	variantsJSON := "["
	for i, v := range variants {
		if i > 0 {
			variantsJSON += ","
		}
		variantsJSON += `{"bitrate":` + itoa(v.Bitrate) + `,"content_type":"` + v.ContentType + `","url":"` + v.URL + `"}`
	}
	variantsJSON += "]"

	return `{
		"mediaDetails": [{
			"type": "video",
			"video_info": { "variants": ` + variantsJSON + ` }
		}]
	}`
}

// itoa converts an int to its string representation without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

func TestFetchVideoURLsWithClient_ReturnsVariantsSortedByBitrate(t *testing.T) {
	mock := &mockHTTPClient{
		getBody: tweetJSON(
			extractors.VideoVariant{Bitrate: 832000, ContentType: "video/mp4", URL: "https://video.twimg.com/low.mp4"},
			extractors.VideoVariant{Bitrate: 0, ContentType: "application/x-mpegURL", URL: "https://video.twimg.com/playlist.m3u8"},
			extractors.VideoVariant{Bitrate: 2176000, ContentType: "video/mp4", URL: "https://video.twimg.com/high.mp4"},
		),
	}

	variants, err := extractors.FetchVideoURLsWithClient("123456789", mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(variants) != 3 {
		t.Fatalf("expected 3 variants, got %d", len(variants))
	}
	// Highest bitrate must come first.
	if variants[0].Bitrate != 2176000 {
		t.Errorf("variants[0].Bitrate = %d, want 2176000", variants[0].Bitrate)
	}
	if variants[1].Bitrate != 832000 {
		t.Errorf("variants[1].Bitrate = %d, want 832000", variants[1].Bitrate)
	}
	// HLS (bitrate 0) must be last.
	if variants[2].Bitrate != 0 {
		t.Errorf("variants[2].Bitrate = %d, want 0 (HLS)", variants[2].Bitrate)
	}
	if variants[2].ContentType != "application/x-mpegURL" {
		t.Errorf("variants[2].ContentType = %q, want application/x-mpegURL", variants[2].ContentType)
	}
}

func TestFetchVideoURLsWithClient_TweetAPIFails(t *testing.T) {
	mock := &mockHTTPClient{
		getErr: errors.New("timeout"),
	}

	_, err := extractors.FetchVideoURLsWithClient("123456789", mock)
	if err == nil {
		t.Fatal("expected error when tweet API call fails, got nil")
	}
}

func TestFetchVideoURLsWithClient_BadTweetJSON(t *testing.T) {
	mock := &mockHTTPClient{
		getBody: `not json at all {{{`,
	}

	_, err := extractors.FetchVideoURLsWithClient("123456789", mock)
	if err == nil {
		t.Fatal("expected error on malformed tweet JSON, got nil")
	}
}

func TestFetchVideoURLsWithClient_NoVideoInTweet(t *testing.T) {
	// Tweet with a photo, not a video.
	mock := &mockHTTPClient{
		getBody: `{"mediaDetails":[{"type":"photo","video_info":{"variants":[]}}]}`,
	}

	_, err := extractors.FetchVideoURLsWithClient("123456789", mock)
	if err == nil {
		t.Fatal("expected error when tweet has no video, got nil")
	}
}

func TestFetchVideoURLsWithClient_AnimatedGIFIsAccepted(t *testing.T) {
	mock := &mockHTTPClient{
		getBody: `{
			"mediaDetails": [{
				"type": "animated_gif",
				"video_info": { "variants": [
					{"bitrate": 0, "content_type": "video/mp4", "url": "https://video.twimg.com/gif.mp4"}
				]}
			}]
		}`,
	}

	variants, err := extractors.FetchVideoURLsWithClient("123456789", mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(variants) != 1 {
		t.Fatalf("expected 1 variant, got %d", len(variants))
	}
	if variants[0].ContentType != "video/mp4" {
		t.Errorf("ContentType = %q, want video/mp4", variants[0].ContentType)
	}
}

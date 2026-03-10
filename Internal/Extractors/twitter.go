// Package extractors implements per-platform media URL extraction.
package extractors

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"

	"github.com/teoba/any-media-convertor/Internal/client"
)

// HTTPClient is the subset of *client.Client used by this package.
// Defined as an interface to allow test injection.
type HTTPClient interface {
	Get(url string, headers map[string]string) ([]byte, error)
	Post(url string, headers map[string]string, body string) ([]byte, error)
}

// VideoVariant represents one downloadable stream for a tweet video.
// MP4 variants have a non-zero Bitrate; the HLS playlist has Bitrate == 0.
type VideoVariant struct {
	Bitrate     int    `json:"bitrate"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

// syndicationResp is the shape returned by cdn.syndication.twimg.com.
type syndicationResp struct {
	MediaDetails []struct {
		Type      string `json:"type"`
		VideoInfo struct {
			Variants []VideoVariant `json:"variants"`
		} `json:"video_info"`
	} `json:"mediaDetails"`
}

// FetchVideoURLs returns all stream variants for the given tweet ID, sorted by
// bitrate descending (index 0 = highest quality MP4, last = HLS playlist).
func FetchVideoURLs(tweetID string) ([]VideoVariant, error) {
	return FetchVideoURLsWithClient(tweetID, client.New())
}

// FetchVideoURLsWithClient is like FetchVideoURLs but accepts a custom HTTPClient.
// Useful for testing with a mock client.
func FetchVideoURLsWithClient(tweetID string, c HTTPClient) ([]VideoVariant, error) {
	// The syndication API is public — no credentials required.
	// A numeric token derived from the tweet ID is expected by the endpoint.
	token := syndicationToken(tweetID)
	url := "https://cdn.syndication.twimg.com/tweet-result?id=" + tweetID + "&lang=en&token=" + token

	body, err := c.Get(url, map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})
	if err != nil {
		return nil, fmt.Errorf("tweet %s: %w", tweetID, err)
	}

	var resp syndicationResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing tweet JSON: %w", err)
	}

	var variants []VideoVariant
	for _, media := range resp.MediaDetails {
		if media.Type == "video" || media.Type == "animated_gif" {
			variants = append(variants, media.VideoInfo.Variants...)
		}
	}

	if len(variants) == 0 {
		return nil, fmt.Errorf("no video found in tweet")
	}

	// Highest bitrate first; HLS (bitrate 0) ends up last.
	sort.Slice(variants, func(i, j int) bool {
		return variants[i].Bitrate > variants[j].Bitrate
	})

	return variants, nil
}

// syndicationToken computes the numeric token expected by the syndication API.
// Formula: round(tweetID / 1e15 * π)
func syndicationToken(tweetID string) string {
	id, ok := new(big.Int).SetString(tweetID, 10)
	if !ok {
		return "0"
	}

	prec := uint(128)
	pi, _, _ := new(big.Float).SetPrec(prec).Parse("3.14159265358979323846", 10)
	e15 := new(big.Float).SetPrec(prec).SetInt64(1_000_000_000_000_000)

	f := new(big.Float).SetPrec(prec).SetInt(id)
	f.Quo(f, e15)
	f.Mul(f, pi)

	// Round to nearest integer.
	result, _ := f.Int(nil)
	return result.String()
}

// Package extractors implements per-platform media URL extraction.
package extractors

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/teoba/any-media-convertor/Internal/client"
)

// bearerToken is the public Twitter web-app bearer token used by the official
// site and open-source tools like yt-dlp. It is not a personal API key.
const bearerToken = "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"

// VideoVariant represents one downloadable stream for a tweet video.
// MP4 variants have a non-zero Bitrate; the HLS playlist has Bitrate == 0.
type VideoVariant struct {
	Bitrate     int    `json:"bitrate"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

// --- internal JSON shapes ---

type guestTokenResp struct {
	GuestToken string `json:"guest_token"`
}

type tweetResp struct {
	ExtendedEntities struct {
		Media []struct {
			Type      string `json:"type"`
			VideoInfo struct {
				Variants []VideoVariant `json:"variants"`
			} `json:"video_info"`
		} `json:"media"`
	} `json:"extended_entities"`
}

// FetchVideoURLs returns all stream variants for the given tweet ID, sorted by
// bitrate descending (index 0 = highest quality MP4, last = HLS playlist).
func FetchVideoURLs(tweetID string) ([]VideoVariant, error) {
	c := client.New()

	guestToken, err := fetchGuestToken(c)
	if err != nil {
		return nil, fmt.Errorf("guest token: %w", err)
	}

	variants, err := fetchVariants(c, guestToken, tweetID)
	if err != nil {
		return nil, fmt.Errorf("tweet %s: %w", tweetID, err)
	}

	return variants, nil
}

// fetchGuestToken activates a short-lived guest session and returns the token.
func fetchGuestToken(c *client.Client) (string, error) {
	body, err := c.Post(
		"https://api.twitter.com/1.1/guest/activate.json",
		map[string]string{"Authorization": "Bearer " + bearerToken},
		"",
	)
	if err != nil {
		return "", err
	}

	var resp guestTokenResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}
	if resp.GuestToken == "" {
		return "", fmt.Errorf("empty guest token in response")
	}
	return resp.GuestToken, nil
}

// fetchVariants queries the Twitter API and returns all video variants in the tweet.
func fetchVariants(c *client.Client, guestToken, tweetID string) ([]VideoVariant, error) {
	url := "https://twitter.com/i/api/1.1/statuses/show.json?id=" + tweetID + "&tweet_mode=extended"

	body, err := c.Get(url, map[string]string{
		"Authorization":              "Bearer " + bearerToken,
		"x-guest-token":              guestToken,
		"x-twitter-active-user":      "yes",
		"x-twitter-client-language":  "en",
		"User-Agent":                 "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})
	if err != nil {
		return nil, err
	}

	var tweet tweetResp
	if err := json.Unmarshal(body, &tweet); err != nil {
		return nil, fmt.Errorf("parsing tweet JSON: %w", err)
	}

	var variants []VideoVariant
	for _, media := range tweet.ExtendedEntities.Media {
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

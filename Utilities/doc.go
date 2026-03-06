// Package utilities provides URL analysis tools for the Any Media Convertor project.
//
// It implements three layers of functionality:
//
//  1. Platform detection — [DetectPlatform] identifies which supported platform
//     (YouTube, Twitter/X, Instagram, TikTok) a given URL belongs to.
//
//  2. Content ID extraction — [ExtractYouTubeID], [ExtractTwitterID],
//     [ExtractInstagramID] and [ExtractTikTokID] parse the URL and return the
//     unique identifier of the content (video ID, tweet ID, shortcode, etc.).
//
//  3. Full validation — [ValidateURL] combines detection, extraction and content
//     type inference into a single call, returning a [ValidationResult].
//
// All functions accept raw user input: missing schemes (https://) and "www."
// prefixes are handled transparently. No external dependencies are used.
package utilities

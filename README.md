# Any Media Convertor

Tool to download media from posts on major platforms: **YouTube**, **X (Twitter)**, **Instagram** and **TikTok** — without relying on yt-dlp. Uses each platform's own public APIs.

> Currently in active development. Twitter/X is the only fully working extractor.

---

## How it works

```
input URL → validator → platform detector → extractor → downloader
```

1. **Validator** — checks the URL is well-formed and belongs to a supported platform, extracts the content ID.
2. **Platform detector** — identifies Twitter, YouTube, Instagram or TikTok from the host.
3. **Extractor** — calls the platform's public API to get the direct stream URLs.
4. **Downloader** — fetches the media (HLS/MP4) and saves it locally. *(in progress)*

### Twitter/X extractor

Uses Twitter's public **syndication API** (`cdn.syndication.twimg.com`) — the same endpoint used by Twitter's own embed widget and by yt-dlp. No credentials or API key required. Returns all available stream variants (multiple MP4 bitrates + HLS playlist), sorted highest quality first.

---

## Project structure

```
Any_media_convertor/
├── main.go                         # Entry point — CLI, dispatches per platform
├── go.mod
├── Internal/
│   ├── Extractors/
│   │   ├── twitter.go              # Twitter/X — fully working
│   │   ├── youtube.go              # stub
│   │   └── instagram.go            # stub
│   ├── Utilities/
│   │   ├── doc.go
│   │   ├── platform.go             # Platform detection from URL
│   │   ├── extractor.go            # Content ID extraction per platform
│   │   └── validator.go            # Full URL validation pipeline
│   ├── client/
│   │   └── http_client.go          # Shared HTTP client
│   ├── downloader/
│   │   ├── file.go                 # File download (stub)
│   │   └── hls.go                  # HLS/m3u8 download (stub)
│   └── Streams/
│       └── m3u8.go                 # m3u8 parsing (stub)
├── Frontend/                       # Web UI (planned)
├── tests/
│   ├── platform_test.go
│   ├── extractor_test.go
│   ├── validator_test.go
│   ├── twitter_test.go
│   └── http_client_test.go
└── README.md
```

---

## Supported URL formats

| Platform | Supported formats |
|---|---|
| Twitter/X | `x.com/user/status/ID`, `twitter.com/user/status/ID`, `mobile.x.com/...` |
| YouTube | `youtube.com/watch?v=ID`, `youtu.be/ID`, `/shorts/`, `/live/`, `/embed/` |
| Instagram | `/p/SHORTCODE/`, `/reel/SHORTCODE/`, `/tv/SHORTCODE/` |
| TikTok | `tiktok.com/@user/video/ID`, `vm.tiktok.com/...`, `vt.tiktok.com/...` |

---

## Usage

```powershell
go run . <URL>
```

Example:

```powershell
go run . https://x.com/user/status/1234567890123456789
```

Output:

```
Tweet 1234567890123456789 — 3 stream(s) found:

  [video/mp4] 2176000 bps
  https://video.twimg.com/...mp4?tag=...

  [video/mp4] 832000 bps
  https://video.twimg.com/...mp4?tag=...

  [application/x-mpegURL]
  https://video.twimg.com/...m3u8
```

---

## Running tests

```powershell
go test ./tests/ -v
```

---

## Tech stack

- **Language:** Go 1.22+
- **Dependencies:** Standard library only — no external packages.
- **No yt-dlp / ffmpeg:** uses each platform's own public APIs directly.

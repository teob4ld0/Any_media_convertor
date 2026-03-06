# Any Media Convertor

Web app that allows you to download media from any post on major platforms: **YouTube**, **X (Twitter)**, **Instagram** and **TikTok** — without relying on yt-dlp.

---

## Project structure

```
Any_media_convertor/
├── Utilities/              # Core utility functions (Go)
│   ├── doc.go              # Package-level documentation
│   ├── platform.go         # Platform detection from URL
│   ├── extractor.go        # Content ID extraction per platform
│   └── validator.go        # Full URL validation pipeline
├── tests/                  # Unit tests (external package: utilities_test)
│   ├── platform_test.go    # Tests for DetectPlatform
│   ├── extractor_test.go   # Tests for Extract*ID functions
│   └── validator_test.go   # Tests for ValidateURL
├── DTOs/                   # Data Transfer Objects (planned)
├── EachPageController/     # Per-platform download logic (planned)
├── Formats/                # Output format handling (planned)
├── Frontend/               # Web UI (planned)
├── go.mod                  # Go module definition
└── README.md
```

---

## Utilities — Reference

### `platform.go` — Platform detection

| Symbol | Description |
|---|---|
| `Platform` | Type alias (`string`) representing a supported platform. |
| `PlatformYouTube`, `PlatformTwitter`, `PlatformInstagram`, `PlatformTikTok`, `PlatformUnknown` | Platform constants. |
| `DetectPlatform(rawURL string) Platform` | Parses the URL, normalises the host, and returns which platform it belongs to. Returns `PlatformUnknown` for invalid or unrecognised URLs. |

**Supported hosts:**

| Platform | Hosts |
|---|---|
| YouTube | `youtube.com`, `youtu.be`, `m.youtube.com` |
| Twitter/X | `twitter.com`, `x.com`, `mobile.x.com` |
| Instagram | `instagram.com`, `m.instagram.com` |
| TikTok | `tiktok.com`, `m.tiktok.com`, `vm.tiktok.com`, `vt.tiktok.com` |

> All hosts are matched without `www.` prefix and case-insensitively. URLs without a scheme (`https://`) are handled automatically.

---

### `extractor.go` — Content ID extraction

Each function takes a raw URL and returns the content identifier.

#### `ExtractYouTubeID(rawURL string) (string, error)`

Returns the 11-character video ID.

| Format | Example |
|---|---|
| Standard watch | `https://www.youtube.com/watch?v=VIDEO_ID` |
| Short URL | `https://youtu.be/VIDEO_ID` |
| Embed | `https://www.youtube.com/embed/VIDEO_ID` |
| Shorts | `https://www.youtube.com/shorts/VIDEO_ID` |
| Live | `https://www.youtube.com/live/VIDEO_ID` |
| Mobile | `https://m.youtube.com/watch?v=VIDEO_ID` |

Extra query parameters (`&list=...`, `&t=...`) are ignored correctly.

#### `ExtractTwitterID(rawURL string) (string, error)`

Returns the numeric tweet ID.

| Format | Example |
|---|---|
| Twitter | `https://twitter.com/user/status/1234567890` |
| X | `https://x.com/user/status/1234567890` |
| Mobile X | `https://mobile.x.com/user/status/1234567890` |

#### `ExtractInstagramID(rawURL string) (string, error)`

Returns the post/reel shortcode.

| Format | Example |
|---|---|
| Post | `https://www.instagram.com/p/SHORTCODE/` |
| Reel | `https://www.instagram.com/reel/SHORTCODE/` |
| Reels | `https://www.instagram.com/reels/SHORTCODE/` |
| IGTV | `https://www.instagram.com/tv/SHORTCODE/` |

#### `ExtractTikTokID(rawURL string) (string, error)`

Returns the numeric video ID or the short-link code.

| Format | Example |
|---|---|
| Standard | `https://www.tiktok.com/@user/video/1234567890123456789` |
| VM short | `https://vm.tiktok.com/ZMxxxxxx/` |
| VT short | `https://vt.tiktok.com/ZSxxxxxx/` |
| Mobile | `https://m.tiktok.com/@user/video/1234567890123456789` |

---

### `validator.go` — Full validation pipeline

#### `ValidateURL(rawURL string) ValidationResult`

One-call function that:
1. Validates the URL format.
2. Detects the platform via `DetectPlatform`.
3. Extracts the content ID via the corresponding `Extract*ID` function.
4. Determines the content type (`video`, `short`, `live`, `tweet`, `post`, `reel`, `igtv`).

Returns a `ValidationResult` struct:

```go
type ValidationResult struct {
    OriginalURL string   // The URL as provided by the user
    Platform    Platform // Detected platform
    ContentID   string   // Extracted ID/shortcode
    ContentType string   // "video", "short", "live", "tweet", "post", "reel", "igtv"
    IsValid     bool     // true if everything succeeded
    Error       string   // Error message (empty on success)
}
```

---

## Running the tests

```powershell
# If Go was installed via ZIP in ~/go_sdk, set PATH first:
$env:GOROOT = "$env:USERPROFILE\go_sdk\go"
$env:Path   = "$env:USERPROFILE\go_sdk\go\bin;$env:Path"

# Run all tests (verbose)
go test ./tests/ -v

# Run specific test suites
go test ./tests/ -v -run TestDetectPlatform
go test ./tests/ -v -run TestExtractYouTubeID
go test ./tests/ -v -run TestExtractTwitterID
go test ./tests/ -v -run TestExtractInstagramID
go test ./tests/ -v -run TestExtractTikTokID
go test ./tests/ -v -run TestValidateURL
```

---

## Tech stack

- **Language:** Go 1.22+
- **Dependencies:** Standard library only (`net/url`, `regexp`, `strings`, `fmt`)
- **No external tools:** yt-dlp, ffmpeg, etc. are intentionally avoided for learning purposes.

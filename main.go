/*
Flujo:

input URL
   ↓
validator
   ↓
platform detector
   ↓
extractor
   ↓
downloader
*/

package main

import (
	"fmt"
	"os"

	extractors "github.com/teoba/any-media-convertor/Internal/Extractors"
	utilities "github.com/teoba/any-media-convertor/Internal/Utilities"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: any-media-convertor <URL>")
		os.Exit(1)
	}

	result := utilities.ValidateURL(os.Args[1])
	if !result.IsValid {
		fmt.Fprintln(os.Stderr, "error:", result.Error)
		os.Exit(1)
	}

	switch result.Platform {
	case utilities.PlatformTwitter:
		variants, err := extractors.FetchVideoURLs(result.ContentID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		fmt.Printf("Tweet %s — %d stream(s) found:\n\n", result.ContentID, len(variants))
		for _, v := range variants {
			if v.Bitrate > 0 {
				fmt.Printf("  [%s] %d bps\n  %s\n\n", v.ContentType, v.Bitrate, v.URL)
			} else {
				fmt.Printf("  [%s]\n  %s\n\n", v.ContentType, v.URL)
			}
		}

	default:
		fmt.Fprintf(os.Stderr, "platform %q not yet supported\n", result.Platform)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/stumburs/vpp/cmd/vpp/download"
)

func main() {
	downloadMode := flag.Bool("dl", false, "-dl <URL|ID> Download a video using a URL or ID.")
	downloadInChunks := flag.Bool("chunk", false, "Download the video in chunks. This might resolve longer videos failing to download due to EOF errors.")

	flag.Parse()

	// Validate that -chunked requires -dl
	if *downloadInChunks && !*downloadMode {
		fmt.Println("ERROR: The -chunk flag requires -dl to be specified.")
		os.Exit(1)
	}

	args := flag.Args()

	if *downloadMode {
		if len(args) < 1 {
			fmt.Println("ERROR: You must provide a video URL or ID when using the -dl flag.")
			os.Exit(1)
		}

		videoURL := args[0]

		dl := download.Downloader{OutputDir: "test"}
		ctx := context.Background()

		video, err := dl.Client.GetVideoContext(ctx, videoURL)
		if err != nil {
			panic(err)
		}

		// for _, format := range video.Formats {
		// 	fmt.Println(format.QualityLabel)
		// }

		// Download highest quality as default (for now)
		qualityLabel := video.Formats[0].QualityLabel

		dl.DownloadComposite(ctx, "", video, qualityLabel, "", "")
	} else {
		fmt.Println("Use -help for usage.")
	}
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/stumburs/vpp/cmd/vpp/download"
)

func main() {
	downloadMode := flag.Bool("dl", false, "-dl <URL|ID> Download a video using a URL or ID.")
	downloadInChunks := flag.Bool("chunk", false, "Download the video in chunks. This might resolve longer videos failing to download due to EOF errors.")

	flag.Parse()

	args := flag.Args()

	if *downloadMode {
		if len(args) < 1 {
			fmt.Println("ERROR: You must provide a video URL or ID when using the -dl flag.")
			os.Exit(1)
		}

		videoURL := args[0]

		download.DownloadVideo(videoURL, *downloadInChunks)
	} else {
		fmt.Println("Use -help for usage.")
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/stumburs/vpp/cmd/vpp/download"
)

func main() {
	downloadMode := flag.Bool("dl", false, "<URL|ID> Download a video using a URL or ID.")
	videoInfo := flag.Bool("info", false, "<URL|ID> Displays all possible formats for specified video.")
	videoQualityFlag := flag.Int("q", 0, "Specifies what quality to download the video as. Use -info to view all possible formats.")
	reencodeFlag := flag.Bool("reencode", false, "After downloading, mixes the video and audio by re-encoding, instead of copying.")

	flag.Parse()

	args := flag.Args()

	if *videoInfo {
		if len(args) < 1 {
			fmt.Println("ERROR: You must provide a video URL or ID when using the -info flag.")
			os.Exit(1)
		}

		// TODO: Add -o flag to specify output file/directory
		dl := download.Downloader{OutputDir: ""}
		ctx := context.Background()
		videoURL := args[0]

		video, err := dl.Client.GetVideoContext(ctx, videoURL)

		if err != nil {
			panic(err)
		}

		download.DisplayFormats(video)
		os.Exit(0)
	}

	if *downloadMode {
		if len(args) < 1 {
			fmt.Println("ERROR: You must provide a video URL or ID when using the -dl flag.")
			os.Exit(1)
		}

		videoURL := args[0]

		dl := download.Downloader{OutputDir: ""}
		ctx := context.Background()

		video, err := dl.Client.GetVideoContext(ctx, videoURL)
		if err != nil {
			panic(err)
		}

		qualityLabel := video.Formats[*videoQualityFlag].QualityLabel

		dl.DownloadComposite(ctx, "", video, qualityLabel, "", "", *reencodeFlag)
	} else {
		fmt.Println("Use -help for usage.")
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/stumburs/vpp/cmd/vpp/convert"
	"github.com/stumburs/vpp/cmd/vpp/download"
	"github.com/stumburs/vpp/cmd/vpp/version"
)

func main() {
	// Download flags
	downloadMode := flag.Bool("dl", false, "<URL|ID> Download a video using a URL or ID.")
	videoInfo := flag.Bool("info", false, "<URL|ID> Displays all possible formats for specified video.")
	videoQualityFlag := flag.Int("q", 0, "Specifies what quality to download the video as. Use -info to view all possible formats.")
	reencodeFlag := flag.Bool("reencode", false, "After downloading, mixes the video and audio by re-encoding using x264/AAC codecs, instead of copying. This fixes embed issues with Discord.")
	downloadMP3 := flag.Bool("mp3", false, "Download only as mp3.")
	downloadWAV := flag.Bool("wav", false, "Download only as wav.")

	// Program meta flags
	versionFlag := flag.Bool("version", false, "Displays executable version.")
	vFlag := flag.Bool("v", false, "Displays executable version.")

	// Convert flags
	convertMode := flag.Bool("conv", false, "Converts a file into another file according to specific parameters.")
	sizeFlag := flag.Bool("size", false, "Convert file to specific size by changing the bitrate.")
	targetSizeMBFlag := flag.Int("mb", 0, "Megabytes")

	flag.Parse()

	// Version
	if *versionFlag || *vFlag {
		fmt.Printf("vpp version %s %s/%s\n", version.VERSION, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	userFlags := download.DownloadFlags{
		ReencodeAfterDownload: *reencodeFlag,
		DownloadAsMP3:         *downloadMP3,
		DownloadAsWAV:         *downloadWAV,
	}

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

		dl.DownloadComposite(ctx, "", video, qualityLabel, "", "", userFlags)
		os.Exit(0)
	}

	if *convertMode {
		if *sizeFlag {
			if *targetSizeMBFlag != 0 {
				if len(args) < 2 {
					fmt.Println("ERROR: You must provide TWO file paths.")
					os.Exit(1)
				}

				inputVideo := args[0]
				outputVideo := args[1]

				err := convert.ChangeVideoSize(inputVideo, outputVideo, *targetSizeMBFlag)
				if err != nil {
					fmt.Printf("ERROR: %s\n", err)
					os.Exit(1)
				}

				os.Exit(0)

			} else {
				fmt.Println("ERROR: Target size must not be 0MB.")
				os.Exit(1)
			}
		} else {
			fmt.Println("ERROR: You must provide a conversion flag.")
			os.Exit(1)
		}
		os.Exit(0)
	}

	fmt.Println("Use -help for usage.")
	os.Exit(0)
}

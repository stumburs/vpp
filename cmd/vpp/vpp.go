package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	youtube "github.com/kkdai/youtube/v2"
)

func downloadVideo(videoID string) {
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels()

	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	file, err := os.Create(video.Title + ".mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}

	fmt.Println("Downloaded video: " + video.Title)
}

func main() {
	downloadMode := flag.Bool("dl", false, "-dl <URL | ID> Set mode to download video")

	flag.Parse()

	args := flag.Args()

	if *downloadMode {
		if len(args) < 1 {
			fmt.Println("ERROR: You must provide a video URL or ID when using the -dl flag.")
			os.Exit(1)
		}

		videoURL := args[0]

		downloadVideo(videoURL)
	} else {
		fmt.Println("Download mode is not set. Use -dl flag to download a video.")
	}
}

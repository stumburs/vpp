package download

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/schollz/progressbar/v3"
)

const (
	chunkSize     = 5 * youtube.Size1Mb // 5 MB per chunk
	maxRetryCount = 3                   // Maximum retry count for each chunk
	retryWait     = 2 * time.Second     // Wait time before retrying
)

func DownloadVideo(videoID string) {
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels()
	format := &formats[0]

	stream, streamSize, err := client.GetStream(video, format)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	file, err := os.Create(video.Title + ".mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(streamSize, "Downloading "+video.Title)

	// Direct download
	if streamSize < 1000*youtube.Size1Mb {
		_, err = io.Copy(file, io.TeeReader(stream, bar))
		if err != nil {
			panic(err)
		}
	} else { // EXPERIMENTAL: Download in chunks
		for offset := int64(0); offset < streamSize; offset += chunkSize {

			// Current size of chunk
			end := offset + chunkSize - 1
			if end >= streamSize {
				end = streamSize - 1
			}

			// Try downloading the chunk with retries
			for attempt := 0; attempt < maxRetryCount; attempt++ {
				err = downloadChunk(&client, video, format, file, offset, end, bar)
				if err == nil {
					break // Continue to next chunk
				}
				fmt.Printf("Error downloading chunk (%d-%d): %v. Retrying...\n", offset, end, err)
				time.Sleep(retryWait)
			}
		}
	}

	fmt.Println("Downloaded video: " + video.Title)
}

func downloadChunk(client *youtube.Client, video *youtube.Video, format *youtube.Format, file *os.File, offset, end int64, bar *progressbar.ProgressBar) error {
	// Get a partial stream for the specified byte range
	stream, _, err := client.GetStream(video, format) // This should ideally support byte-range downloads
	if err != nil {
		return err
	}
	defer stream.Close()

	// Move file cursor to the correct offset for writing
	if _, err = file.Seek(offset, 0); err != nil {
		return err
	}

	// Use a limited reader to only read the desired chunk
	limitedReader := io.LimitReader(stream, end-offset+1)
	_, err = io.Copy(io.MultiWriter(file, bar), limitedReader)

	return err
}

package download

import (
	"fmt"

	"github.com/kkdai/youtube/v2"
)

func DisplayFormats(video *youtube.Video) {
	fmt.Printf("Available formats for: %s\n", video.Title)
	fmt.Printf("--------------------------------------------------\n")

	for idx, format := range video.Formats {

		if format.QualityLabel == "" {
			continue
		}

		fmt.Printf("%d:\n", idx)
		fmt.Printf("\t%s\n", format.QualityLabel)
	}
}

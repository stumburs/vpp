package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type VideoInfo struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
	Streams []struct {
		CodecType string `json:"codec_type"`
		BitRate   string `json:"bit_rate"`
	} `json:"streams"`
}

func getVideoInfo(filePath string) (float64, int, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to run ffprobe: %w", err)
	}

	var info VideoInfo
	if err := json.Unmarshal(out.Bytes(), &info); err != nil {
		return 0, 0, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	// Parse duration
	duration, err := strconv.ParseFloat(info.Format.Duration, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	// Find the audio stream and get its bitrate
	var audioBitrate int
	for _, stream := range info.Streams {
		if stream.CodecType == "audio" && stream.BitRate != "" {
			audioBitrate, err = strconv.Atoi(stream.BitRate)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to parse audio bitrate: %w", err)
			}
			audioBitrate /= 1000 // Convert to kbps
			break
		}
	}

	return duration, audioBitrate, nil
}

func ChangeVideoSize(inputPath string, outputPath string, targetSizeMB int) error {
	// Get video duration and audio bitrate
	durationSec, audioBitrateKbps, err := getVideoInfo(inputPath)
	if err != nil {
		return fmt.Errorf("error getting video info: %v", err)
	}

	// Convert target size from MB to bits (1 MB = 8,388,608 bits)
	targetSizeBits := targetSizeMB * 8388608

	// Calculate the target video bitrate (in kbps) after subtracting audio bitrate
	targetVideoBitrate := (int(targetSizeBits) / int(durationSec) / 1000) - audioBitrateKbps

	// Ensure the bitrate is positive
	if targetVideoBitrate <= 0 {
		return fmt.Errorf("calculated video bitrate is too low; try a larger target size or lower audio bitrate")
	}

	// Convert target video bitrate to a string for ffmpeg
	videoBitrateStr := strconv.Itoa(targetVideoBitrate) + "k"
	audioBitrateStr := strconv.Itoa(audioBitrateKbps) + "k"

	// ffmpeg command
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-b:v", videoBitrateStr,
		"-b:a", audioBitrateStr,
		"-y", outputPath,
		"-v", "quiet",
		"-stats",
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	fmt.Printf("Re-rendering %s into %s with target size: %dMB...\n", inputPath, outputPath, targetSizeMB)
	return cmd.Run()
}

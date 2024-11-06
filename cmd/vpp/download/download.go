package download

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kkdai/youtube/v2"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	youtube.Client
	OutputDir string // optional
}

func (dl *Downloader) getOutputFile(v *youtube.Video, _ *youtube.Format, outputFile string) (string, error) {
	if outputFile == "" {
		outputFile = SanitizeFilename(v.Title)

		// "There's nothing more permanent than a temporary solution"
		// outputFile += pickIdealFileExtension(format.MimeType)
		outputFile += ".mp4"
	}

	if dl.OutputDir != "" {
		if err := os.MkdirAll(dl.OutputDir, 0o755); err != nil {
			return "", err
		}
		outputFile = filepath.Join(dl.OutputDir, outputFile)
	}

	return outputFile, nil
}

func (dl *Downloader) Download(ctx context.Context, v *youtube.Video, format *youtube.Format, outputFile string) error {
	youtube.Logger.Info(
		"Downloading video",
		"id", v.ID,
		"quality", format.Quality,
		"mimeType", format.MimeType,
	)

	destFile, err := dl.getOutputFile(v, format, outputFile)
	if err != nil {
		return err
	}

	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	return dl.videoDLWorker(ctx, out, v, format)
}

func (dl *Downloader) DownloadComposite(ctx context.Context, outputFile string, v *youtube.Video, quality string, mimetype, language string, reencode bool) error {
	videoFormat, audioFormat, err := GetVideoAudioFormats(v, quality, mimetype, language)
	if err != nil {
		return err
	}

	log := youtube.Logger.With("id", v.ID)

	log.Info("Downloading composite video",
		"videoQuality", videoFormat.QualityLabel,
		"videoMimeType", videoFormat.MimeType,
		"audioMimeType", audioFormat.MimeType,
	)

	destFile, err := dl.getOutputFile(v, videoFormat, outputFile)
	if err != nil {
		return err
	}
	outputDir := filepath.Dir(destFile)

	videoFile, err := os.CreateTemp(outputDir, "youtube_*.m4v")
	if err != nil {
		return err
	}
	defer os.Remove(videoFile.Name())

	audioFile, err := os.CreateTemp(outputDir, "youtube_*.m4a")
	if err != nil {
		return err
	}
	defer os.Remove(audioFile.Name())

	log.Debug("Downloading video file...")
	err = dl.videoDLWorker(ctx, videoFile, v, videoFormat)
	if err != nil {
		return err
	}

	log.Debug("Downloading audio file...")
	err = dl.videoDLWorker(ctx, audioFile, v, audioFormat)
	if err != nil {
		return err
	}

	ffmpegVersionCmd := exec.Command("ffmpeg", "-y",
		"-i", videoFile.Name(),
		"-i", audioFile.Name(),
	)

	if reencode {
		ffmpegVersionCmd.Args = append(ffmpegVersionCmd.Args, "-c:v", "libx264", "-c:a", "aac")
	} else {
		ffmpegVersionCmd.Args = append(ffmpegVersionCmd.Args, "-c:v", "copy", "-c:a", "copy")
	}

	ffmpegVersionCmd.Args = append(ffmpegVersionCmd.Args, "-shortest", destFile, "-loglevel", "warning")

	ffmpegVersionCmd.Stderr = os.Stderr
	ffmpegVersionCmd.Stdout = os.Stdout
	log.Info("merging video and audio", "output", destFile)

	// Close files before running ffmpeg as it needs them.
	// Defer wouldn't work.
	videoFile.Close()
	audioFile.Close()

	return ffmpegVersionCmd.Run()
}

func GetVideoAudioFormats(v *youtube.Video, quality string, mimetype string, language string) (*youtube.Format, *youtube.Format, error) {
	formats := v.Formats

	if mimetype != "" {
		formats = formats.Type(mimetype)
	}

	videoFormats := formats.Type("video").AudioChannels(0)
	audioFormats := formats.Type("audio")

	if quality != "" {
		videoFormats = videoFormats.Quality(quality)
	}

	if language != "" {
		audioFormats = videoFormats.Language(language)
	}

	if len(videoFormats) == 0 {
		return nil, nil, errors.New("no video format found after filtering")
	}

	if len(audioFormats) == 0 {
		return nil, nil, errors.New("no audio format found after filtering")
	}

	videoFormats.Sort()
	audioFormats.Sort()

	return &videoFormats[0], &audioFormats[0], nil
}

func (dl *Downloader) videoDLWorker(ctx context.Context, out *os.File, video *youtube.Video, format *youtube.Format) error {
	stream, size, err := dl.GetStreamContext(ctx, video, format)
	if err != nil {
		return err
	}

	bar := progressbar.New(int(size))

	reader := progressbar.NewReader(stream, bar)

	mw := io.MultiWriter(out)

	_, err = io.Copy(mw, &reader)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	return nil
}

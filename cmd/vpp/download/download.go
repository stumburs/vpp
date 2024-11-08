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

type DownloadFlags struct {
	ReencodeAfterDownload bool
	DownloadAsMP3         bool
	DownloadAsWAV         bool
}

func (dl *Downloader) getOutputFile(v *youtube.Video, _ *youtube.Format, outputFile string, downloadFlags DownloadFlags) (string, error) {
	if outputFile == "" {
		outputFile = SanitizeFilename(v.Title)

		// "There's nothing more permanent than a temporary solution"
		// outputFile += pickIdealFileExtension(format.MimeType)
		if downloadFlags.DownloadAsMP3 {
			outputFile += ".mp3"
		} else if downloadFlags.DownloadAsWAV {
			outputFile += ".wav"
		} else {
			outputFile += ".mp4"
		}
	}

	if dl.OutputDir != "" {
		if err := os.MkdirAll(dl.OutputDir, 0o755); err != nil {
			return "", err
		}
		outputFile = filepath.Join(dl.OutputDir, outputFile)
	}

	return outputFile, nil
}

func (dl *Downloader) Download(ctx context.Context, v *youtube.Video, format *youtube.Format, outputFile string, downloadFlags DownloadFlags) error {
	youtube.Logger.Info(
		"Downloading video",
		"id", v.ID,
		"quality", format.Quality,
		"mimeType", format.MimeType,
	)

	destFile, err := dl.getOutputFile(v, format, outputFile, downloadFlags)
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

func (dl *Downloader) DownloadComposite(ctx context.Context, outputFile string, v *youtube.Video, quality string, mimetype, language string, downloadFlags DownloadFlags) error {
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

	destFile, err := dl.getOutputFile(v, videoFormat, outputFile, downloadFlags)
	if err != nil {
		return err
	}
	outputDir := filepath.Dir(destFile)

	// Download only as .mp3
	if downloadFlags.DownloadAsMP3 {
		return dl.downloadVideoAsMP3(ctx, outputDir, destFile, v, audioFormat)
	}

	// Download only as .wav
	if downloadFlags.DownloadAsWAV {
		return dl.downloadVideoAsWAV(ctx, outputDir, destFile, v, audioFormat)
	}

	return dl.downloadVideoAsMP4(ctx, outputDir, destFile, v, videoFormat, audioFormat, downloadFlags)
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

func combineIntoMP4(videoFile *os.File, audioFile *os.File, destFile string) error {
	ffmpegVersionCmd := exec.Command("ffmpeg", "-y",
		"-i", videoFile.Name(),
		"-i", audioFile.Name(),
		"-c:v", "copy",
		"-c:a", "copy",
		"-shortest", destFile,
		"-loglevel", "warning",
	)

	ffmpegVersionCmd.Stderr = os.Stderr
	ffmpegVersionCmd.Stdout = os.Stdout
	fmt.Printf("\nMerging temporary files into: %s\n", destFile)

	videoFile.Close()
	audioFile.Close()

	return ffmpegVersionCmd.Run()
}

func combineIntoMP4Reencode(videoFile *os.File, audioFile *os.File, destFile string) error {
	ffmpegVersionCmd := exec.Command("ffmpeg", "-y",
		"-i", videoFile.Name(),
		"-i", audioFile.Name(),
		"-c:v", "libx264",
		"-c:a", "aac",
		"-shortest", destFile,
		"-loglevel", "warning",
	)

	ffmpegVersionCmd.Stderr = os.Stderr
	ffmpegVersionCmd.Stdout = os.Stdout
	fmt.Printf("\nMerging temporary files into: %s\n", destFile)

	videoFile.Close()
	audioFile.Close()

	return ffmpegVersionCmd.Run()
}

func processIntoMP3(audioFile *os.File, destFile string) error {
	ffmpegVersionCmd := exec.Command("ffmpeg", "-y",
		"-i", audioFile.Name(),
		"-vn",
		"-ar", "44100",
		"-ac", "2",
		"-b:a", "192k",
		destFile, "-loglevel", "warning",
	)

	ffmpegVersionCmd.Stderr = os.Stderr
	ffmpegVersionCmd.Stdout = os.Stdout

	audioFile.Close()

	return ffmpegVersionCmd.Run()
}

func processIntoWAV(audioFile *os.File, destFile string) error {
	ffmpegVersionCmd := exec.Command("ffmpeg", "-y",
		"-i", audioFile.Name(),
		"-vn",
		"-ar", "44100",
		"-ac", "2",
		destFile, "-loglevel", "warning",
	)

	ffmpegVersionCmd.Stderr = os.Stderr
	ffmpegVersionCmd.Stdout = os.Stdout

	audioFile.Close()

	return ffmpegVersionCmd.Run()
}

func (dl *Downloader) downloadVideoAsMP3(ctx context.Context, outputDir, destFile string, v *youtube.Video, audioFormat *youtube.Format) error {
	audioFile, err := os.CreateTemp(outputDir, "youtube_*.m4a")
	if err != nil {
		return err
	}
	defer os.Remove(audioFile.Name())

	fmt.Printf("Downloading audio file...\n")
	if err := dl.videoDLWorker(ctx, audioFile, v, audioFormat); err != nil {
		return err
	}

	fmt.Printf("\nCreating .mp3 file: %s\n", destFile)
	return processIntoMP3(audioFile, destFile)
}

func (dl *Downloader) downloadVideoAsWAV(ctx context.Context, outputDir, destFile string, v *youtube.Video, audioFormat *youtube.Format) error {
	audioFile, err := os.CreateTemp(outputDir, "youtube_*.m4a")
	if err != nil {
		return err
	}
	defer os.Remove(audioFile.Name())

	fmt.Printf("Downloading audio file...\n")
	if err := dl.videoDLWorker(ctx, audioFile, v, audioFormat); err != nil {
		return err
	}

	fmt.Printf("\nCreating .wav file: %s\n", destFile)
	return processIntoWAV(audioFile, destFile)
}

func (dl *Downloader) downloadVideoAsMP4(ctx context.Context, outputDir, destFile string, v *youtube.Video, videoFormat *youtube.Format, audioFormat *youtube.Format, downloadFlags DownloadFlags) error {
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

	fmt.Printf("Downloading video file...\n")
	if err := dl.videoDLWorker(ctx, videoFile, v, videoFormat); err != nil {
		return err
	}
	videoFile.Close()

	fmt.Printf("Downloading audio file...\n")
	if err := dl.videoDLWorker(ctx, audioFile, v, audioFormat); err != nil {
		return err
	}
	audioFile.Close()

	// Combine into final output format
	if downloadFlags.ReencodeAfterDownload {
		return combineIntoMP4Reencode(videoFile, audioFile, destFile)
	}
	return combineIntoMP4(videoFile, audioFile, destFile)
}

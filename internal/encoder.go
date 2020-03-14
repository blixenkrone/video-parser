package internal

import (
	"errors"
	"io"
	"os/exec"
	"strings"
)

var (
	VideoFormatSuffix = []string{".mp4", ".mov"}
)

func SupportedSuffix(fileName string) bool {
	for _, suffix := range VideoFormatSuffix {
		return strings.HasSuffix(fileName, suffix)
	}
	return false
}

type video struct {
}

func NewVideo(rd io.Reader) *video {
	return &video{}
}

func (v *video) Read(p []byte) (int, error) {
	return 0, nil
}

func (v *video) RawMeta(r io.Reader) ([]byte, error) {
	ffprobe, err := exec.LookPath("ffprobe")
	if err != nil {
		return nil, errors.New("ffmpeg no bin in $PATH")
	}
	// cmd = exec.Command(ffprobe, "-v", "error", "-print_format", "json", "-show_format", "-show_streams", "-hide_banner", r)
	cmd := exec.Command(ffprobe, "-v", "quiet", "-print_format", "json", "-show_format", r)
	return cmd.Output()
}

func (v *video) RawMetaString(path string) ([]byte, error) {
	ffprobe, err := exec.LookPath("ffprobe")
	if err != nil {
		return nil, errors.New("ffmpeg no bin in $PATH")
	}
	// cmd = exec.Command(ffprobe, "-v", "error", "-print_format", "json", "-show_format", "-show_streams", "-hide_banner", r)
	cmd := exec.Command(ffprobe, "-v", "quiet", "-print_format", "json", "-show_format", path)
	return cmd.Output()
}

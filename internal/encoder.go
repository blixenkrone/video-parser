package internal

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

var (
	VideoFormatSuffix = []string{"mp4", "mov"}
)

func SupportedSuffix(fileName string) bool {
	fileSuffix := strings.Split(fileName, "/")[1]
	for _, suffix := range VideoFormatSuffix {
		return strings.HasSuffix(fileSuffix, suffix)
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

// ? maybe return []byte, error from cmd.Output()
func (v *video) RawMeta(r io.Reader) ([]byte, error) {
	ffprobe, err := exec.LookPath("ffprobe")
	if err != nil {
		return nil, errors.New("ffmpeg no bin in $PATH")
	}
	var out bytes.Buffer
	// cmd := exec.Command(ffprobe, "-v", "error", "-print_format", "json", "-show_format", "-show_streams", "-hide_banner", "pipe:0")
	cmd := exec.Command(ffprobe, "-v", "quiet", "-print_format", "json", "-show_format", "pipe:0")
	cmd.Stdin = r
	stderr, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	go io.Copy(&out, stderr)
	spew.Dump(cmd.String())
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
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

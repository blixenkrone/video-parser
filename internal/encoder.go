package internal

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

var (
	VideoFormatSuffix = []string{"mp4", "mov", "quicktime"}
)

func SupportedSuffix(fileName string) bool {
	fileSuffix := strings.Split(fileName, "/")[1]
	for _, suffix := range VideoFormatSuffix {
		if fileSuffix == suffix {
			return true
		}
	}
	return false
}

type video struct{}

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
	// cmd := exec.Command(ffprobe, "-v", "error", "-print_format", "json", "-show_format", "-show_streams", "-hide_banner", "pipe:0")

	f, err := ioutil.TempFile(os.TempDir(), "video-*")
	if err != nil {
		return nil, err
	}
	defer v.removeFile(f)

	go io.Copy(f, r)
	cmd := exec.Command(ffprobe, "-v", "quiet", "-print_format", "json", "-show_format", f.Name())
	spew.Dump(cmd.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
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

func (v *video) removeFile(f *os.File) error {
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Remove(f.Name()); err != nil {
		return err
	}
	return nil
}

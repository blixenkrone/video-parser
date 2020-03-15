package encoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	VideoFormatSuffix = []string{"mp4", "mov", "quicktime", "x-m4v", "m4v"}
	fromSecondMark    = "00:00:01.000"
	toSecondMark      = "00:00:01.100"
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

type VideoReader interface {
	io.Reader
}

func RawMeta(r VideoReader) (*FFMPEGMetaOutput, error) {
	ffprobe, err := exec.LookPath("ffprobe")
	if err != nil {
		return nil, errors.New("ffprobe no bin in $PATH")
	}
	cmd := exec.Command(ffprobe, "-v", "quiet", "-print_format", "json", "-show_format", "pipe:0")
	cmd.Stdin = r
	outJSON, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	var ffmpeg FFMPEGMetaOutput
	if err := json.Unmarshal(outJSON, &ffmpeg); err != nil {
		return nil, err
	}
	return &ffmpeg, nil
}

func (fo *FFMPEGMetaOutput) SanitizeOutput() *FFMPEGMetaOutput {
	fo.Format.Tags.ISOLocation = strings.Replace(fo.Format.Tags.ISOLocation, "+", "", 1)
	fo.Format.Tags.ISOLocation = strings.Replace(fo.Format.Tags.ISOLocation, "+", ",", 1)
	fo.Format.Tags.ISOLocation = strings.Replace(fo.Format.Tags.ISOLocation, "/", "", -1)
	return fo
}

// ffprobe output
type FFMPEGMetaOutput struct {
	Format struct {
		Filename       string `json:"filename"`
		NbStreams      int    `json:"nb_streams"`
		NbPrograms     int    `json:"nb_programs"`
		FormatName     string `json:"format_name"`
		FormatLongName string `json:"format_long_name"`
		StartTime      string `json:"start_time"`
		Duration       string `json:"duration"`
		ProbeScore     int    `json:"probe_score"`
		Tags           struct {
			MajorBrand                 string    `json:"major_brand"`
			MinorVersion               string    `json:"minor_version"`
			CompatibleBrands           string    `json:"compatible_brands"`
			CreationTime               time.Time `json:"creation_time"`
			ComAppleQuicktimeArtwork   string    `json:"com.apple.quicktime.artwork"`
			ComAppleQuicktimeIsMontage string    `json:"com.apple.quicktime.is-montage"`
			ISOLocation                string    `json:"com.apple.quicktime.location.ISO6709"`
		} `json:"tags"`
	} `json:"format"`
}

type FFMPEGThumbnail []byte

func Thumbnail(r VideoReader) ([]byte, error) {
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, errors.New("ffmpeg no bin in $PATH")
	}
	cmd := exec.Command(ffmpeg, "-v", "quiet", "pipe:", "-vframes", "1", "pipe:")
	cmd.Stdin = r
	return cmd.CombinedOutput()
}

func removeFile(f *os.File) error {
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Remove(f.Name()); err != nil {
		return err
	}
	return nil
}

func parseLocation(t time.Time) (time.Time, error) {
	tl, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		return t, err
	}
	fmt.Println(tl.String())
	return t.In(tl), nil
}

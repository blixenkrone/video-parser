package internal

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestEncoder(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		var v video
		f, _ := os.Open("../in.mp4")

		meta, err := v.RawMetaString(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		spew.Dump(meta)
	})
}

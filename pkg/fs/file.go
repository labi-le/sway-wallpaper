package fs

import (
	"crypto/sha256"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"os"
)

var ErrUnknownExtension = fmt.Errorf("unknown extension")

func SaveFile(readCloser io.ReadCloser, dir string) (string, error) {
	img, err := io.ReadAll(readCloser)
	if err != nil {
		return "", err
	}

	ext := mimetype.Detect(img).Extension()
	if ext == "" {
		return "", ErrUnknownExtension
	}

	if ioErr := os.MkdirAll(dir, 0755); ioErr != nil {
		return "", ioErr
	}

	sha := sha256.New()
	sha.Write(img)
	// first 7 characters
	short := fmt.Sprintf("%x", sha.Sum(nil))[:7]

	gen := fmt.Sprintf("%s/sw-%s%s", dir, short, ext)
	return gen, os.WriteFile(gen, img, 0600)
}

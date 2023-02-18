package tools

import (
	"crypto/sha256"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"os"
)

var ErrUnknownExtension = fmt.Errorf("unknown extension")

func SaveFile(image []byte, dir string) (string, error) {
	ext := mimetype.Detect(image).Extension()
	if ext == "" {
		return "", ErrUnknownExtension
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	sha := sha256.New()
	sha.Write(image)
	// first 7 characters
	short := fmt.Sprintf("%x", sha.Sum(nil))[:7]

	gen := fmt.Sprintf("%s/hw-%s%s", dir, short, ext)
	return gen, os.WriteFile(gen, image, 0600)
}

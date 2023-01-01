package tools

import (
	"crypto/sha256"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"os"
)

func SaveFile(image []byte, dir string) (string, error) {
	ext := mimetype.Detect(image).Extension()
	if ext == "" {
		return "", fmt.Errorf("unknown extension")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	sha := sha256.New()
	sha.Write(image)
	// first 7 characters
	short := fmt.Sprintf("%x", sha.Sum(nil))[:7]

	gen := fmt.Sprintf("%s/ghw-%s%s", dir, short, ext)
	return gen, os.WriteFile(gen, image, 0644)
}

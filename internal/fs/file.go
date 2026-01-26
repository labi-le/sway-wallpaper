package fs

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

var (
	ErrUnknownExtension = fmt.Errorf("unknown extension")
	unsafeChars         = regexp.MustCompile(`[^a-zA-Z0-9а-яА-Я]+`)
)

func SaveFile(data io.Reader, dir string, tags []string) (string, error) {
	img, err := io.ReadAll(data)
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
	short := fmt.Sprintf("%x", sha.Sum(nil))[:7]

	tagSuffix := buildTagSuffix(tags)
	gen := fmt.Sprintf("%s/sw-%s%s%s", dir, short, tagSuffix, ext)

	return gen, os.WriteFile(gen, img, 0600)
}

func buildTagSuffix(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	var validTags []string
	seen := make(map[string]struct{})

	for _, t := range tags {
		clean := unsafeChars.ReplaceAllString(t, "")
		clean = strings.ToLower(clean)

		if clean == "" || len(clean) > 20 {
			continue
		}

		if _, exists := seen[clean]; !exists {
			validTags = append(validTags, clean)
			seen[clean] = struct{}{}
		}
	}

	if len(validTags) == 0 {
		return ""
	}

	return "__" + strings.Join(validTags, "_")
}

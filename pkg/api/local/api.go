package local

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/labi-le/chiasma/pkg/api/searcher"
	"github.com/rs/zerolog"
)

const Name = "local"

type Local struct {
	dir string
	log zerolog.Logger
}

func NewLocal(log zerolog.Logger, dir string) *Local {
	return &Local{
		dir: dir,
		log: log.With().Str("component", "local_fs").Logger(),
	}
}

func (l *Local) Search(_ context.Context, q string, _ searcher.Resolution) (searcher.Image, error) {
	log := l.log.With().Str("op", "Search").Str("query", q).Logger()

	entries, err := os.ReadDir(l.dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", l.dir, err)
	}

	queryTerms := strings.Fields(strings.ToLower(q))
	var candidates []string

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := strings.ToLower(e.Name())
		matched := true
		for _, term := range queryTerms {
			if !strings.Contains(name, term) {
				matched = false
				break
			}
		}

		if matched {
			candidates = append(candidates, filepath.Join(l.dir, e.Name()))
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no local images found for query: %s", q)
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for _, path := range candidates {
		if err := l.validateImage(path); err == nil {
			f, err := os.Open(path)
			if err != nil {
				log.Warn().Err(err).Str("path", path).Msg("failed to open candidate")
				continue
			}

			img, err := searcher.DetectSize(f)
			if err != nil {
				_ = f.Close()
				log.Warn().Err(err).Str("path", path).Msg("failed to detect image size")
				continue
			}

			return img, nil
		}
	}

	return nil, fmt.Errorf("no valid images found among candidates")
}

func (l *Local) validateImage(path string) error {
	ext := filepath.Ext(path)
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".webp", ".bmp":
		return nil
	}

	mtype, err := mimetype.DetectFile(path)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(mtype.String(), "image/") {
		return fmt.Errorf("not an image")
	}
	return nil
}

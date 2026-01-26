package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/labi-le/chiasma/internal/fs"
	"github.com/labi-le/chiasma/internal/output"
	"github.com/labi-le/chiasma/pkg/api"
	"github.com/labi-le/chiasma/pkg/wallpaper/execute"
	"github.com/rs/zerolog"
)

type QuerySource interface {
	GetLastSearch() (string, error)
}

type WallpaperService struct {
	Log     zerolog.Logger
	API     api.Searcher
	History QuerySource
	Setter  execute.Provider
}

type UpdateParams struct {
	Phrase     string
	Resolution output.Resolution
	SaveDir    string
	OutputID   string
	RetryCount int
}

func (s *WallpaperService) Update(ctx context.Context, params UpdateParams) error {
	log := s.Log.With().Str("op", "Update").Logger()

	phrase := params.Phrase
	if phrase == "" {
		if s.History == nil {
			return errors.New("search phrase is empty and no history source provided")
		}
		var err error
		phrase, err = s.History.GetLastSearch()
		if err != nil {
			return fmt.Errorf("failed to get search phrase from history: %w", err)
		}
		log.Info().Msgf("using phrase from history: %s", phrase)
	}

	img, err := s.fetchImageWithRetry(ctx, phrase, params.Resolution, params.RetryCount)
	if err != nil {
		return err
	}
	defer img.Close()

	path, err := fs.SaveFile(img, params.SaveDir)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}
	log.Debug().Str("path", path).Msg("image saved")

	if err := s.Setter.Change(ctx, path, params.OutputID); err != nil {
		return fmt.Errorf("failed to set wallpaper: %w", err)
	}

	return nil
}

func (s *WallpaperService) fetchImageWithRetry(ctx context.Context, phrase string, res output.Resolution, retries int) (api.Image, error) {
	var lastErr error
	for i := 0; i < retries; i++ {
		img, err := s.API.Search(ctx, phrase, res)
		if err != nil {
			lastErr = err
			continue
		}

		w, h := img.Size()
		if w < res.Width || h < res.Height {
			_ = img.Close()
			lastErr = fmt.Errorf("image too small: %dx%d < %dx%d", w, h, res.Width, res.Height)
			continue
		}

		return img, nil
	}
	return nil, fmt.Errorf("failed to find suitable image after %d attempts: %w", retries, lastErr)
}

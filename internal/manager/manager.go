package manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/labi-le/sway-wallpaper/pkg/browser"
	"github.com/labi-le/sway-wallpaper/pkg/fs"
	"github.com/labi-le/sway-wallpaper/pkg/wallpaper"
	"github.com/rs/zerolog/log"
)

type Manager struct {
	Options       Options
	WallpaperAPI  Searcher
	Analyzer      Analyzer
	WallpaperTool wallpaper.Tool
}

func New(opt Options) (*Manager, error) {
	if opt.SearchPhrase != "" {
		opt.Browser = browser.NoopBrowser
	}

	tool, err := wallpaper.ByName(opt.Tool)
	if err != nil {
		return nil, fmt.Errorf("failed to found wallpaper tool provider: %v", err)
	}
	return &Manager{
		WallpaperAPI:  MustSearcher(opt.API),
		WallpaperTool: tool,
		Analyzer:      NewBrowserHistoryAnalyzer(opt.Browser, opt.HistoryFile),
		Options:       opt,
	}, nil
}

func (m *Manager) Provide(ctx context.Context) error {
	if m.Options.SearchPhrase == "" {
		log.Trace().Msg("searching phrase not provided, trying to get last searched phrase from browser")

		var searchPhErr error
		m.Options.SearchPhrase, searchPhErr = m.Analyzer.Analyze()
		if searchPhErr != nil {
			return searchPhErr
		}
	}

	log.Trace().Msgf("search for %s", m.Options.SearchPhrase)

	img, searchErr := m.WallpaperAPI.Search(ctx, m.Options.SearchPhrase, m.Options.ImageResolution)
	if searchErr != nil {
		return searchErr
	}

	defer img.Close()

	log.Trace().Msgf("save wallpaper to %s", m.Options.SaveWallpaperPath)

	path, saveErr := fs.SaveFile(img, m.Options.SaveWallpaperPath)
	if saveErr != nil {
		return saveErr
	}
	log.Trace().Msgf("provide wallpaper from %s", path)

	if err := m.WallpaperTool.Change(ctx, path, m.Options.Output.ID); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (m *Manager) Close() error {
	if closer, ok := m.WallpaperTool.(wallpaper.ToolDestructor); ok {
		return closer.Close()
	}
	return nil
}

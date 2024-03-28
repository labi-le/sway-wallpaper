package manager

import (
	"context"
	"errors"
	"github.com/labi-le/sway-wallpaper/pkg/browser"
	"github.com/labi-le/sway-wallpaper/pkg/fs"
	"github.com/labi-le/sway-wallpaper/pkg/log"
)

type Manager struct {
	Options       Options
	WallpaperAPI  Searcher
	Analyzer      Analyzer
	WallpaperTool Setter
}

func New(opt Options) *Manager {
	if opt.SearchPhrase != "" {
		opt.Browser = browser.NoopBrowser
	}
	return &Manager{
		WallpaperAPI:  MustSearcher(opt.API),
		WallpaperTool: MustBGTool(opt.Tool),
		Analyzer:      NewBrowserHistoryAnalyzer(opt.Browser, opt.HistoryFile),
		Options:       opt,
	}
}

func (m *Manager) Provide(ctx context.Context) error {
	if m.Options.SearchPhrase == "" {
		log.Info("Searching phrase not provided, trying to get last searched phrase from browser")
		var searchPhErr error
		m.Options.SearchPhrase, searchPhErr = m.Analyzer.Analyze()
		if searchPhErr != nil {
			return searchPhErr
		}
	}

	log.Infof("Search for %s", m.Options.SearchPhrase)
	img, searchErr := m.WallpaperAPI.Search(ctx, m.Options.SearchPhrase, m.Options.ImageResolution)
	if searchErr != nil {
		return searchErr
	}

	defer img.Close()

	log.Infof("Save wallpaper to %s", m.Options.SaveWallpaperPath)
	path, saveErr := fs.SaveFile(img, m.Options.SaveWallpaperPath)
	if saveErr != nil {
		return saveErr
	}

	log.Infof("Provide wallpaper from %s", path)
	if err := m.WallpaperTool.Set(ctx, path, m.Options.Output.ID); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (m *Manager) Close() error {
	return m.WallpaperTool.Clean()
}

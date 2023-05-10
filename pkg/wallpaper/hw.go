package wallpaper

import (
	"context"
	"errors"
	"github.com/labi-le/sway-wallpaper/pkg/api"
	"github.com/labi-le/sway-wallpaper/pkg/browser"
	"github.com/labi-le/sway-wallpaper/pkg/fs"
	"github.com/labi-le/sway-wallpaper/pkg/log"
	"github.com/labi-le/sway-wallpaper/pkg/wptool"
	"os/user"
	"time"
)

type SwayWallpaper struct {
	Options       Options
	WallpaperAPI  api.Finder
	Browser       browser.PhraseFinder
	WallpaperTool wptool.Setter
}

type Options struct {
	SearchPhrase      string
	SaveWallpaperPath string
	FollowDuration    time.Duration
	Resolution        api.Resolution
	API               string
	Tool              string
	HistoryFile       string
	Browser           string
	Usr               *user.User
}

func Must(opt Options) *SwayWallpaper {
	if opt.SearchPhrase != "" {
		opt.Browser = browser.Noop
	}
	return &SwayWallpaper{
		WallpaperAPI:  api.MustFinder(opt.API),
		WallpaperTool: wptool.ParseTool(opt.Tool),
		Browser:       browser.MustBrowser(opt.Browser, opt.Usr, opt.HistoryFile),
		Options:       opt,
	}
}

func (sw *SwayWallpaper) Set(ctx context.Context) error {
	if sw.Options.SearchPhrase == "" {
		log.Info("Searching phrase not provided, trying to get last searched phrase from browser")
		var searchPhErr error
		sw.Options.SearchPhrase, searchPhErr = sw.Browser.LastSearchedPhrase()
		if searchPhErr != nil {
			return searchPhErr
		}
	}

	log.Infof("Search for %s", sw.Options.SearchPhrase)
	img, searchErr := sw.WallpaperAPI.Find(ctx, sw.Options.SearchPhrase, sw.Options.Resolution)
	if searchErr != nil {
		return searchErr
	}

	defer img.Close()

	log.Infof("Save wallpaper to %s", sw.Options.SaveWallpaperPath)
	path, saveErr := fs.SaveFile(img, sw.Options.SaveWallpaperPath)
	if saveErr != nil {
		return saveErr
	}

	log.Infof("Set wallpaper from %s", path)
	if err := sw.WallpaperTool.Set(ctx, path); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

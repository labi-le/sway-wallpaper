package wallpaper

import (
	"context"
	"errors"
	"github.com/labi-le/history-wallpaper/pkg/api"
	"github.com/labi-le/history-wallpaper/pkg/browser"
	"github.com/labi-le/history-wallpaper/pkg/fs"
	"github.com/labi-le/history-wallpaper/pkg/log"
	"github.com/labi-le/history-wallpaper/pkg/wptool"
	"os/user"
	"time"
)

type HW struct {
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

func MustHW(opt Options) *HW {
	return &HW{
		WallpaperAPI:  api.MustFinder(opt.API),
		WallpaperTool: wptool.ParseTool(opt.Tool),
		Browser:       browser.MustBrowser(opt.Browser, opt.Usr, opt.HistoryFile),
		Options:       opt,
	}
}

func (hw *HW) Set(ctx context.Context) error {
	log.Info("Setting wallpaper...")
	if hw.Options.SearchPhrase == "" {
		log.Info("Searching phrase not provided, trying to get last searched phrase from browser")
		var searchPhErr error
		hw.Options.SearchPhrase, searchPhErr = hw.Browser.LastSearchedPhrase()
		if searchPhErr != nil {
			return searchPhErr
		}
	}

	log.Infof("Searching for %s", hw.Options.SearchPhrase)
	img, searchErr := hw.WallpaperAPI.Find(ctx, hw.Options.SearchPhrase, hw.Options.Resolution)
	if searchErr != nil {
		return searchErr
	}

	defer img.Close()

	log.Infof("Saving wallpaper to %s", hw.Options.SaveWallpaperPath)
	path, saveErr := fs.SaveFile(img, hw.Options.SaveWallpaperPath)
	if saveErr != nil {
		return saveErr
	}

	log.Infof("Setting wallpaper to %s", path)
	if err := hw.WallpaperTool.Set(ctx, path); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

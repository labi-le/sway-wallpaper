package wallpaper

import (
	"context"
	"errors"
	"github.com/bearatol/lg"
	"github.com/labi-le/history-wallpaper/pkg/api"
	"github.com/labi-le/history-wallpaper/pkg/browser"
	"github.com/labi-le/history-wallpaper/pkg/fs"
	"github.com/labi-le/history-wallpaper/pkg/wptool"
	"os/user"
	"time"
)

var (
	ErrPhrase = errors.New("")
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
	lg.Info("Setting wallpaper...")
	if hw.Options.SearchPhrase == "" {
		lg.Info("Searching phrase not provided, trying to get last searched phrase from browser")
		var searchPhErr error
		hw.Options.SearchPhrase, searchPhErr = hw.Browser.LastSearchedPhrase()
		if searchPhErr != nil {
			return searchPhErr
		}
	}

	lg.Infof("Searching for %s", hw.Options.SearchPhrase)
	img, searchErr := hw.WallpaperAPI.Find(ctx, hw.Options.SearchPhrase, hw.Options.Resolution)
	if searchErr != nil {
		return searchErr
	}

	defer img.Close()

	lg.Infof("Saving wallpaper to %s", hw.Options.SaveWallpaperPath)
	path, saveErr := fs.SaveFile(img, hw.Options.SaveWallpaperPath)
	if saveErr != nil {
		return saveErr
	}

	lg.Infof("Setting wallpaper to %s", path)
	if err := hw.WallpaperTool.Set(ctx, path); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

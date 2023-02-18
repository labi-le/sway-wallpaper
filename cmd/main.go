package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/labi-le/history-wallpaper/pkg/api"
	"github.com/labi-le/history-wallpaper/pkg/browser"
	"github.com/labi-le/history-wallpaper/pkg/fs"
	"github.com/labi-le/history-wallpaper/pkg/wallpaper"
	"github.com/labi-le/history-wallpaper/pkg/wptool"
	"github.com/nightlyone/lockfile"
	"golang.org/x/net/context"
	_ "modernc.org/sqlite"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"time"
)

func main() {
	lock := MustLock()
	defer Unlock(lock)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	usr, err := user.Current()
	if err != nil {
		Error(err)
	}

	opt := Parse(browser.Available(), usr)

	if opt.FollowDuration == 0 {
		tick(ctx, opt)
		<-ctx.Done()
		return
	}

	for {
		reuse, currentCancel := context.WithTimeout(ctx, opt.FollowDuration)

		tick(reuse, opt)

		<-reuse.Done()
		currentCancel()
	}
}

func Parse(availBrowsers []string, usr *user.User) wallpaper.HW {
	var (
		browserName  string
		historyFile  string
		resolution   string
		wpTool       string
		wpAPI        string
		saveImageDir string

		searchPhrase string

		follow         string
		followDuration time.Duration
	)

	flag.StringVar(
		&browserName,
		"browserName",
		availBrowsers[0],
		"browserName to use. Available: "+fmt.Sprint(availBrowsers),
	)
	flag.StringVar(&historyFile, "history-file", "",
		"browserName history file to use. Auto detect if empty (only for chromium based browsers)\n"+
			"e.g ~/.mozilla/icecat/gxda4hpz.default-1672760493248/formhistory.sqlite",
	)
	flag.StringVar(&resolution, "resolution", "1920x1080", "resolution to use. e.g. 1920x1080")
	flag.StringVar(&wpTool,
		"wp-tool", wptool.Available()[0], "wallpaper tool to use. Available: "+fmt.Sprint(wptool.Available()))
	flag.StringVar(&wpAPI,
		"wp-api", api.Available()[0], "wallpaper api to use. Available: "+fmt.Sprint(api.Available()))
	flag.StringVar(&saveImageDir, "save-image-dir", usr.HomeDir+"/Pictures", "directory to save image to")
	flag.StringVar(&searchPhrase, "search-phrase", "", "search phrase to use")
	flag.StringVar(&follow, "follow", "", "follow a time interval and update wallpaper. e.g. 1h, 1m, 30s")

	flag.Parse()

	if !checkAvailable(browserName, availBrowsers) {
		Error("Invalid browserName")
	}

	if !checkAvailable(wpTool, wptool.Available()) {
		Error("Invalid wallpaper tool")
	}

	if !checkAvailable(wpAPI, api.Available()) {
		Error("Invalid wallpaper api")
	}

	if follow != "" {
		var parseErr error
		followDuration, parseErr = time.ParseDuration(follow)
		if parseErr != nil {
			Error("Invalid follow. e.g. 1h, 1m, 1s")
		}
	}

	return wallpaper.HW{
		WallpaperAPI:      api.MustFinder(wpAPI),
		WallpaperTool:     wptool.ParseTool(wpTool),
		Browser:           browser.MustBrowser(browserName, usr, historyFile),
		Resolution:        api.Resolution(resolution),
		SaveWallpaperPath: saveImageDir,
		SearchPhrase:      searchPhrase,
		FollowDuration:    followDuration,
	}
}

func MustLock() lockfile.Lockfile {
	lock, _ := lockfile.New(filepath.Join(os.TempDir(), "hw.lck"))

	if lockErr := lock.TryLock(); lockErr != nil {
		owner, err := lock.GetOwner()
		if err != nil {
			Error(errors.New("cannot get locked process: " + lockErr.Error()))
		}
		Error(fmt.Errorf("hw is already running. pid %d", owner.Pid))
	}

	return lock
}

func Unlock(lock lockfile.Lockfile) {
	if err := lock.Unlock(); err != nil {
		Error(fmt.Errorf("cannot unlock %q, reason: %w", lock, err))
	}
}

func tick(ctx context.Context, opt wallpaper.HW) {
	if opt.SearchPhrase == "" {
		var searchPhErr error
		opt.SearchPhrase, searchPhErr = SearchedPhraseBrowser(opt.Browser)
		if searchPhErr != nil && !errors.Is(searchPhErr, context.Canceled) {
			Error("Error while getting searched phrase: " + searchPhErr.Error())
		}
	}

	Info("Search phrase: " + opt.SearchPhrase)

	img, searchErr := opt.WallpaperAPI.Find(ctx, opt.SearchPhrase, opt.Resolution)
	if searchErr != nil {
		Error("Error while getting img: " + searchErr.Error())
	}

	defer img.Close()

	path, saveErr := fs.SaveFile(img, opt.SaveWallpaperPath)
	if saveErr != nil {
		Error("Error while saving img: " + saveErr.Error())
	}

	Info("Saved img to: " + path)

	if err := opt.WallpaperTool.Set(ctx, path); err != nil && !errors.Is(err, context.Canceled) {
		Error("Error while setting wallpaper: " + err.Error())
	}
}

func Error(v any) {
	//nolint:forbidigo //dn
	fmt.Printf("%v\n", v)
	os.Exit(1)
}

func Info(v any) {
	//nolint:forbidigo //dn
	fmt.Printf("%v\n", v)
}

func SearchedPhraseBrowser(b browser.PhraseFinder) (string, error) {
	return b.LastSearchedPhrase()
}

func checkAvailable(concrete string, available []string) bool {
	for _, item := range available {
		if concrete == item {
			return true
		}
	}

	return false
}

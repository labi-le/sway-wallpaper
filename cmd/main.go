package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bearatol/lg"
	"github.com/labi-le/history-wallpaper/pkg/api"
	"github.com/labi-le/history-wallpaper/pkg/browser"
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
		lg.Error(err)
	}

	opt := Parse(api.Available(), wptool.Available(), browser.Available(), usr)
	hw := wallpaper.MustHW(opt)

	if opt.FollowDuration == 0 {
		if err := hw.Set(ctx); err != nil {
			lg.Error(err)
		}
		<-ctx.Done()
		return
	}

	for {
		reuse, currentCancel := context.WithTimeout(ctx, opt.FollowDuration)

		if err := hw.Set(reuse); err != nil {
			if errors.Is(err, context.Canceled) {
				lg.Warnf("Handling interrupt signal")
				break
			}
			lg.Error(err)
		}

		<-reuse.Done()
		currentCancel()
	}
}

func Parse(apiAvail []string, wpToolAvail []string, availBrowsers []string, usr *user.User) wallpaper.Options {
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
		"wp-tool", wpToolAvail[0], "wallpaper tool to use. Available: "+fmt.Sprint(wpToolAvail))
	flag.StringVar(&wpAPI,
		"wp-api", apiAvail[0], "wallpaper api to use. Available: "+fmt.Sprint(apiAvail))
	flag.StringVar(&saveImageDir, "save-image-dir", usr.HomeDir+"/Pictures", "directory to save image to")
	flag.StringVar(&searchPhrase, "search-phrase", "", "search phrase to use")
	flag.StringVar(&follow, "follow", "", "follow a time interval and update wallpaper. e.g. 1h, 1m, 30s")

	flag.Parse()

	if !checkAvailable(browserName, availBrowsers) {
		lg.Error("Invalid browserName")
	}

	if !checkAvailable(wpTool, wpToolAvail) {
		lg.Error("Invalid wallpaper tool")
	}

	if !checkAvailable(wpAPI, apiAvail) {
		lg.Error("Invalid wallpaper api")
	}

	if follow != "" {
		var parseErr error
		followDuration, parseErr = time.ParseDuration(follow)
		if parseErr != nil {
			lg.Error("Invalid follow. e.g. 1h, 1m, 1s")
		}
	}

	return wallpaper.Options{
		SearchPhrase:      searchPhrase,
		SaveWallpaperPath: saveImageDir,
		FollowDuration:    followDuration,
		Resolution:        api.Resolution(resolution),
		API:               wpAPI,
		Tool:              wpTool,
		HistoryFile:       historyFile,
		Browser:           browserName,
		Usr:               usr,
	}
}

func MustLock() lockfile.Lockfile {
	lock, _ := lockfile.New(filepath.Join(os.TempDir(), "hw.lck"))

	if lockErr := lock.TryLock(); lockErr != nil {
		owner, err := lock.GetOwner()
		if err != nil {
			lg.Error(errors.New("cannot get locked process: " + lockErr.Error()))
		}
		lg.Error(fmt.Errorf("hw is already running. pid %d", owner.Pid))
	}

	return lock
}

func Unlock(lock lockfile.Lockfile) {
	if err := lock.Unlock(); err != nil {
		lg.Error(fmt.Errorf("cannot unlock %q, reason: %w", lock, err))
	}
}

func checkAvailable(concrete string, available []string) bool {
	for _, item := range available {
		if concrete == item {
			return true
		}
	}

	return false
}

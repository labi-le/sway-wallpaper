package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/labi-le/sway-wallpaper/pkg/api"
	"github.com/labi-le/sway-wallpaper/pkg/browser"
	"github.com/labi-le/sway-wallpaper/pkg/log"
	"github.com/labi-le/sway-wallpaper/pkg/wallpaper"
	"github.com/labi-le/sway-wallpaper/pkg/wptool"
	"github.com/nightlyone/lockfile"
	"github.com/vcraescu/go-xrandr"
	"golang.org/x/net/context"
	_ "modernc.org/sqlite"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"regexp"
	"time"
)

var (
	ErrBrowserNotImplemented       = errors.New("invalid browser not implemented")
	ErrWallpaperToolNotImplemented = errors.New("invalid wallpaper tool not implemented")
	ErrWallpaperAPINotImplemented  = errors.New("invalid wallpaper api not implemented")
	ErrInvalidFollowDuration       = errors.New("invalid follow. e.g. 1h, 1m, 1s")
	ErrSwayWallpaperAlreadyRunning = errors.New("sway-wallpaper is already running. pid %d")
	ErrCannotLock                  = errors.New("cannot get locked process: %s")
	ErrInvalidResolution           = errors.New("invalid resolution. e.g. 1920x1080")
	ErrCannotUnlock                = errors.New("cannot unlock process: %s")
	ErrAutoResolutionNotSupported  = errors.New("xrandr not found. Please install xrandr or set resolution manually")
)

const LockFile = "sway-wallpaper.lck"

func main() {
	lock := MustLock()
	defer Unlock(lock)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	usr, err := user.Current()
	if err != nil {
		log.Error(err)
	}

	opt := Parse(api.Available(), wptool.Available(), browser.Available(), usr)
	wp := wallpaper.Must(opt)

	if opt.FollowDuration == 0 {
		if wpErr := wp.Set(ctx); wpErr != nil {
			if errors.Is(err, context.Canceled) {
				log.Warnf("Handling interrupt signal")
				return
			}
			log.Error(wpErr)
		}
		<-ctx.Done()
		return
	}

	for {
		reuse, currentCancel := context.WithTimeout(ctx, opt.FollowDuration)

		if wpErr := wp.Set(reuse); wpErr != nil {
			if errors.Is(wpErr, context.Canceled) {
				log.Warnf("Handling interrupt signal")
				break
			}
			log.Error(wpErr)
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
	flag.StringVar(&resolution, "resolution", DetectResolution(), "resolution to use. e.g. 1920x1080")
	flag.StringVar(&wpTool,
		"wp-tool", wpToolAvail[0], "wallpaper tool to use. Available: "+fmt.Sprint(wpToolAvail))
	flag.StringVar(&wpAPI,
		"wp-api", apiAvail[0], "wallpaper api to use. Available: "+fmt.Sprint(apiAvail))
	flag.StringVar(&saveImageDir, "save-image-dir", usr.HomeDir+"/Pictures", "directory to save image to")
	flag.StringVar(&searchPhrase, "search-phrase", "", "search phrase to use")
	flag.StringVar(&follow, "follow", "", "follow a time interval and update wallpaper. e.g. 1h, 1m, 30s")

	flag.Parse()

	if !checkAvailable(browserName, availBrowsers) && searchPhrase == "" {
		log.Fatal(ErrBrowserNotImplemented)
	}

	if !checkAvailable(wpTool, wpToolAvail) {
		log.Fatal(ErrWallpaperToolNotImplemented)
	}

	if !checkAvailable(wpAPI, apiAvail) {
		log.Fatal(ErrWallpaperAPINotImplemented)
	}

	if !validateResolution(resolution) {
		log.Fatal(ErrInvalidResolution)
	}

	if follow != "" {
		var parseErr error
		followDuration, parseErr = time.ParseDuration(follow)
		if parseErr != nil {
			log.Fatal(ErrInvalidFollowDuration)
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

func validateResolution(resolution string) bool {
	return regexp.MustCompile(`\d+x\d+`).MatchString(resolution)
}

func DetectResolution() string {
	screens, err := xrandr.GetScreens()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			log.Fatal(ErrAutoResolutionNotSupported)
		}
		log.Fatal(err)
	}

	screen := screens[0]

	resolution := fmt.Sprintf("%.fx%.f", screen.CurrentResolution.Width, screen.CurrentResolution.Height)
	log.Infof("Detected resolution: %s", resolution)

	return resolution
}

func MustLock() lockfile.Lockfile {
	lock, _ := lockfile.New(filepath.Join(os.TempDir(), LockFile))

	if lockErr := lock.TryLock(); lockErr != nil {
		owner, err := lock.GetOwner()
		if err != nil {
			log.Fatalf(ErrCannotLock.Error(), err)
		}
		log.Fatalf(ErrSwayWallpaperAlreadyRunning.Error(), owner.Pid)
	}

	return lock
}

func Unlock(lock lockfile.Lockfile) {
	if err := lock.Unlock(); err != nil {
		log.Fatalf(ErrCannotUnlock.Error(), err)
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

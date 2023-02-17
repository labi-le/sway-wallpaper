package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/labi-le/google-history-wallpaper/pkg/browser/chromium"
	"github.com/labi-le/google-history-wallpaper/pkg/browser/firefox"
	"github.com/labi-le/google-history-wallpaper/pkg/image/unsplash"
	"github.com/labi-le/google-history-wallpaper/pkg/tools"
	"github.com/labi-le/google-history-wallpaper/pkg/wallpaper"
	"github.com/nightlyone/lockfile"
	"golang.org/x/net/context"
	_ "modernc.org/sqlite"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"time"
)

var (
	wallpaperTools = []string{"swaybg", "wbg"}
	wallpaperAPI   = []string{"unsplash"}
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var (
		browser      string
		historyFile  string
		resolution   string
		wpTool       string
		wpAPI        string
		saveImageDir string

		searchPhrase string

		follow         string
		followDuration time.Duration
	)

	usr, usErr := user.Current()
	if usErr != nil {
		Error(usErr)
	}

	flag.StringVar(
		&browser,
		"browser",
		tools.GetAllBrowsers()[0],
		"browser to use. Available: "+fmt.Sprint(tools.GetAllBrowsers()),
	)
	flag.StringVar(&historyFile, "history-file", "",
		"browser history file to use. Auto detect if empty (only for chromium based browsers)\n"+
			"e.g ~/.mozilla/icecat/gxda4hpz.default-1672760493248/formhistory.sqlite",
	)
	flag.StringVar(&resolution, "resolution", "1920x1080", "resolution to use. e.g. 1920x1080")
	flag.StringVar(&wpTool, "wp-tool", wallpaperTools[0], "wallpaper tool to use. Available: "+fmt.Sprint(wallpaperTools))
	flag.StringVar(&wpAPI, "wp-api", wallpaperAPI[0], "wallpaper api to use. Available: "+fmt.Sprint(wallpaperAPI))
	flag.StringVar(&saveImageDir, "save-image-dir", usr.HomeDir+"/Pictures", "directory to save image to")
	flag.StringVar(&searchPhrase, "search-phrase", "", "search phrase to use")
	flag.StringVar(&follow, "follow", "", "follow a time interval and update wallpaper. e.g. 1h, 1m, 30s")

	flag.Parse()

	lock := MustLock()
	defer Unlock(lock)

	if !checkAvailable(browser, tools.GetAllBrowsers()) {
		Error("Invalid browser")
	}

	if !checkAvailable(wpTool, wallpaperTools) {
		Error("Invalid wallpaper tool")
	}

	if !checkAvailable(wpAPI, wallpaperAPI) {
		Error("Invalid wallpaper api")
	}

	if follow != "" {
		var parseErr error
		followDuration, parseErr = parseFollow(follow)
		if parseErr != nil {
			Error("Invalid follow. e.g. 1h, 1m, 1s")
		}
	}

	opt := wallpaper.Options{
		WallpaperAPI:      wpAPI,
		WallpaperSetter:   wpTool,
		Browser:           browser,
		HistoryFile:       historyFile,
		SaveWallpaperPath: saveImageDir,
		Resolution:        resolution,
		SearchPhrase:      searchPhrase,
	}

	for {
		if followDuration == 0 {
			tick(ctx, usr, opt)
			break
		}

		reuse, currentCancel := context.WithTimeout(ctx, followDuration)

		tick(reuse, usr, opt)

		<-reuse.Done()
		currentCancel()
	}
}

func MustLock() lockfile.Lockfile {
	lock, _ := lockfile.New(filepath.Join(os.TempDir(), "hw.lck"))

	if lockErr := lock.TryLock(); lockErr != nil {
		Error(fmt.Errorf("cannot unlock %q, reason: %w", lock, lockErr))
	}

	return lock
}

func Unlock(lock lockfile.Lockfile) {
	if err := lock.Unlock(); err != nil {
		Error(fmt.Errorf("cannot unlock %q, reason: %w", lock, err))
	}
}

func tick(ctx context.Context, usr *user.User, opt wallpaper.Options) {
	if opt.SearchPhrase == "" {
		var searchPhErr error
		opt.SearchPhrase, searchPhErr = SearchedPhraseBrowser(usr, opt.Browser, opt.HistoryFile)
		if errors.Is(searchPhErr, sql.ErrNoRows) {
			searchPhErr = errors.New("browser history is empty")
		}
		if searchPhErr != nil {
			Error("Error while getting searched phrase: " + searchPhErr.Error())
		}
	}

	Info("Search phrase: " + opt.SearchPhrase)

	image, searchErr := GetImage(opt.SearchPhrase, opt.WallpaperAPI, opt.Resolution)
	if searchErr != nil {
		Error("Error while getting image: " + searchErr.Error())
	}

	path, saveErr := tools.SaveFile(image, opt.SaveWallpaperPath)
	if saveErr != nil {
		Error("Error while saving image: " + saveErr.Error())
	}

	Info("Saved image to: " + path)

	if err := SetWallpaper(ctx, path, opt.WallpaperSetter); err != nil {
		Error("Error while setting wallpaper: " + err.Error())
	}
}

func parseFollow(f string) (time.Duration, error) {
	return time.ParseDuration(f)
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

func GetImage(phrase, service, resolution string) ([]byte, error) {
	if service == "unsplash" {
		return unsplash.GetImage(phrase, resolution)
	}

	return nil, nil
}

func SearchedPhraseBrowser(usr *user.User, browser string, file string) (string, error) {
	if tools.IsChromiumBased(browser) {
		return chromium.GetLastSearchedPhrase(usr, browser, file)
	}

	return firefox.GetLastSearchedPhrase(file)
}

func SetWallpaper(ctx context.Context, path string, tool string) error {
	switch tool {
	case "swaybg":
		return tools.SetWallpaperSwayBG(ctx, path)
	case "wbg":
		return tools.SetWallpaperWBG(ctx, path)
	}

	return nil
}

func checkAvailable(concrete string, available []string) bool {
	for _, item := range available {
		if concrete == item {
			return true
		}
	}

	return false
}

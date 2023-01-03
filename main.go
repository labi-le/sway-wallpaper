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
	_ "modernc.org/sqlite"
	"os"
	"os/user"
	"time"
)

var (
	chromiumBasedBrowsers = []string{"vivaldi", "chrome", "chromium", "brave", "opera"}
	firefoxBasedBrowsers  = []string{"firefox"}

	browsers = append(chromiumBasedBrowsers, firefoxBasedBrowsers...)

	wallpaperTools = []string{"swaybg", "wbg"}
	wallpaperAPI   = []string{"unsplash"}
)

func main() {
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

	flag.StringVar(&browser, "browser", browsers[0], "browser to use. Available: "+fmt.Sprint(browsers))
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

	if !checkAvailable(browser, browsers) {
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

	for {
		if followDuration == 0 {
			tick(usr, wpAPI, wpTool, browser, historyFile, saveImageDir, resolution, searchPhrase)
			break
		}

		tick(usr, wpAPI, wpTool, browser, historyFile, saveImageDir, resolution, searchPhrase)
		time.Sleep(followDuration)
	}
}

func tick(usr *user.User, wpAPI, wpTool, browser, historyFile, saveImageDir, resolution, searchPhrase string) {
	if searchPhrase == "" {
		var searchPhErr error
		searchPhrase, searchPhErr = SearchedPhraseBrowser(usr, browser, historyFile)
		if errors.Is(searchPhErr, sql.ErrNoRows) {
			searchPhErr = errors.New("browser history is empty")
		}
		if searchPhErr != nil {
			Error("Error while getting searched phrase: " + searchPhErr.Error())
		}
	}

	Info("Search phrase: " + searchPhrase)

	image, searchErr := GetImage(searchPhrase, wpAPI, resolution)
	if searchErr != nil {
		Error("Error while getting image: " + searchErr.Error())
	}

	path, saveErr := tools.SaveFile(image, saveImageDir)
	if saveErr != nil {
		Error("Error while saving image: " + saveErr.Error())
	}

	Info("Saved image to: " + path)

	if err := SetWallpaper(path, wpTool); err != nil {
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
	if isChromiumBased(browser) {
		return chromium.GetLastSearchedPhrase(usr, browser, file)
	}

	return firefox.GetLastSearchedPhrase(file)
}

func SetWallpaper(path string, tool string) error {
	switch tool {
	case "swaybg":
		return tools.SetWallpaperSwayBG(path)
	case "wbg":
		return tools.SetWallpaperWBG(path)
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

func isChromiumBased(browser string) bool {
	for _, b := range chromiumBasedBrowsers {
		if browser == b {
			return true
		}
	}

	return false
}

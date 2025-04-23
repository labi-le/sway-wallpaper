package manager

import (
	"flag"
	"fmt"
	"github.com/labi-le/sway-wallpaper/pkg/wallpaper"
	"os"
	"time"
)

func ParseOptions() (Options, error) {
	var (
		opt Options
	)

	flag.StringVar(
		&opt.Browser,
		"browserName",
		AvailableBrowsers()[0],
		"Available: "+fmt.Sprint(AvailableBrowsers()),
	)
	flag.StringVar(&opt.HistoryFile, "history-file", "",
		"path to history file. Auto-detect if empty (only for chromium based browsers)\n"+
			"e.g ~/.mozilla/icecat/gxda4hpz.default-1672760493248/formhistory.sqlite",
	)
	flag.Var(&opt.ImageResolution, "image-resolution", "image resolution. e.g. 1920x1080")
	flag.Var(&opt.Output, "output", "output to operate on. e.g eDP-1")
	flag.StringVar(&opt.Tool,
		"tool", wallpaper.AvailableProvider.String(), "wallpaper tool to use. Available: "+fmt.Sprint(wallpaper.AvailableProvider))
	flag.StringVar(&opt.API,
		"api", AvailableAPIs()[0], "wallpaper api to use. Available: "+fmt.Sprint(AvailableAPIs()))
	flag.StringVar(&opt.SaveWallpaperPath, "save-image-dir", os.Getenv("HOME")+"/Pictures", "directory to save image to")
	flag.StringVar(&opt.SearchPhrase, "search-phrase", "", "search phrase to use")
	flag.DurationVar(&opt.FollowDuration, "follow-time", time.Hour, "follow a time interval and update wallpaper. e.g. 1h, 1m, 30s")
	flag.BoolVar(&opt.Follow, "follow", false, "follow a time interval and update wallpaper. e.g. 1h, 1m, 30s")
	flag.BoolVar(&opt.Debug, "debug", false, "enable debug logs")

	flag.Parse()

	if err := opt.Validate(); err != nil {
		return opt, fmt.Errorf("invalid options: %v", err)
	}

	return opt, nil
}

func checkAvailable(concrete string, available []string) bool {
	for _, item := range available {
		if concrete == item {
			return true
		}
	}

	return false
}

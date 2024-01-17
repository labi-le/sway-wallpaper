package manager

import (
	"flag"
	"fmt"
	"github.com/labi-le/sway-wallpaper/pkg/log"
	"os"
)

func ParseOptions() Options {
	var (
		opt Options
	)

	flag.StringVar(
		&opt.Browser,
		"browserName",
		AvailableBrowsers()[0],
		"browserName to use. Available: "+fmt.Sprint(AvailableBrowsers()),
	)
	flag.StringVar(&opt.HistoryFile, "history-file", "",
		"browserName history file to use. Auto detect if empty (only for chromium based browsers)\n"+
			"e.g ~/.mozilla/icecat/gxda4hpz.default-1672760493248/formhistory.sqlite",
	)
	flag.Var(&opt.ImageResolution, "image-resolution", "image resolution. e.g. 1920x1080")
	flag.Var(&opt.Output, "output", "output to operate on. e.g eDP-1")
	flag.StringVar(&opt.Tool,
		"wp-tool", AvailableBGTools()[0], "wallpaper tool to use. Available: "+fmt.Sprint(AvailableBGTools()))
	flag.StringVar(&opt.API,
		"wp-api", AvailableAPIs()[0], "wallpaper api to use. Available: "+fmt.Sprint(AvailableAPIs()))
	flag.StringVar(&opt.SaveWallpaperPath, "save-image-dir", os.Getenv("HOME")+"/Pictures", "directory to save image to")
	flag.StringVar(&opt.SearchPhrase, "search-phrase", "", "search phrase to use")
	flag.DurationVar(&opt.FollowDuration, "follow", 0, "follow a time interval and update wallpaper. e.g. 1h, 1m, 30s")

	flag.Parse()

	if err := opt.Validate(); err != nil {
		log.Fatal(err)
	}

	return opt
}

func checkAvailable(concrete string, available []string) bool {
	for _, item := range available {
		if concrete == item {
			return true
		}
	}

	return false
}

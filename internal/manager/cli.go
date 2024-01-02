package manager

import (
	"errors"
	"flag"
	"fmt"
	"github.com/labi-le/sway-wallpaper/pkg/log"
	"github.com/labi-le/sway-wallpaper/pkg/output"
	"github.com/vcraescu/go-xrandr"
	"os"
	"os/exec"
	"regexp"
	"time"
)

var (
	ErrWallpaperToolNotImplemented = errors.New("invalid wallpaper tool not implemented")
	ErrWallpaperAPINotImplemented  = errors.New("invalid wallpaper api not implemented")
	ErrInvalidFollowDuration       = errors.New("invalid follow. e.g. 1h, 1m, 1s")
	ErrInvalidResolution           = errors.New("invalid resolution. e.g. 1920x1080")
	ErrAutoResolutionNotSupported  = errors.New("xrandr not found. Please install xrandr or set resolution manually")
	ErrDetectResolution            = errors.New("failed to detect resolution. Please set resolution manually")
)

func ParseOptions() Options {
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
		AvailableBrowsers()[0],
		"browserName to use. Available: "+fmt.Sprint(AvailableBrowsers()),
	)
	flag.StringVar(&historyFile, "history-file", "",
		"browserName history file to use. Auto detect if empty (only for chromium based browsers)\n"+
			"e.g ~/.mozilla/icecat/gxda4hpz.default-1672760493248/formhistory.sqlite",
	)
	flag.StringVar(&resolution, "resolution", DetectResolution(), "resolution to use. e.g. 1920x1080")
	flag.StringVar(&wpTool,
		"wp-tool", AvailableBGTools()[0], "wallpaper tool to use. Available: "+fmt.Sprint(AvailableBGTools()))
	flag.StringVar(&wpAPI,
		"wp-api", AvailableAPIs()[0], "wallpaper api to use. Available: "+fmt.Sprint(AvailableAPIs()))
	flag.StringVar(&saveImageDir, "save-image-dir", os.Getenv("HOME")+"/Pictures", "directory to save image to")
	flag.StringVar(&searchPhrase, "search-phrase", "", "search phrase to use")
	flag.StringVar(&follow, "follow", "", "follow a time interval and update wallpaper. e.g. 1h, 1m, 30s")

	flag.Parse()

	if !checkAvailable(browserName, AvailableBrowsers()) && searchPhrase == "" {
		log.Fatal(ErrBrowserNotImplemented)
	}

	if !checkAvailable(wpTool, AvailableBGTools()) {
		log.Fatal(ErrWallpaperToolNotImplemented)
	}

	if !checkAvailable(wpAPI, AvailableAPIs()) {
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

	return Options{
		SearchPhrase:      searchPhrase,
		SaveWallpaperPath: saveImageDir,
		FollowDuration:    followDuration,
		Resolution:        output.MustStringResolution(resolution),
		API:               wpAPI,
		Tool:              wpTool,
		HistoryFile:       historyFile,
		Browser:           browserName,
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
		log.Fatal(ErrDetectResolution)
	}

	// todo add select by monitor
	screen := screens[0].CurrentResolution
	resolution := fmt.Sprintf("%.fx%.f", screen.Width, screen.Height)
	log.Infof("Detected resolution: %s", resolution)

	return resolution
}

//func DetectResolution() string {
//	screens, err := xrandr.GetScreens()
//	if err != nil {
//		if errors.Is(err, exec.ErrNotFound) {
//			log.Fatal(ErrAutoResolutionNotSupported)
//		}
//		log.Fatal(err)
//	}
//
//	screen := screens[0]
//
//	resolution := fmt.Sprintf("%.fx%.f", screen.CurrentResolution.Width, screen.CurrentResolution.Height)
//	log.Infof("Detected resolution: %s", resolution)
//
//	return resolution
//}

func checkAvailable(concrete string, available []string) bool {
	for _, item := range available {
		if concrete == item {
			return true
		}
	}

	return false
}

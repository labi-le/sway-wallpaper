package wallpaper

import (
	"github.com/labi-le/history-wallpaper/pkg/api"
	"github.com/labi-le/history-wallpaper/pkg/browser"
	"github.com/labi-le/history-wallpaper/pkg/wptool"
	"time"
)

type HW struct {
	SaveWallpaperPath string
	SearchPhrase      string
	WallpaperAPI      api.Finder
	Browser           browser.PhraseFinder
	WallpaperTool     wptool.Setter
	FollowDuration    time.Duration
	Resolution        api.Resolution
}

package manager

import (
	"errors"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/labi-le/sway-wallpaper/pkg/api"
	"github.com/labi-le/sway-wallpaper/pkg/browser"
	"github.com/labi-le/sway-wallpaper/pkg/output"
	"github.com/labi-le/sway-wallpaper/pkg/wallpaper"
)

var (
	ErrWallpaperToolNotImplemented = errors.New("invalid wallpaper tool not implemented")
	ErrWallpaperAPINotImplemented  = errors.New("invalid wallpaper api not implemented")
	ErrBrowserNotImplemented       = errors.New("browser not implemented")
)

type Validator interface {
	Validate() error
}

type Options struct {
	ByPhrase          string
	SaveWallpaperPath string
	Follow            bool
	FollowDuration    time.Duration
	ImageResolution   output.Resolution
	API               string
	Tool              string
	HistoryFile       string
	Browser           string
	Output            output.Monitor
	Verbose           bool
}

func (o Options) Validate() error {
	var err error
	if !checkAvailable(o.Browser, browser.AvailableBrowsers()) && o.ByPhrase == "" {
		err = multierror.Append(err, ErrBrowserNotImplemented)
	}

	if _, ok := wallpaper.SupportedProvider[o.Tool]; !ok {
		err = multierror.Append(err, ErrWallpaperToolNotImplemented)
	}

	if !checkAvailable(o.API, api.AvailableAPIs()) {
		err = multierror.Append(err, ErrWallpaperAPINotImplemented)
	}

	return err
}

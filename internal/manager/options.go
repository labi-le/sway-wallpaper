package manager

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	"github.com/labi-le/sway-wallpaper/pkg/output"
	"time"
)

var (
	ErrWallpaperToolNotImplemented = errors.New("invalid wallpaper tool not implemented")
	ErrWallpaperAPINotImplemented  = errors.New("invalid wallpaper api not implemented")
)

type Validator interface {
	Validate() error
}

type Options struct {
	SearchPhrase      string
	SaveWallpaperPath string
	FollowDuration    time.Duration
	ImageResolution   output.Resolution
	API               string
	Tool              string
	HistoryFile       string
	Browser           string
	Output            output.Monitor
}

func (o Options) Validate() error {
	var err error
	if !checkAvailable(o.Browser, AvailableBrowsers()) && o.SearchPhrase == "" {
		err = multierror.Append(err, ErrBrowserNotImplemented)
	}

	if !checkAvailable(o.Tool, AvailableBGTools()) {
		err = multierror.Append(err, ErrWallpaperToolNotImplemented)
	}

	if !checkAvailable(o.API, AvailableAPIs()) {
		err = multierror.Append(err, ErrWallpaperAPINotImplemented)
	}

	return err
}

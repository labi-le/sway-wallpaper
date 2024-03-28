package manager

import (
	"context"
	"errors"
	"github.com/labi-le/sway-wallpaper/pkg/wptool"
)

func AvailableBGTools() []string {
	return []string{"swww", "swaybg"}
}

var ErrUnknownTool = errors.New("unknown tool")

type Setter interface {
	Set(ctx context.Context, path, output string) error
	Clean() error
}

func MustBGTool(tool string) Setter {
	switch tool {
	case "swaybg":
		return wptool.SwayBG{}
	case "swww":
		return wptool.S3W{}
	default:
		panic(ErrUnknownTool)
	}
}

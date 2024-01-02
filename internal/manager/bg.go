package manager

import (
	"context"
	"errors"
	"github.com/labi-le/sway-wallpaper/pkg/wptool"
)

func AvailableBGTools() []string {
	return []string{"swaybg"}
}

var ErrUnknownTool = errors.New("unknown tool")

type Setter interface {
	Set(ctx context.Context, path, output string) error
}

func MustBGTool(tool string) Setter {
	switch tool {
	case "swaybg":
		return wptool.SwayBG{}
	default:
		panic(ErrUnknownTool)
	}
}

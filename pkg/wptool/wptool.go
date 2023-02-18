package wptool

import (
	"errors"
	"golang.org/x/net/context"
)

func Available() []string {
	return []string{"swaybg", "wbg"}
}

var ErrUnknownTool = errors.New("unknown tool")

type Setter interface {
	Set(ctx context.Context, path string) error
}

func ParseTool(tool string) Setter {
	switch tool {
	case "swaybg":
		return SwayBG{}
	case "wbg":
		return WBG{}
	default:
		panic(ErrUnknownTool)
	}
}

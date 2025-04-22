package wallpaper

import "errors"

var ErrUnknownTool = errors.New("unknown tool")

func ByName(tool string) (Tool, error) {
	switch tool {
	case "swaybg":
		return NewSwayBG()
	case "swww":
		return NewSWWW()
	default:
		panic(ErrUnknownTool)
	}
}

var (
	SupportedProvider = map[string]struct{}{"swww": {}, "swaybg": {}}
	AvailableProvider = getAvailableProvider()
)

func getAvailableProvider() Tool {
	if t, err := NewSwayBG(); err == nil {
		return t
	}

	if t, err := NewSWWW(); err == nil {
		return t
	}

	return nil
}

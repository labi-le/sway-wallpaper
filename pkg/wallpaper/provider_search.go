package wallpaper

import (
	"github.com/labi-le/chiasma/pkg/wallpaper/execute"
	"github.com/labi-le/chiasma/pkg/wallpaper/swaybg"
	"github.com/labi-le/chiasma/pkg/wallpaper/swww"
)

func ByNameOrAvailable(tool string) (execute.Provider, error) {
	if tool == "" {
		if availableProvider == nil {
			return nil, execute.ErrUtilityNotFound
		}

		return availableProvider, nil
	}

	switch tool {
	case swaybg.Name:
		return swaybg.NewSwayBG()
	case swww.Name:
		return swww.NewSWWW()
	default:
		return nil, execute.ErrUtilityNotFound
	}
}

var (
	availableProvider = getAvailableProvider()
)

func getAvailableProvider() execute.Provider {
	if t, err := swaybg.NewSwayBG(); err == nil {
		return t
	}

	if t, err := swww.NewSWWW(); err == nil {
		return t
	}

	return nil
}

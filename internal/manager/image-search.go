package manager

import (
	"context"
	"errors"
	"github.com/labi-le/sway-wallpaper/pkg/api"
	"github.com/labi-le/sway-wallpaper/pkg/output"
	"io"
)

var (
	ErrUnknownService = errors.New("unknown service")
)

type Searcher interface {
	Search(ctx context.Context, q string, resolution output.Resolution) (io.ReadCloser, error)
}

func AvailableAPIs() []string {
	return []string{"unsplash"}
}

func MustSearcher(serviceName string) Searcher {
	switch serviceName {
	case "unsplash":
		return &api.Unsplash{}
	default:
		panic(ErrUnknownService)
	}
}

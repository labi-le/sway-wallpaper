package api

import (
	"context"
	"errors"
	"io"
)

type Resolution string // 1920x1080

var ErrUnknownService = errors.New("unknown service")

func Available() []string {
	return []string{"unsplash"}
}

type Finder interface {
	Find(ctx context.Context, q string, r Resolution) (io.ReadCloser, error)
}

func MustFinder(serviceName string) Finder {
	switch serviceName {
	case "unsplash":
		return &Unsplash{}
	default:
		panic(ErrUnknownService)
	}
}

package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Resolution string // 1920x1080

var (
	ErrUnknownService    = errors.New("unknown service")
	ErrConnectionTimeOut = errors.New("api %s connection timeout")
)

func Available() []string {
	return []string{"unsplash"}
}

type Finder interface {
	Find(ctx context.Context, q string, r Resolution) (io.ReadCloser, error)
}

func MustFinder(serviceName string) Finder {
	switch serviceName {
	case "unsplash":
		return &Unsplash{&api{}}
	default:
		panic(ErrUnknownService)
	}
}

// api is a wrapper around http.Client to handle context.Canceled error as ErrConnectionTimeOut
type api struct {
	http.Client
}

func (a *api) Do(req *http.Request) (*http.Response, error) {
	do, err := a.Client.Do(req)
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, fmt.Errorf(ErrConnectionTimeOut.Error(), req.URL.Host)
	}

	return do, err
}

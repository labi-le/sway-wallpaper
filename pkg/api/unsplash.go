package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/labi-le/sway-wallpaper/pkg/output"
	"io"
	"net/http"
)

var (
	ErrConnectionTimeOut = errors.New("api %s connection timeout")
)

type Unsplash struct {
	client api
}

func (u *Unsplash) Search(ctx context.Context, q string, resolution output.Resolution) (io.ReadCloser, error) {
	request, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://source.unsplash.com/%s/?%s", resolution.String(), q),
		nil,
	)

	resp, err := u.client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// api is a wrapper around http.Client to handle context.Canceled error as ErrConnectionTimeOut
type api struct {
	http.Client
}

func (a *api) Do(req *http.Request) (*http.Response, error) {
	do, err := a.Client.Do(req)
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, fmt.Errorf(ErrConnectionTimeOut.Error(), req.URL.String())
	}

	return do, err
}

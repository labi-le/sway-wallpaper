package api

import (
	"bytes"
	"context"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/labi-le/chiasma/pkg/output"
	"github.com/rs/zerolog"
)

var (
	ErrUnknownSearcher = errors.New("unknown searcher")
)

type Image interface {
	io.ReadCloser
	Size() (int, int)
}

type Searcher interface {
	Search(ctx context.Context, q string, resolution output.Resolution) (Image, error)
}

type detectedImage struct {
	io.Reader
	closer io.Closer
	w, h   int
}

func (d *detectedImage) Size() (int, int) { return d.w, d.h }
func (d *detectedImage) Close() error     { return d.closer.Close() }

func DetectSize(img io.Reader) (Image, error) {
	var header bytes.Buffer
	tee := io.TeeReader(img, &header)

	config, _, err := image.DecodeConfig(tee)
	if err != nil {
		return nil, err
	}

	var closer io.Closer
	if c, ok := img.(io.Closer); ok {
		closer = c
	} else {
		closer = io.NopCloser(nil)
	}

	return &detectedImage{
		Reader: io.MultiReader(&header, img),
		closer: closer,
		w:      config.Width,
		h:      config.Height,
	}, nil
}

func AvailableAPIs() []string {
	return []string{"nasa", "unsplash"}
}

func NewSearcher(log zerolog.Logger, name string) (Searcher, error) {
	switch name {
	case "unsplash":
		return NewUnsplash(log), nil
	case "nasa":
		return NewNasa(log), nil
	default:
		return nil, ErrUnknownSearcher
	}
}

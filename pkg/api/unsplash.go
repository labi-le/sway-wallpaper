package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/labi-le/sway-wallpaper/pkg/output"
	"github.com/rs/zerolog"
)

var (
	ErrConnectionTimeOut = errors.New("connection timeout")
)

const (
	searchQuery = "https://unsplash.com/napi/search/photos?page=1&per_page=20&query=%s&xp=reset-search-state%%3Aexperiment"
)

type Unsplash struct {
	log    zerolog.Logger
	client api
}

func NewUnsplash(log zerolog.Logger) *Unsplash {
	return &Unsplash{
		log: log.With().Str("component", "unsplash").Logger(),
	}
}

type SearchResult struct {
	Results []Photo `json:"results"`
}

type Photo struct {
	Id     string `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Urls   struct {
		Full string `json:"full"`
	} `json:"urls"`
	Premium bool `json:"premium"`
}

type unsplashImage struct {
	io.ReadCloser
	w, h int
}

func (i unsplashImage) Size() (int, int) {
	return i.w, i.h
}

func (u *Unsplash) Search(ctx context.Context, q string, resolution output.Resolution) (Image, error) {
	log := u.log.With().Str("op", "Search").Logger()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			searchQuery,
			url.QueryEscape(q),
		),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://unsplash.com/s/photos/"+url.QueryEscape(q))
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36 CrKey/1.54.250320")

	log.Trace().Msgf("requesting unsplash search for: %s", q)
	photo, err := u.tryFetch(req)
	if err != nil {
		return nil, err
	}

	imgURL := fmt.Sprintf("%s&w=%d&h=%d", photo.Urls.Full, resolution.Width, resolution.Height)
	log.Trace().Msgf("requesting image from unsplash: %s", imgURL)
	get, err := u.client.Get(imgURL)
	if err != nil {
		return nil, err
	}
	return unsplashImage{ReadCloser: get.Body, w: photo.Width, h: photo.Height}, nil
}

func (u *Unsplash) tryFetch(req *http.Request) (Photo, error) {
	log := u.log.With().Str("op", "tryFetch").Logger()
	for i := 0; i < 5; i++ {
		resp, err := u.client.Do(req)
		if err != nil {
			return Photo{}, fmt.Errorf("server returned an error: %w", err)
		}

		var r SearchResult
		if decodeErr := json.NewDecoder(resp.Body).Decode(&r); decodeErr != nil {
			resp.Body.Close()
			return Photo{}, fmt.Errorf("error decoding response: %w", decodeErr)
		}
		resp.Body.Close()

		var candidates []Photo
		for _, photo := range r.Results {
			if !photo.Premium {
				candidates = append(candidates, photo)
			}
		}

		if len(candidates) > 0 {
			return candidates[rand.Intn(len(candidates))], nil
		}

		log.Trace().Msg("got a watermarked photo, trying again")
	}

	return Photo{}, errors.New("failed to fetch watermarked photo after multiple attempts")
}

type api struct {
	http.Client
}

func (a *api) Do(req *http.Request) (*http.Response, error) {
	do, err := a.Client.Do(req)
	if err != nil {
		if do != nil {
			do.Body.Close()
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%w: api: %s", ErrConnectionTimeOut, req.URL.String())
		}
		return nil, err
	}
	return do, nil
}

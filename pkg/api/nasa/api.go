package nasa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/labi-le/chiasma/pkg/api/searcher"
	"github.com/rs/zerolog"
)

const Name = "nasa"

const (
	nasaSearchURL = "https://images-api.nasa.gov/search?q=%s&media_type=image"
)

type Nasa struct {
	log    zerolog.Logger
	client http.Client
}

func NewNasa(log zerolog.Logger) *Nasa {
	return &Nasa{
		log: log.With().Str("component", "nasa").Logger(),
	}
}

type nasaSearchResult struct {
	Collection struct {
		Items []struct {
			Href string `json:"href"`
		} `json:"items"`
	} `json:"collection"`
}

type nasaImage struct {
	io.ReadCloser
	w, h int
}

func (i nasaImage) Size() (int, int) {
	return i.w, i.h
}

func (n *Nasa) Search(ctx context.Context, q string, _ searcher.Resolution) (searcher.Image, error) {
	log := n.log.With().Str("op", "Search").Logger()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(nasaSearchURL, url.QueryEscape(q)), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	log.Trace().Msgf("requesting nasa search for: %s", q)
	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nasa api returned status: %d", resp.StatusCode)
	}

	var searchRes nasaSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&searchRes); err != nil {
		return nil, fmt.Errorf("failed to decode search results: %w", err)
	}

	items := searchRes.Collection.Items
	if len(items) == 0 {
		return nil, fmt.Errorf("no images found for query: %s", q)
	}

	item := items[rand.Intn(len(items))]
	log.Trace().Msgf("selected nasa item metadata href: %s", item.Href)

	assetReq, err := http.NewRequestWithContext(ctx, http.MethodGet, item.Href, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset request: %w", err)
	}

	assetResp, err := n.client.Do(assetReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute asset fetch: %w", err)
	}
	defer assetResp.Body.Close()

	var assetURLs []string
	if err := json.NewDecoder(assetResp.Body).Decode(&assetURLs); err != nil {
		return nil, fmt.Errorf("failed to decode asset results: %w", err)
	}

	var imgURL string
	for _, a := range assetURLs {
		if strings.Contains(a, "~orig.jpg") {
			imgURL = a
			break
		}
	}

	if imgURL == "" {
		for _, a := range assetURLs {
			if strings.Contains(a, "~large.jpg") {
				imgURL = a
				break
			}
		}
	}

	if imgURL == "" && len(assetURLs) > 0 {
		imgURL = assetURLs[0]
	}

	if imgURL == "" {
		return nil, fmt.Errorf("no image assets found for nasa item")
	}

	log.Trace().Msgf("selected nasa image url: %s", imgURL)

	imgReq, err := http.NewRequestWithContext(ctx, http.MethodGet, imgURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create image request: %w", err)
	}

	imgResp, err := n.client.Do(imgReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}

	if imgResp.StatusCode != http.StatusOK {
		_ = imgResp.Body.Close()
		return nil, fmt.Errorf("failed to download image, status: %d", imgResp.StatusCode)
	}

	return searcher.DetectSize(imgResp.Body)
}

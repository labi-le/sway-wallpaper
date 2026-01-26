package nasa

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
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

var (
	stopWords = []string{
		// Technical terms
		"chart", "diagram", "plot", "spectrum", "graph",
		"schematic", "profile", "response", "model", "map",
		"histogram", "curve", "data", "survey", "sensor",
		// Visual formats irrelevant for wallpapers
		"mosaic", "panorama", "composite",
		// Sources that often produce "scientific strips" (like PIA05979)
		"galex", "evolution explorer",
	}
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

type nasaData struct {
	Title       string   `json:"title"`
	Keywords    []string `json:"keywords"`
	Description string   `json:"description"`
}

type nasaItem struct {
	Href string     `json:"href"`
	Data []nasaData `json:"data"`
}

type nasaSearchResult struct {
	Collection struct {
		Items []nasaItem `json:"items"`
	} `json:"collection"`
}

func (n *Nasa) Search(ctx context.Context, q string, res searcher.Resolution) (searcher.Image, error) {
	log := n.log.With().Str("op", "Search").Logger()

	items, err := n.fetchSearchResults(ctx, q)
	if err != nil {
		return nil, err
	}

	candidates := n.filterCandidates(items)

	if len(candidates) == 0 {
		log.Warn().Msg("all images were filtered out as technical, falling back to raw results")
		candidates = items
	} else {
		log.Info().Int("total", len(items)).Int("clean", len(candidates)).Msg("filtering complete")
	}

	const maxRetries = 10
	shuffled := make([]nasaItem, len(candidates))
	copy(shuffled, candidates)
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	for i := 0; i < len(shuffled) && i < maxRetries; i++ {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		selected := shuffled[i]
		log.Debug().Str("href", selected.Href).Int("attempt", i+1).Msg("trying candidate")

		imgURL, err := n.resolveImageURL(ctx, selected.Href)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			log.Warn().Err(err).Msg("failed to resolve image url")
			continue
		}

		img, err := n.downloadImage(ctx, imgURL)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			log.Warn().Err(err).Msg("failed to download image")
			continue
		}

		w, h := img.Size()
		if isAspectRatioBad(w, h, res.Width, res.Height) {
			img.Close()
			log.Debug().
				Int("w", w).Int("h", h).
				Msg("image rejected: bad aspect ratio (panorama/strip detected)")
			continue
		}

		return img, nil
	}

	return nil, fmt.Errorf("failed to find suitable image after %d attempts", maxRetries)
}

func isAspectRatioBad(imgW, imgH, targetW, targetH int) bool {
	if targetW == 0 || targetH == 0 {
		ratio := float64(imgW) / float64(imgH)
		return ratio > 2.5 || ratio < 0.4
	}

	targetRatio := float64(targetW) / float64(targetH)
	imgRatio := float64(imgW) / float64(imgH)

	diff := math.Abs(targetRatio - imgRatio)
	return diff > 1.0
}

func (n *Nasa) filterCandidates(items []nasaItem) []nasaItem {
	candidates := make([]nasaItem, 0, len(items))

	for _, item := range items {
		clean, reason := n.inspect(item)

		if len(item.Data) > 0 {
			meta := item.Data[0]
			if !clean {
				n.log.Debug().
					Str("title", meta.Title).
					Str("reject_reason", reason).
					Msg("candidate rejected")
			}
		}

		if clean {
			candidates = append(candidates, item)
		}
	}
	return candidates
}

func (n *Nasa) inspect(item nasaItem) (bool, string) {
	if len(item.Data) == 0 {
		return true, ""
	}

	meta := item.Data[0]
	title := strings.ToLower(meta.Title)
	desc := strings.ToLower(meta.Description)

	if found, word := containsStopWord(title); found {
		return false, "title_" + word
	}

	if found, word := containsStopWord(desc); found {
		return false, "desc_" + word
	}

	for _, k := range meta.Keywords {
		kLower := strings.ToLower(k)
		if found, word := containsStopWord(kLower); found {
			return false, "kw_" + word
		}
	}

	return true, ""
}

func containsStopWord(text string) (bool, string) {
	for _, stop := range stopWords {
		if !strings.Contains(text, stop) {
			continue
		}

		if stop == "graph" && (strings.Contains(text, "photograph") || strings.Contains(text, "graphic")) {
			continue
		}

		return true, stop
	}
	return false, ""
}

func (n *Nasa) fetchSearchResults(ctx context.Context, q string) ([]nasaItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(nasaSearchURL, url.QueryEscape(q)), nil)
	if err != nil {
		return nil, fmt.Errorf("create req: %w", err)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api status: %d", resp.StatusCode)
	}

	var res nasaSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}

	if len(res.Collection.Items) == 0 {
		return nil, fmt.Errorf("no results for %s", q)
	}

	return res.Collection.Items, nil
}

func (n *Nasa) resolveImageURL(ctx context.Context, href string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, href, nil)
	if err != nil {
		return "", fmt.Errorf("create asset req: %w", err)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do asset req: %w", err)
	}
	defer resp.Body.Close()

	var assets []string
	if err := json.NewDecoder(resp.Body).Decode(&assets); err != nil {
		return "", fmt.Errorf("decode assets: %w", err)
	}

	return findBestImage(assets), nil
}

func (n *Nasa) downloadImage(ctx context.Context, url string) (searcher.Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create img req: %w", err)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download img: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("img status: %d", resp.StatusCode)
	}

	return searcher.DetectSize(resp.Body)
}

func findBestImage(urls []string) string {
	for _, u := range urls {
		if strings.Contains(u, "~orig.jpg") {
			return u
		}
	}
	for _, u := range urls {
		if strings.Contains(u, "~large.jpg") {
			return u
		}
	}
	if len(urls) > 0 {
		return urls[0]
	}
	return ""
}

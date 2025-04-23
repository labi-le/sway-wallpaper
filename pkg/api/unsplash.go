package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labi-le/sway-wallpaper/pkg/output"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

var (
	ErrConnectionTimeOut = errors.New("connection timeout")
)

const (
	ua    = "Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Mobile Safari/537.36"
	query = "https://unsplash.com/napi/photos/random?query=%s&xp=new-plus-algorithm:experiment&plus=none&orientation=landscape"
)

type Unsplash struct {
	client api
}

type randomPhoto struct {
	Id               string `json:"id"`
	Slug             string `json:"slug"`
	AlternativeSlugs struct {
		En string `json:"en"`
		Es string `json:"es"`
		Ja string `json:"ja"`
		Fr string `json:"fr"`
		It string `json:"it"`
		Ko string `json:"ko"`
		De string `json:"de"`
		Pt string `json:"pt"`
	} `json:"alternative_slugs"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	PromotedAt     interface{}   `json:"promoted_at"`
	Width          int           `json:"width"`
	Height         int           `json:"height"`
	Color          string        `json:"color"`
	BlurHash       string        `json:"blur_hash"`
	Description    string        `json:"description"`
	AltDescription string        `json:"alt_description"`
	Breadcrumbs    []interface{} `json:"breadcrumbs"`
	Urls           struct {
		Raw     string `json:"raw"`
		Full    string `json:"full"`
		Regular string `json:"regular"`
		Small   string `json:"small"`
		Thumb   string `json:"thumb"`
		SmallS3 string `json:"small_s3"`
	} `json:"urls"`
	Links struct {
		Self             string `json:"self"`
		Html             string `json:"html"`
		Download         string `json:"download"`
		DownloadLocation string `json:"download_location"`
	} `json:"links"`
	Likes                  int           `json:"likes"`
	LikedByUser            bool          `json:"liked_by_user"`
	CurrentUserCollections []interface{} `json:"current_user_collections"`
	Sponsorship            interface{}   `json:"sponsorship"`
	TopicSubmissions       struct {
	} `json:"topic_submissions"`
	AssetType string `json:"asset_type"`
	Premium   bool   `json:"premium"`
	Plus      bool   `json:"plus"`
	User      struct {
		Id              string      `json:"id"`
		UpdatedAt       time.Time   `json:"updated_at"`
		Username        string      `json:"username"`
		Name            string      `json:"name"`
		FirstName       string      `json:"first_name"`
		LastName        string      `json:"last_name"`
		TwitterUsername interface{} `json:"twitter_username"`
		PortfolioUrl    interface{} `json:"portfolio_url"`
		Bio             interface{} `json:"bio"`
		Location        interface{} `json:"location"`
		Links           struct {
			Self      string `json:"self"`
			Html      string `json:"html"`
			Photos    string `json:"photos"`
			Likes     string `json:"likes"`
			Portfolio string `json:"portfolio"`
		} `json:"links"`
		ProfileImage struct {
			Small  string `json:"small"`
			Medium string `json:"medium"`
			Large  string `json:"large"`
		} `json:"profile_image"`
		InstagramUsername          string `json:"instagram_username"`
		TotalCollections           int    `json:"total_collections"`
		TotalLikes                 int    `json:"total_likes"`
		TotalPhotos                int    `json:"total_photos"`
		TotalPromotedPhotos        int    `json:"total_promoted_photos"`
		TotalIllustrations         int    `json:"total_illustrations"`
		TotalPromotedIllustrations int    `json:"total_promoted_illustrations"`
		AcceptedTos                bool   `json:"accepted_tos"`
		ForHire                    bool   `json:"for_hire"`
		Social                     struct {
			InstagramUsername string      `json:"instagram_username"`
			PortfolioUrl      interface{} `json:"portfolio_url"`
			TwitterUsername   interface{} `json:"twitter_username"`
			PaypalEmail       interface{} `json:"paypal_email"`
		} `json:"social"`
	} `json:"user"`
	Exif struct {
		Make         string `json:"make"`
		Model        string `json:"model"`
		Name         string `json:"name"`
		ExposureTime string `json:"exposure_time"`
		Aperture     string `json:"aperture"`
		FocalLength  string `json:"focal_length"`
		Iso          int    `json:"iso"`
	} `json:"exif"`
	Location struct {
		Name     string `json:"name"`
		City     string `json:"city"`
		Country  string `json:"country"`
		Position struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"position"`
	} `json:"location"`
	Meta struct {
		Index bool `json:"index"`
	} `json:"meta"`
	PublicDomain bool `json:"public_domain"`
	Tags         []struct {
		Type  string `json:"type"`
		Title string `json:"title"`
	} `json:"tags"`
	Views     int           `json:"views"`
	Downloads int           `json:"downloads"`
	Topics    []interface{} `json:"topics"`
}

func (u *Unsplash) Search(ctx context.Context, q string, resolution output.Resolution) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			query,
			q,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "https://unsplash.com/s/photos/"+q)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?1")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Android\"")
	req.Header.Set("User-Agent", ua)

	url, err := u.tryFetch(req)
	if err != nil {
		return nil, err
	}

	get, err := u.client.Get(fmt.Sprintf("%s&w=%d&h=%d", url, resolution.Width, resolution.Height))
	if err != nil {
		return nil, err
	}
	return get.Body, nil
}

func (u *Unsplash) tryFetch(req *http.Request) (string, error) {
	for i := 0; i < 5; i++ {
		resp, err := u.client.Do(req)
		if err != nil {
			return "", fmt.Errorf("server returned an error: %v", err)
		}

		var r randomPhoto
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return "", fmt.Errorf("error decoding response: %v", err)
		}

		if !r.Premium {
			return r.Urls.Full, nil
		}
		log.Trace().Msg("got a watermarked photo, trying again")
	}

	return "", errors.New("failed to fetch watermarked photo after multiple attempts")
}

// api is a wrapper around http.Client to handle context.Canceled error as ErrConnectionTimeOut
type api struct {
	http.Client
}

func (a *api) Do(req *http.Request) (*http.Response, error) {
	do, err := a.Client.Do(req)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return nil, fmt.Errorf("%v: api: %s", ErrConnectionTimeOut, req.URL.String())
	default:
		return do, err
	}
}

package unsplash

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetImage(phrase string, resolution string) ([]byte, error) {
	client := http.Client{}
	get, err := client.Get(
		fmt.Sprintf("https://source.unsplash.com/%s/?%s",
			resolution,
			url.QueryEscape(phrase),
		),
	)
	if err != nil {
		return nil, err
	}

	defer get.Body.Close()

	return io.ReadAll(get.Body)
}

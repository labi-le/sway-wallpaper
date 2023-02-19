package browser

import (
	"database/sql"
	"errors"
	"github.com/labi-le/history-wallpaper/pkg/log"
	"net/url"
)

type Chromium struct {
	Name    string
	History *History
}

func (b *Chromium) LastSearchedPhrase() (string, error) {
	defer b.History.Cleanup()

	// It is also possible to search in the keyword_search_terms table,
	// but the last query will not contain the search phrase,
	// because repeated search queries are not displayed as the latest in search history
	rows := b.History.DB.QueryRow(`
		SELECT url
FROM urls
WHERE url LIKE 'https://www.google.com/search?%'
ORDER BY last_visit_time DESC
LIMIT 1;

`)
	var lastURL string
	if err := rows.Scan(&lastURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrHistoryIsEmpty
		}
		return "", err
	}

	ur, parseErr := url.Parse(lastURL)
	if parseErr != nil {
		panic(parseErr)
	}

	phrase := ur.Query().Get("q")
	log.Infof("Last searched phrase: %s", phrase)
	return phrase, nil
}

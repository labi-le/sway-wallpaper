package browser

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"

	_ "modernc.org/sqlite"
)

var ErrHistoryIsEmpty = errors.New("browser history is empty")

type History interface {
	GetLastSearch() (string, error)
	Close() error
}

type chromiumHistory struct{ db *sql.DB }

func (h *chromiumHistory) Close() error { return h.db.Close() }
func (h *chromiumHistory) GetLastSearch() (string, error) {
	var lastURL string
	err := h.db.QueryRow(`
		SELECT url FROM urls
		WHERE url LIKE 'https://www.google.com/search?%'
		ORDER BY last_visit_time DESC LIMIT 1
	`).Scan(&lastURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrHistoryIsEmpty
		}
		return "", err
	}
	u, err := url.Parse(lastURL)
	if err != nil {
		return "", err
	}
	return u.Query().Get("q"), nil
}

type firefoxHistory struct{ db *sql.DB }

func (h *firefoxHistory) Close() error { return h.db.Close() }
func (h *firefoxHistory) GetLastSearch() (string, error) {
	var value string
	err := h.db.QueryRow(`
		SELECT value FROM moz_formhistory
		WHERE fieldname = 'searchbar-history'
		ORDER BY lastUsed DESC LIMIT 1
	`).Scan(&value)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrHistoryIsEmpty
		}
		return "", err
	}
	return value, nil
}

type noopHistory struct{}

func (h *noopHistory) Close() error                   { return nil }
func (h *noopHistory) GetLastSearch() (string, error) { return "", nil }

func openHistoryDB(browserName string, fullPath string) (History, error) {
	path := fullPath
	isChromium := IsChromiumBased(browserName)

	if path == "" && isChromium {
		path = fmt.Sprintf("%s/.config/%s/Default/History", os.Getenv("HOME"), browserName)
	}
	if path == "" && !isChromium {
		return nil, errors.New("firefox-based browsers do not support auto-detecting history file")
	}

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?immutable=1&mode=ro", path))
	if err != nil {
		return nil, err
	}

	if isChromium {
		return &chromiumHistory{db: db}, nil
	}
	return &firefoxHistory{db: db}, nil
}

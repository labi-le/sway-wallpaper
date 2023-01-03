package chromium

import (
	"database/sql"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/user"
)

func GetLastSearchedPhrase(u *user.User, browser string, file string) (string, error) {
	if file == "" {
		file = fmt.Sprintf("%s/.config/%s/Default/History", u.HomeDir, browser)
	}

	open, osErr := os.Open(file)
	if osErr != nil {
		return "", osErr
	}

	defer open.Close()

	temp, tempErr := os.CreateTemp("", "history")
	if tempErr != nil {
		return "", tempErr
	}

	defer temp.Close()
	defer os.Remove(temp.Name())

	if _, ioErr := io.Copy(temp, open); ioErr != nil {
		return "", ioErr
	}

	db, sqlErr := sql.Open("sqlite", temp.Name())
	if sqlErr != nil {
		return "", sqlErr
	}

	defer db.Close()

	// It is also possible to search in the keyword_search_terms table,
	// but the last query will not contain the search phrase,
	// because repeated search queries are not displayed as the latest in search history
	rows := db.QueryRow(`
		SELECT url
FROM urls
WHERE url LIKE 'https://www.google.com/search?%'
ORDER BY last_visit_time DESC
LIMIT 1;

`)
	var lastURL string
	if err := rows.Scan(&lastURL); err != nil {
		return "", err
	}

	ur, parseErr := url.Parse(lastURL)
	if parseErr != nil {
		panic(parseErr)
	}

	phrase := ur.Query().Get("q")

	return phrase, nil
}

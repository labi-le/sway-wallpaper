package chromium

import (
	"database/sql"
	"fmt"
	"io"
	_ "modernc.org/sqlite"
	"net/url"
	"os"
	"os/user"
)

func GetLastSearchedPhrase(u *user.User, browser string) (string, error) {
	browserDir := fmt.Sprintf("%s/.config/%s/Default/History", u.HomeDir, browser)

	open, osErr := os.Open(browserDir)
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

	rows := db.QueryRow(`
		SELECT url
FROM urls
WHERE url LIKE 'https://www.google.com/search?%'
  AND last_visit_time > strftime('%s', 'now', '-1 hours') * 1000000 + 11644473600
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

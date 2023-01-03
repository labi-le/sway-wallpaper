package firefox

import (
	"database/sql"
	"fmt"
	"io"
	"os"
)

func GetLastSearchedPhrase(file string) (string, error) {
	if file == "" {
		return "", fmt.Errorf(
			"firefox-based browsers not support auto-detecting history file. Set formhistory.sqlite path manually")
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

	rows := db.QueryRow(`
		SELECT value
FROM moz_formhistory
WHERE fieldname = 'searchbar-history'
ORDER BY lastUsed DESC
LIMIT 1;

`)
	var value string
	return value, rows.Scan(&value)
}

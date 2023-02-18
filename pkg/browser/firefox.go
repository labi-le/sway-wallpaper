package browser

import (
	"database/sql"
	"errors"
	"os"
)

type Firefox struct {
	Name        string
	HistoryFile *os.File
}

func (f *Firefox) LastSearchedPhrase() (string, error) {
	defer func() {
		f.HistoryFile.Close()
		os.Remove(f.HistoryFile.Name())
	}()

	db, sqlErr := sql.Open("sqlite", f.HistoryFile.Name())
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
	err := rows.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrHistoryIsEmpty
		}
		return "", err
	}

	return value, nil
}

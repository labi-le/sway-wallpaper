package browser

import (
	"database/sql"
	"errors"
)

type Firefox struct {
	Name    string
	History *History
}

func (f *Firefox) Analyze() (string, error) {
	defer f.History.Cleanup()

	rows := f.History.DB.QueryRow(`
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

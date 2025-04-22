package browser

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

var (
	ErrHistoryIsEmpty = errors.New("browser history is empty")
)

var (
	ChromiumBasedBrowsers = []string{"vivaldi", "chrome", "chromium", "brave", "opera"}
	FirefoxBasedBrowsers  = []string{"firefox"}
)

func IsChromiumBased(browser string) bool {
	for _, b := range ChromiumBasedBrowsers {
		if browser == b {
			return true
		}
	}

	return false
}

func OpenHistoryDB(browserName string, fullPath string) *History {
	var path string

	if IsChromiumBased(browserName) {
		if fullPath != "" {
			path = fullPath
		} else {
			path = fmt.Sprintf("%s/.config/%s/Default/History", os.Getenv("HOME"), browserName)
		}

		return copyHistoryFile(path)
	}

	if fullPath == "" {
		panic("firefox-based browsers not support auto-detecting history file. Provide formhistory.sqlite path manually")
	}

	return copyHistoryFile(path)
}

func copyHistoryFile(path string) *History {
	lockedDB, osErr := os.Open(path)
	if osErr != nil {
		panic(osErr)
	}

	defer lockedDB.Close()

	temp, tempErr := os.CreateTemp("", "history")
	if tempErr != nil {
		panic(tempErr)
	}

	written, ioErr := io.Copy(temp, lockedDB)
	if ioErr != nil {
		panic(ioErr)
	}
	log.Info().Msgf("copied %d bytes from %s to %s", written, path, temp.Name())

	db, sqlErr := sql.Open("sqlite", temp.Name()+"?")
	if sqlErr != nil {
		panic(sqlErr)
	}

	return &History{
		DB:      db,
		TmpFile: temp,
	}
}

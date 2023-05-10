package browser

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/labi-le/sway-wallpaper/pkg/log"
	"io"
	"os"
	"os/user"
)

var (
	ChromiumBasedBrowsers = []string{"vivaldi", "chrome", "chromium", "brave", "opera"}
	FirefoxBasedBrowsers  = []string{"firefox"}
)

var (
	ErrHistoryIsEmpty = errors.New("browser history is empty")
)

func Available() []string {
	return append(ChromiumBasedBrowsers, FirefoxBasedBrowsers...)
}

type PhraseFinder interface {
	LastSearchedPhrase() (string, error)
}

func IsChromiumBased(browser string) bool {
	for _, b := range ChromiumBasedBrowsers {
		if browser == b {
			return true
		}
	}

	return false
}

func MustBrowser(browserName string, usr *user.User, historyPath string) PhraseFinder {
	if IsChromiumBased(browserName) {
		return &Chromium{
			Name:    browserName,
			History: openHistoryDB(browserName, usr, historyPath),
		}
	}

	if browserName == Noop {
		return NewNoop()
	}

	return &Firefox{
		Name:    browserName,
		History: openHistoryDB(browserName, usr, historyPath),
	}
}

func openHistoryDB(browserName string, u *user.User, fullPath string) *History {
	var path string

	if IsChromiumBased(browserName) {
		if fullPath != "" {
			path = fullPath
		} else {
			path = fmt.Sprintf("%s/.config/%s/Default/History", u.HomeDir, browserName)
		}

		return copyHistoryFile(path)
	}

	if fullPath == "" {
		panic("firefox-based browsers not support auto-detecting history file. Set formhistory.sqlite path manually")
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
	log.Infof("Copied %d bytes from %s to %s", written, path, temp.Name())

	db, sqlErr := sql.Open("sqlite", temp.Name()+"?")
	if sqlErr != nil {
		panic(sqlErr)
	}

	return &History{
		DB:      db,
		TmpFile: temp,
	}
}

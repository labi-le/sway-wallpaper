package browser

import (
	"errors"
	"fmt"
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
			Name:        browserName,
			HistoryFile: findHistoryFile(browserName, usr, historyPath),
		}
	}

	return &Firefox{
		Name:        browserName,
		HistoryFile: findHistoryFile(browserName, usr, historyPath),
	}
}

func findHistoryFile(browserName string, u *user.User, fullPath string) *os.File {
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

func copyHistoryFile(path string) *os.File {
	open, osErr := os.Open(path)
	if osErr != nil {
		panic(osErr)
	}

	defer open.Close()

	temp, tempErr := os.CreateTemp("", "history")
	if tempErr != nil {
		panic(tempErr)
	}

	if _, ioErr := io.Copy(temp, open); ioErr != nil {
		panic(ioErr)
	}

	return temp
}

package manager

import (
	"errors"
	"github.com/labi-le/sway-wallpaper/pkg/browser"
)

var (
	ErrBrowserNotImplemented = errors.New("browser not implemented")
)

func AvailableBrowsers() []string {
	return append(browser.ChromiumBasedBrowsers, browser.FirefoxBasedBrowsers...)
}

type Analyzer interface {
	Analyze() (string, error)
}

func NewBrowserHistoryAnalyzer(browserName string, historyPath string) Analyzer {
	if browser.IsChromiumBased(browserName) {
		return &browser.Chromium{
			Name:    browserName,
			History: browser.OpenHistoryDB(browserName, historyPath),
		}
	}

	if browserName == browser.NoopBrowser {
		return browser.NewNoop()
	}

	return &browser.Firefox{
		Name:    browserName,
		History: browser.OpenHistoryDB(browserName, historyPath),
	}
}

func NewPhraseAnalyzer(phrase string) Analyzer {
	return &browser.Noop{}
}

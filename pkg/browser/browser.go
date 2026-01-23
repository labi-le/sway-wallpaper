package browser

var (
	ChromiumBasedBrowsers = []string{"google-chrome", "vivaldi", "chromium", "brave", "opera"}
	FirefoxBasedBrowsers  = []string{"firefox"}
)

func AvailableBrowsers() []string {
	return append(ChromiumBasedBrowsers, FirefoxBasedBrowsers...)
}

func IsChromiumBased(browser string) bool {
	for _, b := range ChromiumBasedBrowsers {
		if browser == b {
			return true
		}
	}
	return false
}

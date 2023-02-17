package tools

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

func GetAllBrowsers() []string {
	return append(ChromiumBasedBrowsers, FirefoxBasedBrowsers...)
}

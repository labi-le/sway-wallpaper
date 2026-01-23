package browser

func NewHistoryProvider(name string, customPath string) (History, error) {
	if name == "noop" {
		return &noopHistory{}, nil
	}

	return openHistoryDB(name, customPath)
}

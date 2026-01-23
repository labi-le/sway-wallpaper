package browser

type Firefox struct {
	Name    string
	History History
}

func (f *Firefox) Analyze() (string, error) {
	defer f.History.Close()
	return f.History.GetLastSearch()
}
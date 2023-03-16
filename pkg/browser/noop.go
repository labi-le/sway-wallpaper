package browser

const Noop = "noop"

type noop struct{}

func NewNoop() PhraseFinder {
	return &noop{}
}

func (n *noop) LastSearchedPhrase() (string, error) {
	return "", nil
}

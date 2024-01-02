package browser

const NoopBrowser = "noop"

type Noop struct{}

func NewNoop() *Noop {
	return &Noop{}
}

func (n *Noop) Analyze() (string, error) {
	return "", nil
}

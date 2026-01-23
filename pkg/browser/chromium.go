package browser

import (
	"github.com/rs/zerolog"
)

type Chromium struct {
	Name    string
	History History
	log     zerolog.Logger
}

func NewChromium(log zerolog.Logger, name string, history History) *Chromium {
	return &Chromium{
		Name:    name,
		History: history,
		log:     log.With().Str("component", "chromium").Logger(),
	}
}

func (b *Chromium) Analyze() (string, error) {
	log := b.log.With().Str("op", "Analyze").Logger()
	defer b.History.Close()

	phrase, err := b.History.GetLastSearch()
	if err != nil {
		return "", err
	}

	log.Info().Msgf("last searched phrase: %s", phrase)
	return phrase, nil
}
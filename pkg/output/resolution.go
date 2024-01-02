package output

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidResolution = errors.New("invalid resolution")
)

type Resolution struct {
	Width  int
	Height int
}

func (r Resolution) String() string {
	return fmt.Sprintf("%dx%d", r.Width, r.Height)
}

func NewWithString(s string) (Resolution, error) {
	r := Resolution{}
	_, err := fmt.Sscanf(s, "%dx%d", &r.Width, &r.Height)
	if err != nil {
		return r, errors.Join(err, ErrInvalidResolution)
	}
	return r, nil
}

func MustStringResolution(s string) Resolution {
	r, err := NewWithString(s)
	if err != nil {
		panic(err)
	}
	return r
}

func NewResolution(w, h int) Resolution {
	return Resolution{w, h}
}

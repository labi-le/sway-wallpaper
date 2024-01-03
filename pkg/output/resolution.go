package output

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrInvalidResolution = errors.New("invalid resolution")
)

type Resolution struct {
	Width  int
	Height int
}

func (r *Resolution) String() string {
	return fmt.Sprintf("%dx%d", r.Width, r.Height)
}

func (r *Resolution) Set(s string) error {
	if regexp.MustCompile(`\d+x\d+`).MatchString(s) {
		_, err := fmt.Sscanf(s, "%dx%d", &r.Width, &r.Height)
		if r.Width <= 0 || r.Height <= 0 {
			return ErrInvalidResolution
		}
		return err
	}

	return ErrInvalidResolution
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

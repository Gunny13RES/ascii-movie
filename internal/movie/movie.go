package movie

import (
	"time"
)

func NewMovie() Movie {
	return Movie{}
}

type Movie struct {
	Filename string
	Cap      int
	Frames   []Frame
	Width    int
}

func (m Movie) Duration() time.Duration {
	var totalDuration time.Duration
	for _, f := range m.Frames {
		totalDuration += f.Duration
	}
	return totalDuration
}

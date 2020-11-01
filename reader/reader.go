package reader

import (
	"io"
	"time"
)

type Line struct {
	Text string
	Time time.Time
}

func NewLine(text string) Line {
	return Line{text, time.Now()}
}


type Reader interface {
	io.Closer
	Next() (Line, error)
}


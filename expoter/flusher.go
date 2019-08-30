package expoter

import "io"

type Flusher interface {
	io.Writer
	Flush() error
}
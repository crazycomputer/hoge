package exporter

import "io"

type Flusher interface {
	io.Writer
	Flush() error
}
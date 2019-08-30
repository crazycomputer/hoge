package serializer

import (
	"github.com/crazycomputer/hoge/buffer"
	"github.com/crazycomputer/hoge/log"
)

type Serializer interface {
	SerializeLog(log.Log) (*buffer.Buffer, error)
}

package serializer

import (
	"github.com/crazycomputer/hoge/buffer"
	"github.com/crazycomputer/hoge/log"
	"sync"
)

var (
	_serialPool = sync.Pool{New: func() interface{} {
		return &commonSerializer{}
	}}

)

const form = "2006-01-01 15:04:05"
type commonSerializer struct {
	buf *buffer.Buffer
}

func getCommonSerializer() *commonSerializer {
	return _serialPool.Get().(*commonSerializer)
}

func (cs *commonSerializer) clone() *commonSerializer {
	clone := getCommonSerializer()
	clone.buf = buffer.Get()
	return clone
}

func (cs *commonSerializer) SerializeLog(l log.Log) (*buffer.Buffer, error) {
	final:=cs.clone()
	final.buf.AppendByte('[')
	final.buf.AppendString(l.Level.String())
	final.buf.AppendString("][")
	final.buf.AppendString(l.Time.Format(form))
	final.buf.AppendString("][")
	if l.Caller.Defined{
		final.buf.AppendString("go-")
		final.buf.AppendInt(int64(l.Caller.Go))
		final.buf.AppendString("][")
		final.buf.AppendString(l.Caller.TrimmedPath())
		final.buf.AppendString("][")
		//todo 函数名
		final.buf.AppendString("]\n")
	}
	final.buf.AppendString(l.Message)
	if l.Stack != "" {
		final.buf.AppendString(l.Stack)

	}
	return final.buf,nil
}

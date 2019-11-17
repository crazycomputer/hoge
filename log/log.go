package log

import (
	"github.com/crazycomputer/hoge/buffer"
	"github.com/crazycomputer/hoge/exporter"
	"github.com/crazycomputer/hoge/level"
	"strings"
	"sync"
	"time"
)

var _wlPool = sync.Pool{New: func() interface{} {
	// Pre-allocate some space for cores.
	return &WrapLog{
		exporter: make([]exporter.Flusher, 4),
	}
}}

type Log struct {
	Level   level.Level
	Time    time.Time
	Message string
	Caller  Caller
	Stack   string
}

type Caller struct {
	Defined bool
	PC      uintptr
	File    string
	Line    int
	Go      int
}

type WrapLog struct {
	Log
	dirty    bool // best-effort detection of pool misuse
	exporter []exporter.Flusher
}

func (wl *WrapLog) Wrap(l Log, out exporter.Flusher) *WrapLog {
	if wl == nil {
		wl = getWrapLog()
		wl.Log = l
	}
	wl.exporter = append(wl.exporter, out)
	return wl
}

func getWrapLog() *WrapLog {
	ce := _wlPool.Get().(*WrapLog)
	ce.reset()
	return ce
}

func (ec Caller) TrimmedPath() string {
	idx := strings.LastIndexByte(ec.File, '/')
	if idx == -1 {
		return ec.FullPath()
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(ec.File[:idx], '/')
	if idx == -1 {
		return ec.FullPath()
	}
	buf := buffer.Get()
	// Keep everything after the penultimate separator.
	buf.AppendString(ec.File[idx+1:])
	buf.AppendByte(':')
	buf.AppendInt(int64(ec.Line))
	caller := buf.String()
	buf.Free()
	return caller
}

func (ec Caller) FullPath() string {
	buf := buffer.Get()
	buf.AppendString(ec.File)
	buf.AppendByte(':')
	buf.AppendInt(int64(ec.Line))
	caller := buf.String()
	buf.Free()
	return caller
}

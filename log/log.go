package log

import (
	"github.com/crazycomputer/hoge/level"
	"time"
)

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

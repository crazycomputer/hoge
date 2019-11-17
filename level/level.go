package level

import (
	"github.com/crazycomputer/hoge/atomic"
	"strconv"
)

type Level int8

const (
	TraceLevel Level = iota - 2

	DebugLevel Level = iota - 1

	InfoLevel

	WarnLevel

	ErrorLevel

	_minLevel = DebugLevel
	_maxLevel = ErrorLevel
)

const TRACE = "TRACE"
const DEBUG = "DEBUG"
const INFO = "INFO"
const WARN = "WARN"
const ERROR = "ERROR"

type Checker struct {
	l *atomic.Int32
}

func NewChecker() Checker {
	return Checker{
		l: atomic.NewInt32(int32(InfoLevel)),
	}
}

func (c Checker) SetLevel(l Level) {
	c.l.Store(int32(l))
}

func NewLevelChecker(l Level) Checker {
	c := NewChecker()
	c.SetLevel(l)
	return c
}

func (c Checker) Enabled(l Level) bool {
	return c.Level().Enabled(l)
}

func (c Checker) Level() Level {
	return Level(int8(c.l.Load()))
}

func (lvl Level) Enabled(l Level) bool {
	return lvl >= l
}

type Enabler interface {
	Enabled(Level) bool
}

func (lvl Level) String() string {
	switch lvl {
	case TraceLevel:
		return TRACE
	case DebugLevel:
		return DEBUG
	case InfoLevel:
		return INFO
	case WarnLevel:
		return WARN
	case ErrorLevel:
		return ERROR
	default:
		return "unknown" + strconv.FormatInt(int64(lvl), 10)
	}

}

package hoge

import (
	"github.com/crazycomputer/hoge/expoter"
	"github.com/crazycomputer/hoge/level"
	"github.com/crazycomputer/hoge/serializer"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	level.Enabler
	s   serializer.Serializer
	out expoter.Flusher
	addCaller bool
	addStack  zapcore.LevelEnabler
	callerSkip int
}

type Core struct {
	level.Enabler
	s   serializer.Serializer
	out expoter.Flusher
}

package hoge

import (
	"fmt"
	"github.com/crazycomputer/hoge/exporter"
	"github.com/crazycomputer/hoge/expoter"
	"github.com/crazycomputer/hoge/level"
	"github.com/crazycomputer/hoge/log"
	"github.com/crazycomputer/hoge/serializer"
	"go.uber.org/zap/zapcore"
	"runtime"
	"time"
)

type Logger struct {
	level.Enabler
	s          serializer.Serializer
	out        exporter.Flusher
	addCaller  bool
	addStack   zapcore.LevelEnabler
	callerSkip int
}

func (logger *Logger) clone() *Logger {
	log := *logger
	return &log
}

func (logger *Logger) check(lvl level.Level, msg string) *zapcore.CheckedEntry {
	const callerSkipOffset = 2

	ent := log.Log{
		Time:       time.Now(),
		Level:      lvl,
		Message:    msg,
	}
	ce := log.wrap(ent, nil)
	willWrite := ce != nil

	// Set up any required terminal behavior.
	switch ent.Level {
	case zapcore.PanicLevel:
		ce = ce.Should(ent, zapcore.WriteThenPanic)
	case zapcore.FatalLevel:
		ce = ce.Should(ent, zapcore.WriteThenFatal)
	case zapcore.DPanicLevel:
		if log.development {
			ce = ce.Should(ent, zapcore.WriteThenPanic)
		}
	}

	// Only do further annotation if we're going to write this message; checked
	// entries that exist only for terminal behavior don't benefit from
	// annotation.
	if !willWrite {
		return ce
	}

	// Thread the error output through to the CheckedEntry.
	ce.ErrorOutput = log.errorOutput
	if log.addCaller {
		ce.Entry.Caller = zapcore.NewEntryCaller(runtime.Caller(log.callerSkip + callerSkipOffset))
		if !ce.Entry.Caller.Defined {
			fmt.Fprintf(log.errorOutput, "%v Logger.check error: failed to get caller\n", time.Now().UTC())
			log.errorOutput.Sync()
		}
	}
	if log.addStack.Enabled(ce.Entry.Level) {
		ce.Entry.Stack = Stack("").String
	}

	return ce
}


func (logger *Logger) wrap(l log.Log, wl *log.WrapLog) *log.WrapLog {
	if logger.Enabled(l.Level) {
		return wl.Wrap(l, logger.out)
	}
	return wl
}

func (logger *Logger) Trace(args ...interface{}) {
	logger.log(DebugLevel, "", args, nil)
}


func (logger *Logger) Debug(args ...interface{}) {
	logger.log(DebugLevel, "", args, nil)
}


func (logger *Logger) Info(args ...interface{}) {
	logger.log(InfoLevel, "", args, nil)
}


func (logger *Logger) Warn(args ...interface{}) {
	logger.log(WarnLevel, "", args, nil)
}


func (logger *Logger) Error(args ...interface{}) {
	logger.log(ErrorLevel, "", args, nil)
}


func (logger *Logger) log(lvl level.Level, template string, fmtArgs []interface{}, context []interface{}) {



	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}

	if ce := logger.Check(lvl, msg); ce != nil {
		ce.Write(s.sweetenFields(context)...)
	}
}


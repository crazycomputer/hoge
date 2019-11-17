package hoge

import (
	"github.com/crazycomputer/hoge/level"
	"go.uber.org/zap/zapcore"
	"os"

	"time"
)

type FlowCtrlCfg struct {
	Period    time.Duration `json:"period" yaml:"period"`
	Threshold int           `json:"threshold" yaml:"threshold"`
	OverFreq  int           `json:"frequency" yaml:"frequency"`
}
type ExporterConfig struct {
	Mode string `json:"mode" yaml:"mode"`
	Path string `json:"path" yaml:"path"`
}
type Config struct {
	Level level.Checker `json:"level" yaml:"level"`

	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`

	FlowControl *FlowCtrlCfg `json:"flowControl" yaml:"flowControl"`

	OutputExporters []ExporterConfig `json:"output" yaml:"output"`
}

func (cfg Config) Build(opts ...Option) (*Logger, error) {

	sink, errSink, err := cfg.openSinks()
	if err != nil {
		return nil, err
	}

	log := new(
		newCore(enc, sink, cfg.Level),
		cfg.buildOptions(errSink)...,
	)
	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}
	return log, nil
}


func new(core zapcore.Core, options ...Option) *Logger {
	log := &Logger{
		core:        core,
		errorOutput: zapcore.Lock(os.Stderr),
		addStack:    zapcore.FatalLevel + 1,
	}
	return log.WithOptions(options...)
}

type Option interface {
	apply(*Logger)
}

type optionFunc func(*Logger)

package hoge

import (
	"github.com/crazycomputer/hoge/level"
	"go.uber.org/zap/zapcore"
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

	Serializer string `json:"serializer" yaml:"serializer"`

	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`

	OutputExporters []ExporterConfig `json:"output" yaml:"output"`
}

func (cfg Config) Build(opts ...Option) (*Logger, error) {
	enc, err := cfg.buildEncoder()
	if err != nil {
		return nil, err
	}


	sink, errSink, err := cfg.openSinks()
	if err != nil {
		return nil, err
	}

	log := New(
		zapcore.NewCore(enc, sink, cfg.Level),
		cfg.buildOptions(errSink)...,
	)
	if len(opts) > 0 {
		log = log.WithOptions(opts...)
	}
	return log, nil
}

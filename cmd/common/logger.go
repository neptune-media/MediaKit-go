package common

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() *zap.Logger {
	loggerCfg := NewLoggerConfig()
	logger, err := loggerCfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func NewLoggerConfig() *zap.Config {
	level, err := zap.ParseAtomicLevel(viper.GetString(ArgLogLevel))
	if err != nil {
		panic(err)
	}

	cfg := &zap.Config{
		Development:      false,
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		ErrorOutputPaths: []string{"stderr"},
		Level:            level,
		OutputPaths:      []string{"stdout"},
	}

	switch viper.GetString(ArgLogFormat) {
	case "json":
		cfg.Encoding = "json"
	default:
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	}

	return cfg
}

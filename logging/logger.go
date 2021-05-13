package logging

import (
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the global use log instance
var Logger *zap.Logger

// InitLogger is intended to init logger
func InitLogger() {
	var err error
	var c zap.Config
	if config.Debug {
		c = zap.NewDevelopmentConfig()
	} else {
		c = zap.NewProductionConfig()
	}
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	c.OutputPaths = []string{"stdout"}
	c.ErrorOutputPaths = []string{"stderr"}
	Logger, err = c.Build()
	if err != nil {
		panic(errors.WithMessage(err, "can't initialize logger driver, program exited."))
	}
}

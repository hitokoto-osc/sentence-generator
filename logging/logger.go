package logging

import (
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Logger is the global use log instance
var Logger *zap.Logger

// InitLogger is intended to init logger
func InitLogger() {
	var err error
	if config.Debug {
		Logger, err = zap.NewDevelopment()
	} else {
		Logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(errors.WithMessage(err, "can't initialize logger driver, program exited."))
	}
}

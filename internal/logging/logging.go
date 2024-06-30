package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger() (*zap.SugaredLogger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

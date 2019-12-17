package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
)

// Log returns a configured logger
func Log() *zap.SugaredLogger {
	debug := os.Getenv("DEBUG")

	var logger *zap.Logger
	logger, err := zap.NewProduction()

	if debug == "true" {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Printf("could not init logger: %v", err)
	}

	sugar := logger.Sugar()

	return sugar
}

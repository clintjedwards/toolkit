package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitGlobalLogger sets an application wide logger using zap
// after running this in main you can use Zap.S() to log things
func InitGlobalLogger(level string, development bool) {

	config := zap.NewProductionConfig()

	if development {
		config = zap.NewDevelopmentConfig()
	}

	config.Level.SetLevel(parseLogLevel(level))

	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}

	zap.ReplaceGlobals(logger)

}

func parseLogLevel(loglevel string) zapcore.Level {
	switch loglevel {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

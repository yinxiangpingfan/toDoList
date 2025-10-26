package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func LoggerInit() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.ErrorLevel)
	logger := zap.New(core)
	Logger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.OpenFile("../../logs/test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	return zapcore.AddSync(file)
}

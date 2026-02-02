package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	// 1. 配置zap的编码器配置（保持你原有配置不变）
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoder := zapcore.NewJSONEncoder(encodeConfig)
	level := zapcore.InfoLevel
	logFilePath := "./app.log"
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("打开日志文件失败：" + err.Error())
	}
	stdoutWriteSyncer := zapcore.AddSync(os.Stdout)
	fileWriteSyncer := zapcore.AddSync(file)
	multiWriteSyncer := zapcore.NewMultiWriteSyncer(stdoutWriteSyncer, fileWriteSyncer)
	core := zapcore.NewCore(
		encoder,
		multiWriteSyncer,
		level,
	)

	Logger = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}
func Info(msg string, field ...zap.Field) {
	Logger.Info(msg, field...)
}
func Error(msg string, field ...zap.Field) {
	Logger.Error(msg, field...)
}

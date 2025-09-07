package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: Replace logger with slog as that will not force your to include zap.Fields which is kinda ugly

// Gl is the global logger
var Gl *zap.Logger

func init() {
	Gl = zap.NewNop()
}

func Initialize() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, err := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	Gl = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func Error(op, msg string, err error, fields ...zap.Field) {
	errF := zap.Error(err)
	opF := zap.String("op", op)
	if len(fields) == 0 {
		Gl.Error(msg, errF, opF)
	} else {
		Gl.Error(
			msg,
			append([]zap.Field{errF, opF}, fields...)...,
		)
	}
}

func Fatal(op, msg string, err error, fields ...zap.Field) {
	errF := zap.Error(err)
	opF := zap.String("op", op)
	if len(fields) == 0 {
		Gl.Error(msg, errF, opF)
	} else {
		Gl.Error(
			msg,
			append([]zap.Field{errF, opF}, fields...)...,
		)
	}
}

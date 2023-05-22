package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var core zapcore.Core

var Logger *zap.Logger

var size uint16

func init() {
	// 设置一些基本日志格式
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.999"))
		},
		CallerKey: "caller",
		//EncodeCaller: zapcore.FullCallerEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现两个判断日志等级的interface
	printLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	// 创建Logger
	core = zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), printLevel), //打印到控制台
	)
	//分配log
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // 显示打日志点的文件名和行数
}

type LogWriter struct {
}

const (
	line          = ""
	ColorRed      = line + "\033[31m"
	ColorGreen    = line + "\033[32m"
	ColorYellow   = line + "\033[33m"
	ColorBlue     = line + "\033[34m"
	ColorPurple   = line + "\033[35m"
	ColorWhite    = line + "\033[37m"
	ColorHiRed    = line + "\033[91m"
	ColorHiGreen  = line + "\033[92m"
	ColorHiYellow = line + "\033[93m"
	ColorHiBlue   = line + "\033[94m"
	ColorHiPurple = line + "\033[95m"
	ColorHiWhite  = line + "\033[97m"
	ColorReset    = line + "\033[0m"
)

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

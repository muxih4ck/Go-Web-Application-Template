// zap是uber开发的一个高性能的日志包
// lumberjack是用于持久化日志的包
package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
)

type Logger interface {
	Info(string, ...zap.Field)
	Fatal(string, ...zap.Field)
	Debug(string, ...zap.Field)
	Error(string, ...zap.Field)
}

func init() {
	hook := lumberjack.Logger{
		Filename:   "./logs/api_server1.log", // 日志文件路径
		MaxSize:    128,                      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 5,                        // 日志文件最多保存多少个备份
		MaxAge:     30,                       // 文件最多保存多少天
		Compress:   true,                     // 是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)
	// 新建一个ZapCore
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "serviceName"))
	// 构造日志
	logger = zap.New(core, caller, development, filed)
}

func SyncLogger() {
	logger.Sync()
}

// info level log
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// fatal level log
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// debug level log
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// error level log
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

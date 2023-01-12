// Package logger 处理日志相关逻辑
package logger

import (
	"fmt"
	"go-api-demo/pkg/app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"time"
)

var Logger *zap.Logger

// InitLogger 日志初始化
func InitLogger(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string, level string) {
	//1. 获取日志写入介质
	writeSyncer := getLogWriter(filename, maxSize, maxBackup, maxAge, compress, logType)

	//2. 设置日志等级， 具体请见 config/log.go文件
	logLevel := new(zapcore.Level)
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		fmt.Println("日志初始化错误，日志级别设置有误。请修改config/log.go文件中的log.level配置项")
	}

	//3. 初始化core
	core := zapcore.NewCore(getEncoder(), writeSyncer, logLevel)

	//4. 初始化Logger
	Logger = zap.New(core,
		zap.AddCaller(),                   //调用文件和行号，内部使用runtime.Caller
		zap.AddCallerSkip(1),              //封装了一层，调用文件去除一层(runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), //Error时才会显示 stacktrace
	)

	//5. 将自定义的logger替换为全局的logger，zap.L().Fatal()调用时，就会使用我们自定的Logger
	zap.ReplaceGlobals(Logger)
}

// getEncoder 设置日志存储格式
func getEncoder() zapcore.Encoder {
	//1. 日志格式规则
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller", // 代码调用，如 paginator/paginator.go:148
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,      // 每行日志的结尾添加 "\n"
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 日志级别名称大写，如 ERROR、INFO
		EncodeTime:     customTimeEncoder,              // 时间格式，我们自定义为 2006-01-02 15:04:05
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行时间，以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Caller 短格式，如：types/converter.go:17，长格式为绝对路径

	}

	//2. 本地环境配置
	if app.IsLocal() {
		//终端输出的关键词高亮
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		//本地设置内置的Console 解码器 支持stacktrace换行
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	//3. 线上环境使用 JSON 编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// customTimeEncoder 自定义友好的时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// getLogWriter 日志记录介质。 使用两种介质：os.Stdout 和 文件
func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string) zapcore.WriteSyncer {
	//1. 如果配置了按照日期记录日志文件
	if logType == "daily" {
		logName := time.Now().Format("2006-01-02.log")
		filename = strings.ReplaceAll(filename, "logs.log", logName)
	}

	//2. 滚动日志，详见 config/log.go
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	//3. 配置输出介质，本地开发终端打印和记录文件
	if app.IsLocal() {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	//4. 生产环境只记录文件
	return zapcore.AddSync(lumberJackLogger)
}

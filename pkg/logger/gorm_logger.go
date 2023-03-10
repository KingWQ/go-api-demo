package logger

import (
	"context"
	"errors"
	"go-api-demo/pkg/helpers"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// GormLogger 操作对象，实现 gormLogger.Interface
type GormLogger struct {
	ZapLogger     *zap.Logger
	SlowThreshold time.Duration
}

// NewGormLogger 外部调用。实例化一个 GormLogger 对象，示例：
//     DB, err := gorm.Open(dbConfig, &gorm.Config{
//         Logger: logger.NewGormLogger(),
//     })
func NewGormLogger() GormLogger {
	return GormLogger{
		ZapLogger:     Logger,
		SlowThreshold: 200 * time.Millisecond,
	}
}

// LogMode 实现 gormLogger.Interface 的 LogMode 方法
func (l GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return GormLogger{
		ZapLogger:     l.ZapLogger,
		SlowThreshold: l.SlowThreshold,
	}
}

// Info 实现 gormLogger.Interface 的 Info 方法
func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Debugf(str, args...)
}

// Warn 实现 gormLogger.Interface 的 Warn 方法
func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Warnf(str, args...)
}

// Error 实现 gormLogger.Interface 的 Error 方法
func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Errorf(str, args...)
}

// Trace 实现 gormLogger.Interface 的 Trace 方法
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	//1. 获取运行时间
	elapsed := time.Since(begin)

	//2. 获取SQL 请求和返回条数
	sql, rows := fc()

	//3. 通用字段
	logFields := []zap.Field{
		zap.String("sql", sql),
		zap.String("time", helpers.MicrosecondsStr(elapsed)),
		zap.Int64("rows", rows),
	}

	//4. Gorm 错误
	if err != nil {
		//记录未找到的错误使用 warning 等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.logger().Warn("Database ErrRecordNotFound", logFields...)
		} else {
			logFields = append(logFields, zap.Error(err))
			l.logger().Error("Database Error", logFields...)
		}
	}

	//5. 慢查询日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger().Warn("Database Slow log", logFields...)
	}

	//6. 记录所有的SQL请求
	l.logger().Debug("Database Query", logFields...)
}

func (l GormLogger) logger() *zap.Logger {
	// 跳过 gorm 内置的调用
	var (
		gormPackage    = filepath.Join("gorm.io", "gorm")
		zapgormPackage = filepath.Join("moul.io", "zapgorm2")
	)

	// 减去一次封装，以及一次在 logger 初始化里添加 zap.AddCallerSkip(1)
	clone := l.ZapLogger.WithOptions(zap.AddCallerSkip(-2))

	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapgormPackage):
		default:
			// 返回一个附带跳过行号的新的 zap logger
			return clone.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return l.ZapLogger

}

package logger

import (
	"context"
	"errors"
	"gohub/pkg/helpers"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

//GormLogger 操作对象 gormlogger.Interface
type GormLogger struct {
	ZapLogger    *zap.Logger
	SlowThrshold time.Duration
}

//NewGormLogger 外部调用 实例化一个GormLogger 对象 示例
//     DB, err := gorm.Open(dbConfig, &gorm.Config{
//         Logger: logger.NewGormLogger(),
//     })
func NewGormLogger() GormLogger {
	return GormLogger{
		ZapLogger:    Logger,                 // 使用全局的 logger.Logger 对象
		SlowThrshold: 200 * time.Microsecond, // 慢查询阈值，单位为千分之一秒
	}
}

//LogMod 实现gormLogger.Interface 的LogMod 方法
func (l GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return GormLogger{
		ZapLogger:    l.ZapLogger,
		SlowThrshold: l.SlowThrshold,
	}
}

// Info 实现 gormlogger.Interface 的 Info 方法
func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Debugf(str, args...)
}

// Warn 实现 gormlogger.Interface 的 Warn 方法
func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Warnf(str, args...)
}

// Error 实现 gormlogger.Interface 的 Error 方法
func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Errorf(str, args...)
}

// Trace 实现gormlogger.Interface 的Trace方法
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	//获取运行时间
	elapsed := time.Since(begin)
	//获取SQL 请求和返回条数
	sql, rows := fc()

	//通用字段
	logFields := []zap.Field{
		zap.String("sql", sql),
		zap.String("time", helpers.MicrosecondsStr(elapsed)),
		zap.Int64("rows", rows),
	}

	//Gorm错误
	if err != nil {

		//记录未找到的错误 使用warning 等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.logger().Warn("Database ErrRecordNotFound", logFields...)
		} else {
			//其他错误使用error 等级
			logFields = append(logFields, zap.Error(err))

			l.logger().Error("Database Error", logFields...)
		}
	}

	//慢查询日志
	if l.SlowThrshold != 0 && elapsed > l.SlowThrshold {
		l.logger().Warn("Database Slow Log", logFields...)
	}

	//记录所有 SQL请求
	l.logger().Debug("Database Query", logFields...)

}

func (l GormLogger) logger() *zap.Logger {
	//跳过gorm 内置的调用
	var (
		gormPackage    = filepath.Join("gorm.io", "gorm")
		zapgormPackage = filepath.Join("moul.io", "zapgorm2")
	)

	//TODO 暂未理解这里的调用
	//减去一次封装，以及一次在 logger初始化里添加 zap.AddCallerSkip(1)
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

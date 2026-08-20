package database

import (
	"context"
	"errors"
	"time"

	applogger "github.com/akozadaev/go_todo_service/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type correlatedGORMLogger struct {
	level         logger.LogLevel
	slowThreshold time.Duration
}

func newCorrelatedGORMLogger() logger.Interface {
	return &correlatedGORMLogger{
		level:         logger.Info,
		slowThreshold: 200 * time.Millisecond,
	}
}

func (l *correlatedGORMLogger) LogMode(level logger.LogLevel) logger.Interface {
	clone := *l
	clone.level = level
	return &clone
}

func (l *correlatedGORMLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= logger.Info {
		applogger.FromContext(ctx).Sugar().Infof(msg, args...)
	}
}

func (l *correlatedGORMLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= logger.Warn {
		applogger.FromContext(ctx).Sugar().Warnf(msg, args...)
	}
}

func (l *correlatedGORMLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= logger.Error {
		applogger.FromContext(ctx).Sugar().Errorf(msg, args...)
	}
}

func (l *correlatedGORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level == logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.String("component", "gorm"),
		zap.Duration("duration", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}
	log := applogger.FromContext(ctx)

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && l.level >= logger.Error:
		log.Error("SQL query failed", append(fields, zap.Error(err))...)
	case elapsed > l.slowThreshold && l.level >= logger.Warn:
		log.Warn("Slow SQL query", append(fields, zap.Duration("slow_threshold", l.slowThreshold))...)
	case l.level >= logger.Info:
		log.Info("SQL query", fields...)
	}
}

// Package db 直接替换掉gorm中的日志记录，实现自带的日志记录和等级
package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/yhhaiua/engine/logzap"
	gorm_log "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

var replaceLog = NewLog(gorm_log.Config{
	SlowThreshold: 200 * time.Millisecond,
	LogLevel:      gorm_log.Warn,
	Colorful:      true,
})

type accessorLog struct {
	gorm_log.Config
	LogLevel                            gorm_log.LogLevel
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewLog(config gorm_log.Config) gorm_log.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	return &accessorLog{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}
func (l *accessorLog) LogMode(level gorm_log.LogLevel) gorm_log.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}
func (l *accessorLog) Info(ctx context.Context, msg string, data ...interface{}) {
	logzap.LOGGER("LOG").Sugar().Infof(msg, data)
}
func (l *accessorLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	logzap.LOGGER("LOG").Sugar().Warnf(msg, data)
}
func (l *accessorLog) Error(ctx context.Context, msg string, data ...interface{}) {
	logzap.LOGGER("LOG").Sugar().Errorf(msg, data)
}
func (l *accessorLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gorm_log.Error && (!errors.Is(err, gorm_log.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logzap.LOGGER("LOG").Sugar().Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logzap.LOGGER("LOG").Sugar().Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gorm_log.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logzap.LOGGER("LOG").Sugar().Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logzap.LOGGER("LOG").Sugar().Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == gorm_log.Info:
		sql, rows := fc()
		if rows == -1 {
			logzap.LOGGER("LOG").Sugar().Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logzap.LOGGER("LOG").Sugar().Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

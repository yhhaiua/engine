package log

import (
	"engine/logzap"
	"fmt"
	"go.uber.org/zap/zapcore"
	"runtime"
	"sync"
)

type Logger struct {
}

var (
	globalLog *Logger
	once      sync.Once
)

func GetLogger() HandlerLog {
	once.Do(func() {
		globalLog = &Logger{}
	})
	return globalLog
}

func (l *Logger) Debugf(arg0 string, args ...interface{}) {
	logzap.LOGGERSUGER("LOG").Debugf(arg0, args...)
}

func (l *Logger) Infof(arg0 string, args ...interface{}) {
	logzap.LOGGERSUGER("LOG").Infof(arg0, args...)
}
func (l *Logger) Warnf(arg0 string, args ...interface{}) {
	logzap.LOGGERSUGER("LOG").Warnf(arg0, args...)
}
func (l *Logger) Errorf(arg0 string, args ...interface{}) {
	logzap.LOGGERSUGER("LOG").Errorf(arg0, args...)
}

//Config 读取日志配置文件
func (l *Logger) LoadConfig(homeRote, dir string) {
	logzap.LoadConfig(homeRote, dir)
}

//InfoLog 特殊日志记录文件 name 文件标识，arg0 (func() string)
func (l *Logger) InfoLog(name string, arg0 string) {
	logzap.LOGGER(name).Info(arg0)
}

//TraceErr 输出错误，跟踪代码
func (l *Logger) TraceErr(args ...interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	data := ""
	for _, v := range args {
		data += fmt.Sprintf("%v", v)
	}
	data += "\n"
	data += string(buf[:n])
	logzap.LOGGER("LOG").Error(data)
}

func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	logzap.LOGGER("LOG").Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	logzap.LOGGER("LOG").Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	logzap.LOGGER("LOG").Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	logzap.LOGGER("LOG").Error(msg, fields...)
}

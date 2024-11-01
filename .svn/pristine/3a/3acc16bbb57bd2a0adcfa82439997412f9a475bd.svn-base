package log

import "go.uber.org/zap/zapcore"

type HandlerLog interface {

	//普通日志输出
	Debugf(arg0 string, args ...interface{})
	Infof(arg0 string, args ...interface{})
	Warnf(arg0 string, args ...interface{})
	Errorf(arg0 string, args ...interface{})

	//结构日志输出
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)

	//TraceErr 输出错误，跟踪代码
	TraceErr(args ...interface{})
	//InfoLog 特殊日志记录文件 name 文件标识，arg0 string)
	InfoLog(name string, arg0 string)
	//Config 读取日志配置文件
	LoadConfig(homeRote, dir string)
}

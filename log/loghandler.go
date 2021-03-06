package log

type LogHandler interface {
	Debug(arg0 interface{}, args ...interface{})
	Info(arg0 interface{}, args ...interface{})
	Warn(arg0 interface{}, args ...interface{})
	Error(arg0 interface{}, args ...interface{})
	InfoLog(name string,arg interface{})
	Config(dir string)
}


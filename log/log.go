package log

import "github.com/yhhaiua/log4go"

type Logger struct {

}
var globalLog = &Logger{}
func GetLogger() LogHandler {
	return globalLog
}

func (l *Logger)Info(arg0 interface{}, args ...interface{}){
	log4go.Info(arg0,args...)
}
func (l *Logger)Warn(arg0 interface{}, args ...interface{}){
	log4go.Warn(arg0,args...)
}
func (l *Logger)Error(arg0 interface{}, args ...interface{}){
	log4go.Error(arg0,args...)
}
func (l *Logger)InfoLog(name string,arg interface{}){
	log4go.InfoLog(name,arg)
}

package jsonx

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/yhhaiua/engine/log"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var logger = log.GetLogger()

//MarshalToString 结构到字符串
func MarshalToString(v interface{}) string {
	s, err := json.MarshalToString(v)
	if err != nil {
		logger.TraceErr(err)
	}
	return s
}

//Marshal 结构到字符串
func Marshal(v interface{}) []byte {
	s, err := json.Marshal(v)
	if err != nil {
		logger.TraceErr(err)
	}
	return s
}

//UnmarshalFromString 解析json字符串到需要的结构
func UnmarshalFromString(str string, v interface{}) bool {
	err := json.UnmarshalFromString(str, v)
	if err != nil {
		logger.TraceErr(str, err)
		return false
	}
	return true
}

//Unmarshal 解析json字符串到需要的结构
func Unmarshal(data []byte, v interface{}) bool {
	err := json.Unmarshal(data, v)
	if err != nil {
		logger.TraceErr(err)
		return false
	}
	return true
}

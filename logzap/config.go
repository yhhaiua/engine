package logzap

import "go.uber.org/zap/zapcore"

type ZapConfig struct {
	Level         string 			`yaml:"level"`                           	// 级别
	Format        string 			`yaml:"format"`                        		// 输出
	Prefix        string 			`yaml:"prefix"`                       	 	// 日志前缀
	Director      string 			`yaml:"director"`                 			// 日志文件夹
	FileName      string 			`yaml:"file-name"`                			// 文件名称
	Pattern       string 			`yaml:"pattern"`                			// 时间后缀
	ShowLine      bool   			`yaml:"showLine"`                 			// 显示行
	EncodeLevel   string 			`yaml:"encode-level"`       				// 编码级
	StacktraceKey string 			`yaml:"stacktrace-key"` 					// 栈名
	LogInConsole  bool   			`yaml:"log-in-console"`  					// 输出控制台
	Category	  string			`yaml:"category"`							// 日志标签
	ZapLevel     zapcore.Level
}

type ZapConfigs struct {
	Zaps   		[]*ZapConfig			`yaml:"log"`  							// 所有的日志控制
}
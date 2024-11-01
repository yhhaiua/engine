package logzap

import (
	"engine/util"
	"fmt"
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math"
	"os"
	"path"
	"time"
)

func (z *ZapConfig) zap() (logger *zap.Logger) {
	if ok, _ := util.PathExists(z.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", z.Director)
		_ = os.Mkdir(z.Director, os.ModePerm)
	}
	//初始化等级
	z.zapLevel()

	if z.Category == "LOG" || z.Category == "stdout" {
		z.EncodeLevel = "LowercaseLevelEncoder"
		z.Prefix = "2006-01-02 15:04:05.000"
		z.StacktraceKey = "stacktrace"
	} else {
		z.EncodeLevel = ""
		z.Prefix = ""
		z.StacktraceKey = ""
	}
	z.Pattern = "%Y%m%d"
	z.Format = "format"
	if z.ZapLevel == zap.ErrorLevel {
		logger = zap.New(z.getEncoderCore(), zap.AddStacktrace(z.ZapLevel), zap.AddCallerSkip(1))
	} else {
		logger = zap.New(z.getEncoderCore(), zap.AddCallerSkip(1))
	}
	if z.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

func (z *ZapConfig) zapLevel() {
	var level zapcore.Level
	switch z.Level { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	z.ZapLevel = level
}

// getEncoderConfig 获取zapcore.EncoderConfig
func (z *ZapConfig) getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: z.StacktraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		//EncodeLevel:    zapcore.LowercaseLevelEncoder,
		//EncodeTime:     z.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case z.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case z.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case z.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = nil
	}
	if z.Prefix != "" {
		config.EncodeTime = z.CustomTimeEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func (z *ZapConfig) getEncoder() zapcore.Encoder {
	if z.Format == "json" {
		return zapcore.NewJSONEncoder(z.getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(z.getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func (z *ZapConfig) getEncoderCore() (core zapcore.Core) {
	writer, err := z.GetWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(z.getEncoder(), writer, z.ZapLevel)
}

// 自定义日志输出时间格式
func (z *ZapConfig) CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(z.Prefix))
}

func (z *ZapConfig) GetWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(z.Director, z.FileName+"."+z.Pattern),
		zaprotatelogs.WithRotationCount(math.MaxUint),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if z.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}

package logzap

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Loggers map[string]*zap.Logger
type LoggerSuger map[string]*zap.SugaredLogger

var (
	Global      Loggers
	GlobalSuger LoggerSuger
)

func init() {
	Global = make(Loggers)
	GlobalSuger = make(LoggerSuger)
	Global["stdout"], _ = zap.NewDevelopment(zap.AddCallerSkip(1))
	GlobalSuger["stdout"] = Global["stdout"].Sugar()
}

func LoadConfig(homeRote, file string) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not open %q for reading: %s\n", file, err)
		return
	}

	configs := &ZapConfigs{}
	if err := yaml.Unmarshal(yamlFile, configs); err != nil {
		fmt.Fprintf(os.Stderr, "LoadJsonConfiguration: Error: Could not parse json configuration in %q: %s\n", file, err)
		return
	}
	for _, c := range configs.Zaps {
		c.Director = homeRote + c.Director
		z := c.zap()
		Global[c.Category] = z
		GlobalSuger[c.Category] = z.Sugar()
	}
}

// LOGGER get the log Filter by category
func LOGGER(category string) *zap.Logger {
	f, ok := Global[category]
	if !ok {
		return Global["stdout"]
	}
	return f
}

// LOGGER get the log Filter by category
func LOGGERSUGER(category string) *zap.SugaredLogger {
	f, ok := GlobalSuger[category]
	if !ok {
		return GlobalSuger["stdout"]
	}
	return f
}

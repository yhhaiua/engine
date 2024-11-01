package conf

import (
	"github.com/yhhaiua/engine/log"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

var logger = log.GetLogger()

func LoadYamlConfig(file string, config interface{}) bool {

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Errorf("config error:%s", err.Error())
		return false
	}

	if err := yaml.Unmarshal(yamlFile, config); err != nil {
		logger.Errorf("config error :%s", err.Error())
		return false
	}
	return true
}

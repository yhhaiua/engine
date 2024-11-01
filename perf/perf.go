package perf

import (
	"engine/log"
	"net/http"
	_ "net/http/pprof"
)

var logger = log.GetLogger()

func Init(perfPort string) {
	if len(perfPort) != 0 {
		go startPerf(perfPort)
	}

}

func startPerf(perfPort string) {
	logger.Infof("性能分析工具启动...端口 %s", perfPort)
	err := http.ListenAndServe("localhost:"+perfPort, nil)
	if err != nil {
		logger.Errorf("性能分析工具启动失败:%s", err.Error())
	}
}

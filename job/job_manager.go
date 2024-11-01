package job

import (
	"github.com/panjf2000/ants/v2"
	"github.com/yhhaiua/engine/log"
)

var logger = log.GetLogger()

var defaultAntsPool, _ = ants.NewPool(ants.DefaultAntsPoolSize, ants.WithPanicHandler(func(p interface{}) {
	logger.TraceErr(p)
}))

func Submit(task func()) {
	_ = defaultAntsPool.Submit(task)
}

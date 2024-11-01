package job

import (
	"engine/log"
	"github.com/panjf2000/ants/v2"
)

var logger = log.GetLogger()

var defaultAntsPool, _ = ants.NewPool(ants.DefaultAntsPoolSize, ants.WithPanicHandler(func(p interface{}) {
	logger.TraceErr(p)
}))

func Submit(task func()) {
	_ = defaultAntsPool.Submit(task)
}

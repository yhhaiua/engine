package queue

import (
	"github.com/yhhaiua/engine/job"
	"github.com/yhhaiua/engine/util"
	"sync"
)

// DefaultQueue 消息队列处理 T为指针或接口
type DefaultQueue[T any] struct {
	sync.Mutex
	queue      []T
	temp       []T
	queueIndex int
	tempIndex  int
	state      util.AtomicInteger
	task       func(i T)
}

// NewDefault 新建默认队列
func NewDefault[T any](task func(i T)) *DefaultQueue[T] {
	q := new(DefaultQueue[T])
	q.queue = make([]T, 4)
	q.temp = make([]T, 4)
	q.task = task
	return q
}

// Add 向队列中添加数据
func (d *DefaultQueue[T]) Add(entry T) {
	d.Lock()
	if d.queueIndex+1 > len(d.queue) {
		d.queue = append(d.queue, entry)
		d.queueIndex++
	} else {
		d.queue[d.queueIndex] = entry
		d.queueIndex++
	}
	d.Unlock()
	if !d.state.CompareAndSet(0, 1) {
		return
	}
	job.Submit(d.Run)
}

// Run 新协程运行数据
func (d *DefaultQueue[T]) Run() {
	//==============================================
	d.Lock()
	if len(d.temp) < len(d.queue) {
		d.temp = make([]T, len(d.queue))
	}
	var t T
	for i := 0; i < d.queueIndex; i++ {
		d.temp[i] = d.queue[i]
		d.queue[i] = t
	}
	d.tempIndex = d.queueIndex
	d.queueIndex = 0
	d.Unlock()
	//==============================================
	for i := 0; i < d.tempIndex; i++ {
		d.task(d.temp[i])
		d.temp[i] = t
	}
	//==============================================
	d.Lock()
	if d.queueIndex > 0 {
		job.Submit(d.Run)
	} else {
		d.state.Set(0)
	}
	d.Unlock()
}

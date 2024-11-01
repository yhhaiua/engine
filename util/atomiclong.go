package util

import "sync/atomic"

type AtomicLong struct {
	value int64
}

func (a *AtomicLong) Set(newValue int64) {
	atomic.StoreInt64(&a.value, newValue)
}
func (a *AtomicLong) Get() int64 {
	return atomic.LoadInt64(&a.value)
}

func (a *AtomicLong) IncrementAndGet() int64 {
	return atomic.AddInt64(&a.value, 1)
}

func (a *AtomicLong) DecrementAndGet() int64 {
	return atomic.AddInt64(&a.value, -1)
}

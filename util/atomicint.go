package util

import "sync/atomic"

type AtomicInteger struct {
	value int32
}

func (a *AtomicInteger) Set(newValue int32) {
	atomic.StoreInt32(&a.value, newValue)
}

func (a *AtomicInteger) Get() int {
	return int(atomic.LoadInt32(&a.value))
}

func (a *AtomicInteger) CompareAndSet(old, new int32) bool {
	return atomic.CompareAndSwapInt32(&a.value, old, new)
}
func (a *AtomicInteger) IncrementAndGet() int32 {
	return atomic.AddInt32(&a.value, 1)
}

func (a *AtomicInteger) DecrementAndGet() int32 {
	return atomic.AddInt32(&a.value, -1)
}

func (a *AtomicInteger) Decrement(count int32) {
	atomic.AddInt32(&a.value, -count)
}

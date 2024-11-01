package concurrent

import (
	"github.com/yhhaiua/engine/util"
	"sync"
)

// ConcurrentHashMap 此方法弃用，后续使用HashMap
// Deprecated
type ConcurrentHashMap struct {
	hashMap sync.Map
	state   util.AtomicInteger
}

// NewConcurrentHashMap 此方法弃用，后续使用HashMap
// Deprecated
func NewConcurrentHashMap() *ConcurrentHashMap {
	m := &ConcurrentHashMap{}
	return m
}

func (h *ConcurrentHashMap) Put(key interface{}, value interface{}) {
	_, ok := h.hashMap.LoadOrStore(key, value)
	if !ok {
		h.state.IncrementAndGet()
	} else {
		h.hashMap.Store(key, value)
	}
}

func (h *ConcurrentHashMap) Get(key interface{}) (interface{}, bool) {
	return h.hashMap.Load(key)
}

func (h *ConcurrentHashMap) Remove(key interface{}) {
	_, ok := h.hashMap.LoadAndDelete(key)
	if ok {
		h.state.DecrementAndGet()
	}
}
func (h *ConcurrentHashMap) LoadAndRemove(key interface{}) (value interface{}, loaded bool) {
	v, ok := h.hashMap.LoadAndDelete(key)
	if ok {
		h.state.DecrementAndGet()
	}
	return v, ok
}

func (h *ConcurrentHashMap) PutIfAbsent(key interface{}, value interface{}) (actual interface{}, loaded bool) {
	v, ok := h.hashMap.LoadOrStore(key, value)
	if !ok {
		h.state.IncrementAndGet()
	}
	return v, ok
}

func (h *ConcurrentHashMap) Count() int {
	return h.state.Get()
}

func (h *ConcurrentHashMap) Values() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	h.hashMap.Range(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

func (h *ConcurrentHashMap) Range(f func(key, value any) bool) {
	h.hashMap.Range(f)
}
func (h *ConcurrentHashMap) Keys() []interface{} {
	var m []interface{}
	h.hashMap.Range(func(key, value interface{}) bool {
		m = append(m, key)
		return true
	})
	return m
}

func (h *ConcurrentHashMap) Clear() {
	m := h.Keys()
	if len(m) > 0 {
		for _, key := range m {
			h.hashMap.Delete(key)
		}
	}
	h.state.Decrement(int32(len(m)))
}

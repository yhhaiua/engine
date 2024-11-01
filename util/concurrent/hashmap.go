package concurrent

import (
	"github.com/yhhaiua/engine/util"
	"sync"
)

type HashKey interface {
	~int | ~int64 | ~int32 | ~int8 | ~uint | ~float32 | ~float64 | ~string
}
type HashMap[K HashKey, V any] struct {
	hashMap sync.Map
	state   util.AtomicInteger
}

func (h *HashMap[K, V]) Put(key K, value V) {
	_, ok := h.hashMap.LoadOrStore(key, value)
	if !ok {
		h.state.IncrementAndGet()
	} else {
		h.hashMap.Store(key, value)
	}
}

func (h *HashMap[K, V]) Get(key K) (V, bool) {
	v, ok := h.hashMap.Load(key)
	if ok {
		return v.(V), ok
	}
	var t V
	return t, ok
}

func (h *HashMap[K, V]) Remove(key K) {
	_, ok := h.hashMap.LoadAndDelete(key)
	if ok {
		h.state.DecrementAndGet()
	}
}
func (h *HashMap[K, V]) LoadAndRemove(key K) (value V, loaded bool) {
	v, ok := h.hashMap.LoadAndDelete(key)
	if ok {
		h.state.DecrementAndGet()
		return v.(V), ok
	}
	var t V
	return t, ok
}

func (h *HashMap[K, V]) PutIfAbsent(key K, value V) (actual V, loaded bool) {
	v, ok := h.hashMap.LoadOrStore(key, value)
	if !ok {
		h.state.IncrementAndGet()
	}
	return v.(V), ok
}

func (h *HashMap[K, V]) Count() int {
	return h.state.Get()
}

func (h *HashMap[K, V]) Values() map[K]V {
	m := make(map[K]V)
	h.hashMap.Range(func(key, value interface{}) bool {
		m[key.(K)] = value.(V)
		return true
	})
	return m
}

func (h *HashMap[K, V]) Range(f func(key, value any) bool) {
	h.hashMap.Range(f)
}
func (h *HashMap[K, V]) Keys() []K {
	var m []K
	h.hashMap.Range(func(key, value interface{}) bool {
		m = append(m, key.(K))
		return true
	})
	return m
}

func (h *HashMap[K, V]) ValueArray() []V {
	var m []V
	h.hashMap.Range(func(key, value interface{}) bool {
		m = append(m, value.(V))
		return true
	})
	return m
}

func (h *HashMap[K, V]) Clear() {
	m := h.Keys()
	if len(m) > 0 {
		for _, key := range m {
			h.hashMap.Delete(key)
		}
	}
	h.state.Decrement(int32(len(m)))
}

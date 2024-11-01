package cache

import (
	"sync"
)

type segment struct {
	sync.Map
	mutex sync.Mutex
}

func (s *segment) get(k Key) *entry {
	v, ok := s.Load(k)
	if ok {
		return v.(*entry)
	}
	return nil
}

func (s *segment) create(k Key, v Value, c *localCache) (*entry, *entry) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	en := s.get(k)
	var cen *entry = nil
	if en == nil {
		en = newEntry(k, v, sum(k))
		cen = c.cache.getOrSet(en)
		//createNewEntry = true
	} else if v != nil {
		en.setValue(v)
	}
	return en, cen
}

func (s *segment) createSimple(k Key, v Value, c *simpleCache) *entry {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	en := s.get(k)
	if en == nil {
		en = newEntry(k, v, sum(k))
		c.cache.getOrSet(en)
	} else if v != nil {
		en.setValue(v)
	}
	return en
}

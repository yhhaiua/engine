package cache

type simpleCache struct {
	cache  cache
	loader LoaderFunc
}

func (s *simpleCache) GetIfPresent(k Key) (Value, bool) {
	en := s.cache.get(k, sum(k))

	if en == nil || en.getValue() == nil {
		return nil, false
	}
	return en.getValue(), true
}

func (s *simpleCache) Put(k Key, v Value) {
	seg := s.cache.segment(sum(k))
	en := seg.get(k)
	if en == nil || en.getValue() == nil {
		en = seg.createSimple(k, v, s)
	} else {
		en.setValue(v)
	}
}

func (s *simpleCache) Invalidate(k Key) {
	s.cache.deleteKey(k)
}

func (s *simpleCache) InvalidateAll() {
	var m []Key
	s.cache.walk(func(en *entry) {
		m = append(m, en.key)
	})
	if len(m) > 0 {
		for _, k := range m {
			s.cache.deleteKey(k)
		}
	}
}

func (s *simpleCache) CacheAllValue() []Value {
	var m []Value
	s.cache.walk(func(en *entry) {
		m = append(m, en.getValue())
	})
	return m
}

func (s *simpleCache) Get(k Key) (Value, error) {
	seg := s.cache.segment(sum(k))
	en := seg.get(k)

	if en == nil || en.getValue() == nil {
		return s.load(k, seg)
	}
	return en.getValue(), nil
}

func (s *simpleCache) load(k Key, seg *segment) (Value, error) {
	if s.loader == nil {
		panic("cache loader function must be set")
	}

	en := seg.createSimple(k, nil, s)

	return s.create(en)
}

func (s *simpleCache) create(en *entry) (Value, error) {
	en.mutex.Lock()
	defer en.mutex.Unlock()

	if en.getValue() == nil {
		v, err := s.loader(en.key)
		if err != nil {
			return nil, err
		}
		en.setValue(v)
		return v, nil
	} else {
		return en.getValue(), nil
	}
}

func (s *simpleCache) Count() int {
	return s.cache.len()
}

func newSimpleCache() *simpleCache {
	return &simpleCache{
		cache: cache{},
	}
}

func NewSimpleCache(loader LoaderFunc) LoadingCache {
	c := newSimpleCache()
	c.loader = loader
	return c
}

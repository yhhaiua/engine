package cache

import (
	"github.com/yhhaiua/engine/log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// Default maximum number of cache entries.
	maximumCapacity = 1 << 30
	// Buffer size of entry channels
	chanBufSize = 64
	// Maximum number of entries to be drained in a single clean up.
	drainMax = 16
	// Number of cache access operations that will trigger clean up.
	drainThreshold = 64
)

// currentTime is an alias for time.Now, used for testing.
var currentTime = time.Now

var logger = log.GetLogger()

// localCache is an asynchronous LRU cache.
type localCache struct {
	// internal data structure
	cache cache // Must be aligned on 32-bit

	// user configurations
	expireAfterAccess time.Duration
	expireAfterWrite  time.Duration
	refreshAfterWrite time.Duration //TODO 不做刷新处理，此字段为0
	policyName        string

	onInsertion Func
	onRemoval   Func

	loader   LoaderFunc
	reloader Reloader
	stats    StatsCounter

	// cap is the cache capacity.
	cap int

	// accessQueue is the cache retention policy, which manages entries by access time.
	accessQueue policy
	// writeQueue is for managing entries by write time.
	// It is only fulfilled when expireAfterWrite or refreshAfterWrite is set.
	writeQueue policy
	// events is the cache event queue for processEntries
	events chan entryEvent

	// readCount is a counter of the number of reads since the last write.
	readCount int32

	// for closing routines created by this cache.
	closing int32
	closeWG sync.WaitGroup

	timer         *time.Ticker
	lastClearTime int64
}

// newLocalCache returns a default localCache.
// init must be called before this cache can be used.
func newLocalCache() *localCache {
	return &localCache{
		cap:   maximumCapacity,
		cache: cache{},
		stats: &nullCounter{}, //不需要统计数据时处理类
	}
}

// init initializes cache replacement policy after all user configuration properties are set.
func (c *localCache) init() {
	c.accessQueue = newPolicy(c.policyName)
	c.accessQueue.init(&c.cache, c.cap)
	if c.expireAfterWrite > 0 || c.refreshAfterWrite > 0 {
		c.writeQueue = &recencyQueue{}
	} else {
		c.writeQueue = discardingQueue{}
	}
	c.writeQueue.init(&c.cache, c.cap)
	c.events = make(chan entryEvent, chanBufSize)
	c.timer = time.NewTicker(10 * time.Minute)
	//c.closeWG.Add(1)
	go c.processEntries()
}

// Close implements io.Closer and always returns a nil error.
// Caller would ensure the cache is not being used (reading and writing) before closing.
func (c *localCache) Close() error {
	//if atomic.CompareAndSwapInt32(&c.closing, 0, 1) {
	//	// Do not close events channel to avoid panic when cache is still being used.
	//	//c.events <- entryEvent{nil, eventClose}
	//	// Wait for the goroutine to close this channel
	//	//c.closeWG.Wait()
	//}
	return nil
}

// Count 本地缓存中的数据
func (c *localCache) Count() int {
	return c.cache.len()
}

// GetIfPresent gets cached value from entries list and updates
// last access time for the entry if it is found.
func (c *localCache) GetIfPresent(k Key) (Value, bool) {
	en := c.cache.get(k, sum(k))
	if en == nil || en.getValue() == nil {
		c.stats.RecordMisses(1)
		return nil, false
	}

	en = c.get(en)

	now := currentTime()
	if en == nil {
		c.stats.RecordMisses(1)
		return nil, false
	}
	c.setEntryAccessTime(en, now)
	c.sendEvent(eventAccess, en)
	c.stats.RecordHits(1)
	return en.getValue(), true
}

// Put adds new entry to entries list.
func (c *localCache) Put(k Key, v Value) {
	seg := c.cache.segment(sum(k))
	en := seg.get(k)
	if en != nil && en.getValue() != nil {
		en = c.get(en)
	}
	now := currentTime()
	if en == nil || en.getValue() == nil {
		en, _ = seg.create(k, v, c)
		c.setEntryWriteTime(en, now)
		c.setEntryAccessTime(en, now)
		//if cen != nil {
		//	cen.setValue(v)
		//	c.setEntryWriteTime(cen, now)
		//	en = cen
		//}
	} else {
		// Update value and send notice
		en.setValue(v)
		c.setEntryWriteTime(en, now)
	}
	c.sendEvent(eventWrite, en)
}

// Invalidate removes the entry associated with key k.
func (c *localCache) Invalidate(k Key) {
	en := c.cache.get(k, sum(k))
	if en != nil {
		en.setInvalidated(true)
		c.sendEvent(eventDelete, en)
	}
}

// InvalidateAll resets entries list.
func (c *localCache) InvalidateAll() {
	c.cache.walk(func(en *entry) {
		en.setInvalidated(true)
	})
	c.sendEvent(eventDelete, nil)
}

func (s *localCache) CacheAllValue() []Value {
	var m []Value
	s.cache.walk(func(en *entry) {
		m = append(m, en.getValue())
	})
	return m
}

// Get returns value associated with k or call underlying loader to retrieve value
// if it is not in the cache. The returned value is only cached when loader returns
// nil error.
func (c *localCache) Get(k Key) (Value, error) {
	seg := c.cache.segment(sum(k))

	en := seg.get(k)
	if en != nil && en.getValue() != nil {
		en = c.get(en)
	}
	if en == nil || en.getValue() == nil {
		c.stats.RecordMisses(1)
		v, err, ev := c.load(k, seg)
		if err == nil {
			c.sendEvent(ev, v)
		}
		return v.getValue(), err
	}
	// Check if this entry needs to be refreshed
	c.sendEvent(eventAccess, en)
	return en.getValue(), nil
}

// get 获取数据
func (c *localCache) get(en *entry) *entry {
	now := currentTime()
	if c.isExpired(en, now) {
		send := false
		en.mutex.Lock()
		if !en.delete {
			c.cache.delete(en)
			send = true
		}
		en.mutex.Unlock()
		if send {
			c.sendEvent(eventDelete, en)
		}
		c.stats.RecordMisses(1)
		return nil
	} else {
		c.setEntryAccessTime(en, now)
		c.stats.RecordHits(1)
	}
	return en
}

// Refresh asynchronously reloads value for Key if it existed, otherwise
// it will synchronously load and block until it value is loaded.
func (c *localCache) Refresh(k Key) {
	//if c.loader == nil {
	//	return
	//}
	//en := c.cache.get(k, sum(k))
	//if en == nil {
	//	c.load(k)
	//} else {
	//	c.refreshAsync(en)
	//}
}

// Stats copies cache stats to t.
func (c *localCache) Stats(t *Stats) {
	c.stats.Snapshot(t)
}

func (c *localCache) processEntries() {
	//defer c.closeWG.Done()
	defer func() {
		if r := recover(); r != nil {
			logger.TraceErr(r)
			go c.processEntries()
		}
	}()
	for {
		select {
		case e := <-c.events:
			switch e.event {
			case eventWrite:
				c.write(e.entry)
				c.postWriteCleanup()
			case eventAccess:
				c.access(e.entry)
				c.postReadCleanup()
			case eventDelete:
				if e.entry == nil {
					c.removeAll()
				} else {
					c.remove(e.entry)
				}
				c.postReadCleanup()
			}
		case <-c.timer.C:
			c.postTimeCleanup()
		}
	}
	//for e := range c.events {
	//	switch e.event {
	//	case eventWrite:
	//		c.write(e.entry)
	//		c.postWriteCleanup()
	//	case eventAccess:
	//		c.access(e.entry)
	//		c.postReadCleanup()
	//	case eventDelete:
	//		if e.entry == nil {
	//			c.removeAll()
	//		} else {
	//			c.remove(e.entry)
	//		}
	//		c.postReadCleanup()
	//		//case eventClose:
	//		//	if c.reloader != nil {
	//		//		// Stop all refresh tasks.
	//		//		c.reloader.Close()
	//		//	}
	//		//	c.removeAll()
	//		//	return
	//	}
	//}
}

// sendEvent sends event only when the cache is not closing/closed.
func (c *localCache) sendEvent(typ event, en *entry) {
	if atomic.LoadInt32(&c.closing) == 0 {
		c.events <- entryEvent{en, typ}
	}
}

// This function must only be called from processEntries goroutine.
func (c *localCache) write(en *entry) {
	ren := c.accessQueue.write(en)
	c.writeQueue.write(en)
	if c.onInsertion != nil {
		c.onInsertion(en.key, en.getValue())
	}
	if ren != nil {
		c.writeQueue.remove(ren)
		// An entry has been evicted
		c.stats.RecordEviction()
		if c.onRemoval != nil {
			c.onRemoval(ren.key, ren.getValue())
		}
	}
}

// removeAll remove all entries in the cache.
// This function must only be called from processEntries goroutine.
func (c *localCache) removeAll() {
	c.accessQueue.iterate(func(en *entry) bool {
		c.remove(en)
		return true
	})
}

// remove removes the given element from the cache and entries list.
// It also calls onRemoval callback if it is set.
func (c *localCache) remove(en *entry) {
	ren := c.accessQueue.remove(en)
	c.writeQueue.remove(en)
	if ren != nil && c.onRemoval != nil {
		c.onRemoval(ren.key, ren.getValue())
	}
}

// access moves the given element to the top of the entries list.
// This function must only be called from processEntries goroutine.
func (c *localCache) access(en *entry) {
	c.accessQueue.access(en)
}

// create 加锁创建，避免重复创建
func (c *localCache) create(en *entry, cen *entry) (*entry, error, event) {
	en.mutex.Lock()
	defer en.mutex.Unlock()

	now := currentTime()
	if en.getValue() == nil {
		c.setEntryWriteTime(en, now)
		c.setEntryAccessTime(en, now)
		start := currentTime()
		v, err := c.loader(en.key)
		loadTime := now.Sub(start)
		if err != nil {
			c.stats.RecordLoadError(loadTime)
			return nil, err, eventNo
		}
		en.setValue(v)
		//if cen != nil {
		//	cen.setValue(v)
		//	c.setEntryWriteTime(cen, now)
		//	en = cen
		//}
		//c.sendEvent(eventWrite, en)
		c.stats.RecordLoadSuccess(loadTime)
		return en, nil, eventWrite
	} else {
		c.setEntryAccessTime(en, now)
		//c.sendEvent(eventAccess, en)
		c.stats.RecordHits(1)
		return en, nil, eventAccess
	}
}

// load uses current loader to synchronously retrieve value for k and adds new
// entry to the cache only if loader returns a nil error.
func (c *localCache) load(k Key, seg *segment) (*entry, error, event) {
	if c.loader == nil {
		panic("cache loader function must be set")
	}
	// TODO: Poll the value instead when the entry is loading.

	en, cen := seg.create(k, nil, c)

	return c.create(en, cen)
}

// refreshAsync reloads value in a go routine or using custom executor if defined.
func (c *localCache) refreshAsync(en *entry) bool {
	//if en.setLoading(true) {
	//	// Only do refresh if it isn't running.
	//	if c.reloader == nil {
	//		go c.refresh(en)
	//	} else {
	//		c.reload(en)
	//	}
	//	return true
	//}
	return false
}

// refresh reloads value for the given key. If loader returns an error,
// that error will be omitted. Otherwise, the entry value will be updated.
// This function would only be called by refreshAsync.
func (c *localCache) refresh(en *entry) {
	//defer en.setLoading(false)
	//if c.refreshAfterWrite > 0 {
	//	// TODO 不做刷新处理  设置刷新时间的重新读取数据
	//	//start := currentTime()
	//	//v, err := c.loader(en.key)
	//	//now := currentTime()
	//	//loadTime := now.Sub(start)
	//	//if err == nil {
	//	//	en.setValue(v)
	//	//	c.setEntryWriteTime(en, now)
	//	//	c.sendEvent(eventWrite, en)
	//	//	c.stats.RecordLoadSuccess(loadTime)
	//	//} else {
	//	//	// TODO: Log error
	//	//	c.stats.RecordLoadError(loadTime)
	//	//}
	//} else {
	//	now := currentTime()
	//	c.setEntryWriteTime(en, now)
	//	c.sendEvent(eventWrite, en)
	//}

}

// reload uses user-defined reloader to reloads value.
func (c *localCache) reload(en *entry) {
	//start := currentTime()
	//setFn := func(newValue Value, err error) {
	//	defer en.setLoading(false)
	//	now := currentTime()
	//	loadTime := now.Sub(start)
	//	if err == nil {
	//		en.setValue(newValue)
	//		c.setEntryWriteTime(en, now)
	//		c.sendEvent(eventWrite, en)
	//		c.stats.RecordLoadSuccess(loadTime)
	//	} else {
	//		c.stats.RecordLoadError(loadTime)
	//	}
	//}
	//c.reloader.Reload(en.key, en.getValue(), setFn)
}

// postReadCleanup is run after entry access/delete event.
// This function must only be called from processEntries goroutine.
func (c *localCache) postReadCleanup() {
	if atomic.AddInt32(&c.readCount, 1) > drainThreshold {
		atomic.StoreInt32(&c.readCount, 0)
		c.expireEntries()
	}
}

// postWriteCleanup is run after entry add event.
// This function must only be called from processEntries goroutine.
func (c *localCache) postWriteCleanup() {
	atomic.StoreInt32(&c.readCount, 0)
	c.expireEntries()
}

func (c *localCache) postTimeCleanup() {
	if c.expireAfterAccess <= 0 {
		return
	}
	now := currentTime()
	expiry := now.Add(-10 * time.Minute).UnixNano()
	if c.lastClearTime <= expiry {
		c.expireEntries()
	}
}

// expireEntries removes expired entries.
func (c *localCache) expireEntries() {
	remain := drainMax
	now := currentTime()
	if c.expireAfterAccess > 0 {
		c.lastClearTime = now.UnixNano()
		expiry := now.Add(-(c.expireAfterAccess + time.Minute)).UnixNano()
		c.accessQueue.iterate(func(en *entry) bool {
			if remain == 0 || en.getAccessTime() >= expiry {
				// Can stop as the entries are sorted by access time.
				// (the next entry is accessed more recently.)
				return false
			}
			// accessTime + expiry passed
			en.mutex.Lock()
			if en.getAccessTime() < expiry && !en.delete {
				c.cache.delete(en)
				c.remove(en)
				remain--
			}
			en.mutex.Unlock()
			c.stats.RecordEviction()
			return remain > 0
		})
	}
	//if remain > 0 && c.expireAfterWrite > 0 {
	//	expiry := now.Add(-c.expireAfterWrite).UnixNano()
	//	c.writeQueue.iterate(func(en *entry) bool {
	//		if remain == 0 || en.getWriteTime() >= expiry {
	//			return false
	//		}
	//		// writeTime + expiry passed
	//		c.remove(en)
	//		c.stats.RecordEviction()
	//		remain--
	//		return remain > 0
	//	})
	//}
	//if remain > 0 && c.loader != nil && c.refreshAfterWrite > 0 {
	//	expiry := now.Add(-c.refreshAfterWrite).UnixNano()
	//	c.writeQueue.iterate(func(en *entry) bool {
	//		if remain == 0 || en.getWriteTime() >= expiry {
	//			return false
	//		}
	//		// FIXME: This can cause deadlock if the custom executor runs refresh in current go routine.
	//		// The refresh function, when finish, will send to event channels.
	//		if c.refreshAsync(en) {
	//			// TODO: Maybe move this entry up?
	//			remain--
	//		}
	//		return remain > 0
	//	})
	//}
}

func (c *localCache) isExpired(en *entry, now time.Time) bool {
	if en.getInvalidated() {
		return true
	}
	if c.expireAfterAccess > 0 && en.getAccessTime() < now.Add(-c.expireAfterAccess).UnixNano() {
		// accessTime + expiry passed
		return true
	}
	if c.expireAfterWrite > 0 && en.getWriteTime() < now.Add(-c.expireAfterWrite).UnixNano() {
		// writeTime + expiry passed
		return true
	}
	return false
}

func (c *localCache) needRefresh(en *entry, now time.Time) bool {
	if en.getLoading() {
		return false
	}
	if c.refreshAfterWrite > 0 {
		tm := en.getWriteTime()
		if tm > 0 && tm < now.Add(-c.refreshAfterWrite).UnixNano() {
			// writeTime + refresh passed
			return true
		}
	}
	return false
}

// setEntryAccessTime sets access time if needed.
func (c *localCache) setEntryAccessTime(en *entry, now time.Time) {
	if c.expireAfterAccess > 0 {
		en.setAccessTime(now.UnixNano())
	}
}

// setEntryWriteTime sets write time if needed.
func (c *localCache) setEntryWriteTime(en *entry, now time.Time) {
	if c.expireAfterWrite > 0 || c.refreshAfterWrite > 0 {
		en.setWriteTime(now.UnixNano())
	}
}

// New returns a local in-memory Cache.
func New(options ...Option) Cache {
	c := newLocalCache()
	for _, opt := range options {
		opt(c)
	}
	c.init()
	return c
}

// NewLoadingCache returns a new LoadingCache with given loader function
// and cache options.
func NewLoadingCache(loader LoaderFunc, options ...Option) LoadingCache {
	c := newLocalCache()
	c.loader = loader
	for _, opt := range options {
		opt(c)
	}
	c.init()
	return c
}

// Option add options for default Cache.
type Option func(c *localCache)

// WithMaximumSize returns an Option which sets maximum size for the cache.
// Any non-positive numbers is considered as unlimited.
func WithMaximumSize(size int) Option {
	if size < 0 {
		size = 0
	}
	if size > maximumCapacity {
		size = maximumCapacity
	}
	return func(c *localCache) {
		c.cap = size
	}
}

// WithRemovalListener returns an Option to set cache to call onRemoval for each
// entry evicted from the cache.
func WithRemovalListener(onRemoval Func) Option {
	return func(c *localCache) {
		c.onRemoval = onRemoval
	}
}

// WithExpireAfterAccess returns an option to expire a cache entry after the
// given duration without being accessed.
func WithExpireAfterAccess(d time.Duration) Option {
	return func(c *localCache) {
		c.expireAfterAccess = d
	}
}

// WithExpireAfterWrite returns an option to expire a cache entry after the
// given duration from creation.
func WithExpireAfterWrite(d time.Duration) Option {
	return func(c *localCache) {
		c.expireAfterWrite = d
	}
}

// WithRefreshAfterWrite returns an option to refresh a cache entry after the
// given duration. This option is only applicable for LoadingCache.
//func WithRefreshAfterWrite(d time.Duration) Option {
//	return func(c *localCache) {
//		c.refreshAfterWrite = d
//	}
//}

// WithStatsCounter returns an option which overrides default cache stats counter.
func WithStatsCounter(st StatsCounter) Option {
	return func(c *localCache) {
		c.stats = st
	}
}

// WithPolicy returns an option which sets cache policy associated to the given name.
// Supported policies are: lru, slru, tinylfu.
func WithPolicy(name string) Option {
	return func(c *localCache) {
		c.policyName = name
	}
}

// WithReloader returns an option which sets reloader for a loading cache.
// By default, each asynchronous reload is run in a go routine.
// This option is only applicable for LoadingCache.
func WithReloader(reloader Reloader) Option {
	return func(c *localCache) {
		c.reloader = reloader
	}
}

// withInsertionListener is used for testing.
func withInsertionListener(onInsertion Func) Option {
	return func(c *localCache) {
		c.onInsertion = onInsertion
	}
}

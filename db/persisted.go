package db

import (
	"engine/util/concurrent"
	"sync"
	"time"
)

type TimingPersisted struct {
	elements concurrent.HashMap[string, *Element]
	name     string
	timer    *time.Ticker
	shut     chan bool
	wg       *sync.WaitGroup
	stop     bool
	m        *sync.RWMutex
	locks    concurrent.HashMap[string, *sync.Mutex]
}

func newTimingPersisted(cron time.Duration) *TimingPersisted {
	persisted := new(TimingPersisted)
	persisted.name = cron.String() + "db saver"
	persisted.shut = make(chan bool)
	persisted.timer = time.NewTicker(cron)
	persisted.m = new(sync.RWMutex)
	go persisted.run()
	return persisted
}

func (t *TimingPersisted) lockIdLock(id string) *sync.Mutex {
	lock, ok := t.locks.Get(id)
	if !ok {
		lock = &sync.Mutex{}
		v, vOk := t.locks.PutIfAbsent(id, lock)
		if vOk {
			lock = v
		}
	}
	lock.Lock()
	return lock
}
func (t *TimingPersisted) releaseIdLock(id string, lock *sync.Mutex) {
	lock.Unlock()
	t.locks.Remove(id)
}

func (t *TimingPersisted) shutDown() {
	t.stop = true
	logger.Warnf("关闭程序保存数据开始:%s", t.name)
	t.wg = &sync.WaitGroup{}
	t.wg.Add(1)
	t.shut <- true
	t.wg.Wait()
	logger.Warnf("关闭程序保存数据完成:%s", t.name)
}
func (t *TimingPersisted) put(element *Element) {
	if element == nil || t.stop {
		return
	}
	//读写锁
	t.m.RLock()
	defer t.m.RUnlock()

	//玩家各自表锁
	lock := t.lockIdLock(element.getIdentity())
	defer t.releaseIdLock(element.getIdentity(), lock)

	if element.event == update {
		s, ok := t.elements.Get(element.getIdentity())
		if ok {
			if s.event == save {
				element.event = save
			} else if s.event == remove {
				return
			}
		}
	}
	t.elements.Put(element.getIdentity(), element)
}

func (t *TimingPersisted) run() {
	defer func() {
		if r := recover(); r != nil {
			logger.TraceErr(r)
			go t.run() //出现退出异常，再次启动
		}
	}()
	for {
		select {
		//数据处理
		case <-t.timer.C:
			//定时器处理
			t.timerProcessing()
		case <-t.shut:
			ret := t.timerProcessing() //对数据进行保存
			logger.Warnf("关闭程序保存数据:%s,保存数据:%d条", t.name, ret)
			t.wg.Done()
			return
		}
	}
}
func (t *TimingPersisted) clearElements() map[string]*Element {
	t.m.Lock()
	ret := t.elements.Values()
	t.elements.Clear()
	t.m.Unlock()
	return ret
}

func (t *TimingPersisted) timerProcessing() int {
	defer func() {
		if r := recover(); r != nil {
			logger.TraceErr(r)
		}
	}()
	count := t.elements.Count()
	if count == 0 {
		return 0
	}
	el := t.clearElements()
	for _, v := range el {
		element := v
		switch element.event {
		case save:
			if element.ceObject != nil {
				element.ceObject.Before()
				accessor.Create(element.ceObject.GetEntity())
			} else {
				accessor.Create(element.dbObject)
			}

		case update:
			if element.ceObject != nil {
				element.ceObject.Before()
				accessor.Save(element.ceObject.GetEntity())
			} else {
				accessor.Save(element.dbObject)
			}

		case remove:
			accessor.Delete(element.dbObject.GetId(), element.dbObject)
		}
	}
	return len(el)
}

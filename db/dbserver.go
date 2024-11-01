package db

import (
	"time"
)

type dbServer struct {
	persisted map[time.Duration]*TimingPersisted
}

var globalDbServer = &dbServer{}

// InitDB 初始化数据库缓存信息
func (d *dbServer) InitDBCache() {
	d.persisted = make(map[time.Duration]*TimingPersisted)
	d.persisted[PRE5SECOND] = newTimingPersisted(PRE5SECOND)
	d.persisted[PRE30SECOND] = newTimingPersisted(PRE30SECOND)
	d.persisted[PRE1MINUTE] = newTimingPersisted(PRE1MINUTE)
	d.persisted[PRE5MINUTE] = newTimingPersisted(PRE5MINUTE)
}

// AddDBCache 添加数据到保存队列
func (d *dbServer) AddDBCache(eventType EventType, entity IEntity) {
	e := new(Element)
	e.dbObject = entity
	e.event = eventType
	d.addDirty(e)
}

// AddDBCacheElement 添加数据到保存队列
func (d *dbServer) AddDBCacheElement(eventType EventType, entity CElement) {
	e := new(Element)
	e.dbObject = entity.GetEntity()
	e.ceObject = entity
	e.event = eventType
	d.addDirty(e)
}

func (d *dbServer) addDirty(element *Element) {
	p, ok := d.persisted[element.dbObject.GetCron()]
	if !ok {
		logger.Errorf("Duration not have :%v", element.dbObject.GetCron())
		return
	}
	//if element.event == save || element.event == update{
	//	element.dbObject.Before()
	//}
	p.put(element)
}

// ShutDown 关闭程序保存数据
func (d *dbServer) ShutDown() {
	for _, v := range d.persisted {
		v.shutDown()
	}
}

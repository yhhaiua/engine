package db

import (
	"strconv"
)

type Element struct {
	dbObject IEntity
	event    EventType
	//cron		time.Duration
	ceObject CElement
}

func (e *Element) getIdentity() string {
	//return reflect.TypeOf(e.dbObject).Elem().String() + ":" + strconv.FormatInt(e.dbObject.GetId(), 10)
	return e.dbObject.TableName() + ":" + strconv.FormatInt(e.dbObject.GetId(), 10)
}

func (e *Element) update(s *Element) {
	e.dbObject = s.dbObject
	e.event = s.event
}

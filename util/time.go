package util

import "time"

//TimeMillis 获取系统的毫秒时间戳
func TimeMillis() int64 {
	return time.Now().UnixMilli()
}

//TimeSecond 获取系统的秒时间戳
func TimeSecond() int64 {
	return time.Now().Unix()
}

//GetFirstDateOfMonth 获取传入的时间所在月份的第一天，即某月第一天的0点
func GetFirstDateOfMonth(timeMillis int64) time.Time {
	d := time.UnixMilli(timeMillis)
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//GetFirstDateOfWeek 获取传入的时间所在周的周一，即某周周一的0点
func GetFirstDateOfWeek(timeMillis int64) time.Time {
	d := time.UnixMilli(timeMillis)
	offset := int(time.Monday - d.Weekday())
	if offset > 0 {
		offset = -6
	}
	d = GetZeroTime(d)
	d = d.AddDate(0, 0, offset)
	return d
}

//GetMorning 获取时间的零点
func GetMorning(timeMillis int64) int64 {
	t := time.UnixMilli(timeMillis)
	newTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return newTime.UnixMilli()
}

//GetZeroTime 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

//IsSameDay 是否同一天
func IsSameDay(tM1, tM2 int64) bool {
	t1 := time.UnixMilli(tM1)
	t2 := time.UnixMilli(tM2)
	if t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day() {
		return true
	}
	return false
}

//IsToDay 是否是今天
func IsToDay(millis int64) bool {
	t1 := time.Now()
	t2 := time.UnixMilli(millis)
	if t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day() {
		return true
	}
	return false
}

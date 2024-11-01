package util

import (
	"github.com/yhhaiua/engine/util/cast"
	"strconv"
	"strings"
)

const (
	Second int64 = 1000
	Minute       = 60 * Second
	Hour         = 60 * Minute
)

// ParseArray 一串以逗号分割的数字字符串，转换成int数组
func ParseArray(str string) []int {
	if len(str) == 0 {
		return nil
	}
	s := strings.Split(str, ",")
	if len(s) == 0 {
		return nil
	}
	var result []int
	for _, ss := range s {
		v, _ := strconv.Atoi(ss)
		result = append(result, v)
	}
	return result
}

func ParseArrayFloat(str string) []float64 {
	if len(str) == 0 {
		return nil
	}
	s := strings.Split(str, ",")
	if len(s) == 0 {
		return nil
	}
	var result []float64
	for _, ss := range s {
		v := cast.ToFloat64(ss)
		result = append(result, v)
	}
	return result
}

// ConvertToInt32 int切片转换成int32切片
func ConvertToInt32(source []int) []int32 {
	length := len(source)
	if length <= 0 {
		return nil
	}
	temp := make([]int32, length)
	for i := 0; i < length; i++ {
		temp[i] = int32(source[i])
	}
	return temp
}

// ConvertToInt int32切片转换成int切片
func ConvertToInt(source []int32) []int {
	length := len(source)
	if length <= 0 {
		return nil
	}
	temp := make([]int, length)
	for i := 0; i < length; i++ {
		temp[i] = int(source[i])
	}
	return temp
}

// Remove 移除切片中对应位置的元素，返回新切片
func Remove[T any](slice []T, removeIndex int) []T {
	l := len(slice)
	if l == 0 {
		return slice
	}
	if l <= removeIndex {
		return slice
	}
	newList := make([]T, len(slice)-1)
	for i := 0; i < l; i++ {
		if i == removeIndex {
			continue
		}
		if i < removeIndex {
			newList[i] = slice[i]
		} else {
			newList[i-1] = slice[i]
		}
	}
	return newList
}

// RandomIndexByRates 这个是权重切片，随机获取位置信息
func RandomIndexByRates(rates []int) int {
	l := len(rates)
	if l <= 0 {
		return 0
	}
	if l == 0 {
		return 0
	}

	allRate := 0
	newRates := make([]int, l)
	for i, v := range rates {
		allRate += v
		newRates[i] = allRate
	}
	num := Intn(allRate)
	for i, v := range newRates {
		if num < v {
			return i
		}
	}

	return -1
}

// RandomList 随机列表
func RandomList(all []int64) int64 {
	l := len(all)
	if l <= 0 {
		return 0
	}
	index := Intn(l)
	return all[index]
}

// RandomListInt 随机列表
func RandomListInt(all []int) int {
	l := len(all)
	if l <= 0 {
		return 0
	}
	index := Intn(l)
	return all[index]
}

type SliceCompare interface {
	~int | ~int64 | ~int32 | ~int8 | ~uint | ~float32 | ~float64 | ~string
}

// Contains 查找目标切片中是否有对应元素
func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func RemoveSlice[T SliceCompare](slice []T, value T) []T {
	var newList []T
	for _, v := range slice {
		if v == value {
			continue
		}
		newList = append(newList, v)
	}
	return newList
}
func RemoveSliceList[T SliceCompare](slice []T, value []T) []T {
	var newList []T
	for _, v := range slice {
		if Contains(value, v) {
			continue
		}
		newList = append(newList, v)
	}
	return newList
}

// RemoveRepeated 去除切片中重复的元素
func RemoveRepeated[T SliceCompare](slice []T) []T {
	newSlice := make([]T, 0)
	if len(slice) == 0 {
		return newSlice
	}
	m := make(map[T]T)
	for _, v := range slice {
		m[v] = v
	}

	for _, value := range m {
		newSlice = append(newSlice, value)
	}
	return newSlice
}

// HasSameValue 检测数组是否有重复数据
func HasSameValue[T SliceCompare](slice []T) bool {
	length := len(slice)
	if length == 0 {
		return false
	}
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			if slice[i] == slice[j] {
				return true
			}
		}
	}
	return false
}

// GetMillion 获取值对应的万分比
func GetMillion(value int) float64 {
	return float64(value) / float64(10000)
}

// ParseHourMinuteSecond 12:30:30 时分秒转换成时间戳
func ParseHourMinuteSecond(str string) int64 {
	ss := strings.Split(str, ":")
	var time int64 = 0
	if len(ss) >= 1 {
		time += cast.ToInt64(ss[0]) * Hour
	}
	if len(ss) >= 2 {
		time += cast.ToInt64(ss[1]) * Minute
	}
	if len(ss) >= 3 {
		time += cast.ToInt64(ss[2]) * Second
	}
	return time
}

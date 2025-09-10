package util

type intRange struct {
	from int
	to   int
}

// RangeMap 范围映射
// 范围映射是一个有序的映射，它将一个范围映射到一个值。
// 范围映射的键是一个闭区间，例如 [from, to]，值是一个任意类型的值。
// 范围映射的作用是将一个范围映射到一个值，例如将一个数字映射到一个字符串。
type RangeMap[T any] struct {
	ranges []intRange
	values []T
}

func (m *RangeMap[T]) Get(key int) (T, bool) {
	for i, rm := range m.ranges {
		if rm.from <= key && key <= rm.to {
			return m.values[i], true
		}
	}
	var zero T
	return zero, false
}
func (m *RangeMap[T]) GetOrDefault(key int, defaultValue T) T {
	v, ok := m.Get(key)
	if ok {
		return v
	}
	return defaultValue
}
func (m *RangeMap[T]) Put(from int, to int, value T) {
	m.ranges = append(m.ranges, intRange{from, to})
	m.values = append(m.values, value)
}

func (m *RangeMap[T]) Len() int {
	return len(m.ranges)
}

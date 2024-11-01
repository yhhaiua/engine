package concurrent

import (
	"github.com/yhhaiua/engine/util"
	"github.com/yhhaiua/engine/util/cast"
	"sync"
)

// CopyOnWriteArrayList 线程安全list
type CopyOnWriteArrayList[V comparable] struct {
	mutex sync.Mutex
	array []V
}

// Add 添加元素
func (c *CopyOnWriteArrayList[V]) Add(e V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	elements := c.getArray()
	l := len(elements)
	dst := make([]V, l+1)
	copy(dst, elements)
	dst[l] = e
	c.setArray(dst)
}

// AddAll 添加所有
func (c *CopyOnWriteArrayList[V]) AddAll(all []V) bool {
	csL := len(all)
	if csL == 0 {
		return false
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	elements := c.getArray()
	l := len(elements)
	if l == 0 {
		c.setArray(all)
	} else {
		dst := make([]V, l+csL)
		copy(dst, elements)
		copy(dst[l:], all)
		c.setArray(dst)
	}
	return true
}

// AddIndex 添加元素到某一位置，index小于0或者index大于Size，将添加失败panic
func (c *CopyOnWriteArrayList[V]) AddIndex(index int, e V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	elements := c.getArray()
	l := len(elements)
	if index > l || index < 0 {
		panic("IndexOutOfBoundsException---" + "Index: " + cast.ToString(index) + ", Size: " + cast.ToString(l))
	}
	newElements := make([]V, l+1)
	numMoved := l - index
	if numMoved == 0 {
		copy(newElements, elements)
	} else {
		copy(newElements, elements[0:index])
		copy(newElements[index+1:], elements[index:])
	}
	newElements[index] = e
	c.setArray(newElements)
}

// Set 更改list某一位置的值，返回旧的值和是否成功
func (c *CopyOnWriteArrayList[V]) Set(index int, e V) (V, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	elements := c.getArray()
	oldValue, ok := c.get(elements, index)
	if !ok {
		return oldValue, false
	}
	if oldValue != e {
		l := len(elements)
		dst := make([]V, l)
		copy(dst, elements)
		dst[index] = e
		c.setArray(dst)
	} else {
		c.setArray(elements)
	}
	return oldValue, true
}

// RemoveIndex 移除某一位置的值，返回旧值
func (c *CopyOnWriteArrayList[V]) RemoveIndex(index int) (V, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	elements := c.getArray()
	l := len(elements)
	oldValue, ok := c.get(elements, index)
	if !ok {
		return oldValue, false
	}
	numMoved := l - index - 1
	newElements := make([]V, l-1)
	if numMoved == 0 {
		copy(newElements, elements)
	} else {
		copy(newElements, elements[0:index])
		copy(newElements[index:], elements[index+1:])
	}
	c.setArray(newElements)
	return oldValue, true
}

// Remove 移除list中的某个元素
func (c *CopyOnWriteArrayList[V]) Remove(o V) bool {
	index := c.IndexOf(o)
	if index < 0 {
		return false
	}
	return c.remove(o, index)
}

func (c *CopyOnWriteArrayList[V]) remove(o V, index int) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	current := c.getArray()
	l := len(current)
	if index < l && current[index] == o {
		goto findIndex
	} else {
		prefix := util.Min(index, l)
		for i := 0; i < prefix; i++ {
			if o == current[i] {
				index = i
				goto findIndex
			}
		}
		if index >= l {
			return false
		}
		index = c.indexOf(o, current, index, l)
		if index < 0 {
			return false
		}
	}
findIndex:
	numMoved := l - index - 1
	newElements := make([]V, l-1)
	if numMoved == 0 {
		copy(newElements, current)
	} else {
		copy(newElements, current[0:index])
		copy(newElements[index:], current[index+1:])
	}
	c.setArray(newElements)
	return true
}

// AddIfAbsent 不重复添加，返回是否添加成功
func (c *CopyOnWriteArrayList[V]) AddIfAbsent(o V) bool {
	index := c.IndexOf(o)
	if index >= 0 {
		return false
	}
	return c.addIfAbsent(o)
}

// ContainsAll 检测list是否包含all中的所有元素
func (c *CopyOnWriteArrayList[V]) ContainsAll(all []V) bool {
	elements := c.getArray()
	l := len(elements)
	for _, e := range all {
		if c.indexOf(e, elements, 0, l) < 0 {
			return false
		}
	}
	return true
}

// RemoveAll 移除切片中的所有元素
func (c *CopyOnWriteArrayList[V]) RemoveAll(all []V) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	elements := c.getArray()
	l := len(elements)
	if l > 0 {
		temp := make([]V, 0, l)
		for _, v := range elements {
			if !util.Contains(all, v) {
				temp = append(temp, v)
			}
		}
		if len(temp) != l {
			c.setArray(temp)
			return true
		}
	}
	return false
}

func (c *CopyOnWriteArrayList[V]) addIfAbsent(o V) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	current := c.getArray()
	index := c.indexOf(o, current, 0, len(current))
	if index >= 0 {
		return false
	}
	l := len(current)
	dst := make([]V, l+1)
	copy(dst, current)
	dst[l] = o
	c.setArray(dst)
	return true
}

func (c *CopyOnWriteArrayList[V]) getArray() []V {
	return c.array
}
func (c *CopyOnWriteArrayList[V]) setArray(a []V) {
	c.array = a
}

// Size list元素数量
func (c *CopyOnWriteArrayList[V]) Size() int {
	return len(c.getArray())
}

// IsEmpty 判断列表是否为空
func (c *CopyOnWriteArrayList[V]) IsEmpty() bool {
	return c.Size() == 0
}

// Contains 判断元素是否存在列表中
func (c *CopyOnWriteArrayList[V]) Contains(o V) bool {
	elements := c.getArray()
	return c.indexOf(o, elements, 0, len(elements)) >= 0
}

// ToArray 获取列表中的所有元素
func (c *CopyOnWriteArrayList[V]) ToArray() []V {
	elements := c.getArray()
	dst := make([]V, len(elements))
	copy(dst, elements)
	return dst
}

// IndexOf 查找元素的位置
func (c *CopyOnWriteArrayList[V]) IndexOf(o V) int {
	elements := c.getArray()
	return c.indexOf(o, elements, 0, len(elements))
}

// LastIndexOf 从后查找元素的位置
func (c *CopyOnWriteArrayList[V]) LastIndexOf(o V) int {
	elements := c.getArray()
	return c.lastIndexOf(o, elements, len(elements)-1)
}
func (c *CopyOnWriteArrayList[V]) indexOf(o V, elements []V, index, fence int) int {
	for i := index; i < fence; i++ {
		if elements[i] == o {
			return i
		}
	}
	return -1
}
func (c *CopyOnWriteArrayList[V]) lastIndexOf(o V, elements []V, index int) int {
	for i := index; i >= 0; i-- {
		if elements[i] == o {
			return i
		}
	}
	return -1
}

func (c *CopyOnWriteArrayList[V]) get(a []V, index int) (V, bool) {
	if len(a) <= index {
		var e V
		return e, false
	}
	return a[index], true
}

// Get 通过位置获取对应元素
func (c *CopyOnWriteArrayList[V]) Get(index int) (V, bool) {
	return c.get(c.getArray(), index)
}

package concurrent

import (
	"testing"
)

type name struct {
	v int
}

func (n *name) GetValue() int {
	return n.v
}

type cc interface {
	GetValue() int
}

func TestHashMap(t *testing.T) {
	var m HashMap[int, cc]
	for i := 0; i < 100; i++ {
		m.Put(i, &name{v: i})
		m.PutIfAbsent(i, &name{v: i})
	}
	a := 1000
	//m.PutIfAbsent(a, &name{v: a})
	b, ok := m.Get(a)
	if ok {
		t.Log(b)
	}
	t.Log(b)
	c := m.Values()
	for k, v := range c {
		t.Log(k)
		t.Log(v)
	}
}

func TestList(t *testing.T) {
	var m CopyOnWriteArrayList[*name]
	m.AddIfAbsent(&name{v: 1})
	m.AddIfAbsent(&name{v: 2})
	m.AddIfAbsent(&name{v: 3})
	for _, v := range m.ToArray() {
		m.Remove(v)
	}
	t.Log(m.ToArray())
}

package linkedlist

import (
	"cmp"
	"fmt"
	"github.com/crrow/reona/util"
	"sync/atomic"
)

type Map[K cmp.Ordered, V any] struct {
	bSize  uint64
	size   atomic.Uint64
	mp     []*LinkedList[K, V]
	hasher func(K) uintptr
}

func NewMap[K cmp.Ordered, V any](opts ...util.Option[Map[K, V]]) *Map[K, V] {
	var r = new(Map[K, V])
	r.hasher = util.GetHasher[K]()
	util.ApplyOptions[Map[K, V]](r, opts...)
	return r
}

func WithCapacity[K cmp.Ordered, V any](nBucket uint64) util.Option[Map[K, V]] {
	return util.OptionFunc[Map[K, V]](func(t *Map[K, V]) {
		t.bSize = nBucket
		for i := 0; i < int(nBucket); i++ {
			t.mp = append(t.mp, New[K, V]())
		}
	})
}

func (m *Map[K, V]) Len() uint64 {
	return m.size.Load()
}

func (m *Map[K, V]) IsEmpty() bool {
	return m.size.Load() == 0
}

func (m *Map[K, V]) Insert(k K, v V) {
	ndx := uint64(m.hasher(k)) % m.bSize
	m.mp[ndx].Insert(k, v)
	m.size.Add(1)
	fmt.Printf("insert %v, index: %d \n", k, ndx)
}

func (m *Map[K, V]) Get(k K) (*V, bool) {
	ndx := uint64(m.hasher(k)) % m.bSize
	r := m.mp[ndx].Get(k)
	fmt.Printf("try get %v, index: %d \n", k, ndx)
	if r == nil {
		return nil, false
	}
	return r.Load(), true
}

func (m *Map[K, V]) Remove(k K) bool {
	ndx := uint64(m.hasher(k)) % m.bSize
	fmt.Printf("try remove %v, index: %d \n", k, ndx)
	if m.mp[ndx].Remove(k) {
		var cur = m.size.Load()
		for !m.size.CompareAndSwap(cur, cur-1) {
			cur = m.size.Load()
		}
		return true
	}
	return false
}

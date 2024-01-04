package linkedlist

import (
	"cmp"
	"hash/maphash"
	"sync/atomic"
	"unsafe"

	"github.com/crrow/reona/util"
)

type Map[K cmp.Ordered, V any] struct {
	bSize uint64
	size  atomic.Uint64
	mp    []*LinkedList[K, V]
	seed  maphash.Seed
}

func NewMap[K cmp.Ordered, V any](opts ...util.Option[Map[K, V]]) *Map[K, V] {
	var r = new(Map[K, V])
	r.seed = maphash.MakeSeed()
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
	ndx := m.calculateHash(k) % m.bSize
	m.mp[ndx].Insert(k, v)
	m.size.Add(1)
}

func (m *Map[K, V]) Get(k K) (*V, bool) {
	ndx := m.calculateHash(k) % m.bSize
	r := m.mp[ndx].Get(k)
	if r == nil {
		return nil, false
	}
	return r.Load(), true
}

func (m *Map[K, V]) Remove(k K) bool {
	ndx := m.calculateHash(k) % m.bSize
	return m.mp[ndx].Remove(k)
}

func (m *Map[K, V]) calculateHash(v any) uint64 {
	s := unsafe.Slice((*byte)(unsafe.Pointer(&v)), unsafe.Sizeof(v))
	return maphash.Bytes(m.seed, s)
}

package linkedlist

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockFreeLinkedList(t *testing.T) {
	l := New[int, int]()
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		var r *atomic.Pointer[int]
		for r == nil {
			r = l.Get(1)
		}
	}()
	go func() {
		wg.Add(1)
		r := l.Get(2)
		assert.Nil(t, r)
	}()
	go func() {
		wg.Add(1)
		l.Insert(1, 1)
	}()

	wg.Wait()
}

package linkedlist

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	mem := NewMap[string, int](WithCapacity[string, int](10))
	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		for {
			r, ok := mem.Get("hello")
			if !ok {
				continue
			}
			assert.Equal(t, 1, *r)
			return
		}
	}()

	r, ok := mem.Get("hello")
	assert.False(t, ok)
	assert.Nil(t, r)

	go func() {
		wg.Add(1)
		mem.Insert("hello", 1)
	}()

	go func() {
		wg.Add(1)
		mem.Insert("hello2", 2)
	}()

	go func() {
		wg.Add(1)
		for {
			r, ok := mem.Get("hello2")
			if !ok {
				continue
			}
			assert.Equal(t, 2, *r)
			return
		}
	}()

	wg.Wait()
}

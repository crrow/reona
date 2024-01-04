package linkedlist

import (
	"fmt"
	"hash/maphash"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	mem := NewMap[string, int](WithCapacity[string, int](10))
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		for {
			r, ok := mem.Get("hello")
			if !ok {
				fmt.Println("cannot get hello1")
				continue
			}
			assert.Equal(t, 1, *r)
			return
		}
	}()

	r, ok := mem.Get("hello")
	assert.False(t, ok)
	assert.Nil(t, r)
	r, ok = mem.Get("hello2")
	assert.False(t, ok)
	assert.Nil(t, r)

	go func() {
		defer wg.Done()
		mem.Insert("hello", 1)
	}()

	go func() {
		defer wg.Done()
		mem.Insert("hello2", 2)
	}()

	go func() {
		defer wg.Done()
		for {
			r, ok := mem.Get("hello2")
			if !ok {
				fmt.Println("cannot get hello2")
				continue
			}
			assert.Equal(t, 2, *r)
			return
		}
	}()

	wg.Wait()
}

func TestDeleteWhileRead(t *testing.T) {
	mem := NewMap[string, int](WithCapacity[string, int](10))
	mem.Insert("hello", 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		r, ok := mem.Get("hello")
		assert.True(t, ok)
		time.Sleep(1 * time.Second) // how to yield cpu time?
		assert.NotNil(t, r)
		assert.Equal(t, 1, *r) // we should be able to read the value

		// read again, we should fail
		_, ok = mem.Get("hello")
		assert.False(t, ok)
	}()
	time.Sleep(100 * time.Millisecond)
	ok := mem.Remove("hello")
	assert.True(t, ok)
	wg.Wait()
}

func TestHash(t *testing.T) {
	seed := maphash.MakeSeed()
	fmt.Println(calculateHash(seed, "hello"))
	fmt.Println(calculateHash(seed, "hello"))
	fmt.Println(calculateHash(seed, "hello1"))
	fmt.Println(calculateHash(seed, "hello1"))
	fmt.Println(calculateHash(seed, "hello2"))
	fmt.Println(calculateHash(seed, "hello2"))
}
func calculateHash(seed maphash.Seed, v any) uint64 {
	s := unsafe.Slice((*byte)(unsafe.Pointer(&v)), unsafe.Sizeof(v))
	return maphash.Bytes(seed, s)
}

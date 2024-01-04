# Lock-free data structure in go.

This repo is used for learning lock-free data structures in go.

At present, I implemented a lock-free linked-list based map.

TODO:

- [ ] benchmark
- [ ] skiplist

Be honest, it's indeed simpler to implement lock-free data structure without worrying about memory reclamation,
but go's atomic looks like a little wonky, all atomic is seq cst, no fetch add wrap method... 

And I'd rather operate raw pointer in rust rather than go, actually compiler diff. 

```go
package demo_test

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

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

```
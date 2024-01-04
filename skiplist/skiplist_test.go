//go:build exclude

package skiplist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSkipList(t *testing.T) {
	list := New[int, int]()
	assert.Equal(t, list.IsEmpty(), true)

	pos := list.searchPosition(1)
	fmt.Println(pos)
	//n, ok := list.Front()
	//assert.False(t, ok)
	//assert.Nil(t, n)
}

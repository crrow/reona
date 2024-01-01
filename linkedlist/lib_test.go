package linkedlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedList_Insert(t *testing.T) {
	l := New[int, int]()
	l.Insert(1, 1)
	ptr := l.Get(1)
	assert.NotNil(t, ptr)
	assert.Equal(t, 1, *ptr.Load())

	ptr = l.Get(2)
	assert.Nil(t, ptr)
}

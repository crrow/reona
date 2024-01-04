package linkedlist

import (
	"github.com/crrow/reona/linkedlist/lock"
	"github.com/crrow/reona/linkedlist/thread_unsafe"
	"math/rand"
	"testing"
)

func BenchmarkLockFree(b *testing.B) {
	b.Run("lockfree_insert", lockFreeInsert)
	b.Run("lockfeee_get", lockFreeGet)
	b.Run("thread_unsafe_insert", threadUnsafeInsert)
	b.Run("thread_unsafe_get", threadUnsafeGet)
	b.Run("lock_insert", lockInsert)
	b.Run("lock_get", lockGet)
}

func lockFreeInsert(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	l := New[int, int]()
	for i := 0; i < b.N; i++ {
		l.Insert(i, rand.Int())
	}
}
func lockFreeGet(b *testing.B) {
	l := New[int, int]()
	for i := 0; i < 10000; i++ {
		l.Insert(i, rand.Int())
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Get(i)
	}
}

func threadUnsafeInsert(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	l := thread_unsafe.New[int]()
	for i := 0; i < b.N; i++ {
		l.PushFront(i)
	}
}
func threadUnsafeGet(b *testing.B) {
	l := thread_unsafe.New[int]()
	for i := 0; i < 10000; i++ {
		l.PushFront(i)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.PopBack()
	}
}

func lockInsert(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	l := lock.NewLinkedList[int]()
	for i := 0; i < b.N; i++ {
		l.PushFront(i)
	}
}
func lockGet(b *testing.B) {
	l := lock.NewLinkedList[int]()
	for i := 0; i < 10000; i++ {
		l.Push(rand.Int())
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Pop()
	}
}

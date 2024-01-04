package linkedlist

import (
	"cmp"
	"sync/atomic"
)

type Node[K cmp.Ordered, V any] struct {
	key    K
	val    atomic.Pointer[V]
	next   atomic.Pointer[Node[K, V]]
	prev   atomic.Pointer[Node[K, V]]
	active atomic.Bool
}

func newNode[K cmp.Ordered, V any](k K, v V) *Node[K, V] {
	n := &Node[K, V]{
		key: k,
	}
	n.val.Store(&v)
	n.active.Store(true)
	return n
}

type LinkedList[K cmp.Ordered, V any] struct {
	head atomic.Pointer[Node[K, V]]
}

func New[K cmp.Ordered, V any]() *LinkedList[K, V] {
	return &LinkedList[K, V]{}
}

func (l *LinkedList[K, V]) Insert(k K, v V) {
	var curAtomicPtrToNode = &l.head // &AtomicPtr
	for {
		var curNode = curAtomicPtrToNode.Load() // *Node
		if curNode == nil {                     // the current node is null, store point here
			newNodePtr := newNode(k, v)
			if !l.head.CompareAndSwap(curNode, newNodePtr) {
				continue // TODO: will we actually load it again ?
			}
			return
		}
		// find the same key
		if curNode.key == k && curNode.active.Load() {
			originalNodeValPtr := curNode.val.Load()
			if !curNode.val.CompareAndSwap(originalNodeValPtr, &v) {
				// the value has been changed, go back again
				continue
			}
			// cas succeed
			_ = originalNodeValPtr
			return
		}
		curAtomicPtrToNode = &curNode.next // point to next
		// key does not exist
		nextNodePtr := curNode.next.Load()
		if nextNodePtr == nil {
			insNodePtr := newNode(k, v)
			insNodePtr.prev.Store(curNode)               // we won't fail here
			curNode.next.CompareAndSwap(nil, insNodePtr) // we may fail here
			return
		}
	}
}

func (l *LinkedList[K, V]) Get(k K) *atomic.Pointer[V] {
	var curAtomicPtrToNode = &l.head
	for {
		curNode := curAtomicPtrToNode.Load()
		if curNode == nil {
			return nil
		}
		if curNode.key == k && curNode.active.Load() {
			return &curNode.val
		}
		curAtomicPtrToNode = &curNode.next
	}
}

func (l *LinkedList[K, V]) Remove(k K) bool {
	var curAtomicPtrToNode = &l.head
	for {
		curNode := curAtomicPtrToNode.Load()
		if curNode == nil {
			return false
		}

		if curNode.key == k && curNode.active.Load() {
			if curNode.active.CompareAndSwap(true, false) {
				return false
			}

			next := curNode.next.Load()
			prev := curNode.prev.Load()
			if next != nil {
				if prev != nil {
					if !prev.next.CompareAndSwap(curNode, next) {
						// someone may delete it already
						return false
					}
					if !next.prev.CompareAndSwap(curNode, next) {
						return false
					}
				} else { // prev == nil
					if !next.prev.CompareAndSwap(curNode, nil) {
						return false
					}
					if !l.head.CompareAndSwap(curNode, next) {
						return false
					}
				}
			} else {
				if prev != nil {
					if !prev.next.CompareAndSwap(curNode, nil) {
						return false
					}
				} else {
					if !l.head.CompareAndSwap(curNode, nil) {
						return false
					}
				}
			}

			_ = curNode
			return true
		}
		curAtomicPtrToNode = &curNode.next
	}
}

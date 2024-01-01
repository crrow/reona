//go:build exclude
package skiplist

import (
	"cmp"
	"sync/atomic"
)

// heightBits Number of bits needed to store height.
const heightBits uint64 = 5

// maxHeight height of a skip list tower.
const maxHeight uint64 = 1 << heightBits

// heightMask The bits of `refs_and_height` that keep the height.
const heightMask uint64 = (1 << heightBits) - 1

// Tower The tower of atomic pointers.
// The actual size of the tower will vary depending on the height that a node
// was allocated with.
type Tower[K cmp.Ordered, V any] struct {
	// disgusting
	pointers []atomic.Pointer[Node[K, V]]
}

// Node A skip list node.
//
// TODO: Will go reorder the struct fields?
// TODO: should K always be comparable ?
type Node[K cmp.Ordered, V any] struct {
	// the value
	value V
	// the key
	key K
	// Keeps the reference count and the height of its tower.
	//
	// The reference count is equal to the number of Entry pointing to this node, plus the
	// number of levels in which this node is installed.
	refsAndHeight atomic.Uint64
	// Whether this node is marked as deleted.
	markedAsDeleted atomic.Bool // it is behind an atomic, so we just use normal bool
	// The tower of atomic pointers.
	tower Tower[K, V]
}

// Height Returns the height of this node's tower.
func (n *Node[K, V]) Height() uint64 {
	return n.refsAndHeight.Load()&heightMask + 1
}

type SkipList[K cmp.Ordered, V any] struct {
	// The head of the skip list (just a dummy node, not a real entry).
	head Tower[K, V]
	// The seed for random height generation.
	seed atomic.Uint64
	// The number of entries in the skip list.
	len atomic.Uint64
	// the highest tower currently in use.
	// This value is used as a hint for where
	// to start lookups and never decreases.
	maxHeight atomic.Uint64
}

func New[K cmp.Ordered, V any]() *SkipList[K, V] {
	sl := &SkipList[K, V]{
		head: Tower[K, V]{
			pointers: make([]atomic.Pointer[Node[K, V]], maxHeight),
		},
	}

	sl.seed.Store(1)
	sl.maxHeight.Store(1)

	return sl
}

// Returns the number of entries in the skip list.
//
// If the skip list is being concurrently modified, consider the returned number just an
// approximation without any guarantees.
func (sl *SkipList[K, V]) Len() uint64 {
	return sl.len.Load()
}

// Returns `true` if the skip list is empty.
func (sl *SkipList[K, V]) IsEmpty() bool {
	return sl.Len() == 0
}

// Front returns the entry with the smallest key.
func (sl *SkipList[K, V]) Front() (*Node[K, V], bool) {
	return sl.nextNode(&sl.head, reona.newUnbound[K]())
}

// Returns the successor of a node.
//
// This will keep searching until a non-deleted node is found. If a deleted
// node is reached then a search is performed using the given key.
func (sl *SkipList[K, V]) nextNode(pred *Tower[K, V], lowerBounder reona.bounder[K]) (*Node[K, V], bool) {
	// Load the level 0 successor of the current node.
	curr := pred.pointers[0].Load()
	// FIXME: may pred be deleted? If so, what should we do?
	if curr == nil || curr.markedAsDeleted {
		return sl.searchBound(lowerBounder, false)
	}
	// current node is not null
	for curr != nil {
		// Loads its level 0 successor.
		succ := curr.tower.pointers[0].Load()
		// If the successor's successor is marked,
		if succ.markedAsDeleted { // the successor has been deleted.
			// attempts to help with removal using help_unlink.
			if newNext, ok := sl.helpUnlink(&pred.pointers[0], curr, succ); ok {
				// If help_unlink succeeds, continues searching from the updated successor.
				// On success, continue searching through the current level.
				curr = newNext
				continue
			} else {
				// On failure, we cannot do anything reasonable to continue
				// searching from the current position. Restart the search.
				return sl.searchBound(lowerBounder, false)
			}
		}
		// If a non-marked successor is found, returns it as the valid successor.
		return curr, true
	}
	// If no valid successor is found (end of skiplist), returns None.
	return nil, false
}

// Searches for first/last node that is greater/less/equal to a key in the skip list.
//
// If `upper_bound == true`: the last node less than (or equal to) the key.
//
// If `upper_bound == false`: the first node greater than (or equal to) the key.
func (sl *SkipList[K, V]) searchBound(bounder reona.bounder[K], upperBound bool) (*Node[K, V], bool) {
	panic("implement me")
}

// If we encounter a deleted node while searching, help with the deletion
// by attempting to unlink the node from the list.
//
// If the unlinking is succeeded, then this function returns the next node
// with which the search should continue at the current level.
func (sl *SkipList[K, V]) helpUnlink(pred *atomic.Pointer[Node[K, V]], curr, next *Node[K, V]) (*Node[K, V], bool) {
	panic("implement me")
}

// Inserts an entry with the specified `key` and `value`.
// If `replace` is `true`, then any existing entry with this key will first be removed.
func (sl *SkipList[K, V]) doInsert(key K, val V, replace bool) *Node[K, V] {
	var searchResult *_Position[K, V]
	for {
		searchResult = sl.searchPosition(key)
		if searchResult.found == nil {
			break
		}
		if replace {
			// FIXME: cannot update there
			searchResult.found..
			markedAsDeleted = true
		}
	}
}

// Searches for a key in the skip list and returns a list of all adjacent nodes.
func (sl *SkipList[K, V]) searchPosition(key K) *_Position[K, V] {
search:
	for {
		result := _Position[K, V]{}
		// The current level we're at.
		level := sl.maxHeight.Load()
		// Fast loop to skip empty tower levels.
		for ; level >= 1 && sl.head.pointers[level-1].Load() == nil; level-- {
		}

		pred := sl.head

		for ; level >= 1; level-- {
			// Two adjacent nodes at the current level.
			currPtr := &pred.pointers[level]
			curr := currPtr.Load()
			if curr != nil && curr.markedAsDeleted {
				// If `curr` is marked, that means `pred` is removed and we have to restart the
				// search.
				continue search
			}

			// Iterate through the current level until we reach a node with a key greater
			// than or equal to `key`.
		walk:
			for curr != nil {
				succ := curr.tower.pointers[level].Load()
				if succ == nil || succ.markedAsDeleted {
					if next, ok := sl.helpUnlink(&pred.pointers[level], curr, succ); ok {
						curr = next
						continue
					} else {
						// On failure, we cannot do anything reasonable to continue
						// searching from the current position. Restart the search.
						continue search
					}
				}

				switch cmp.Compare(curr.key, key) {
				case 1:
					break walk
				case 0:
					result.found = currPtr
					break walk
				case -1:
					// do nothing
				}

				// Move one step forward.
				pred = curr.tower
				curr = succ
			}
			result.left[level] = &pred
			result.right[level] = currPtr
		}
		return &result
	}
}

// A search result.
//
// The result indicates whether the key was found, as well as what were the adjacent nodes to the
// key on each level of the skip list.
type _Position[K cmp.Ordered, V any] struct {
	// reference a node with the given key, if found.
	// If this is not nil then it will point to the same node as `right[0]`.
	found *atomic.Pointer[Node[K, V]]
	// Adjacent nodes with smaller keys (predecessors).
	left [maxHeight]*Tower[K, V]
	// Adjacent nodes with equal or greater keys (successors).
	right [maxHeight]*atomic.Pointer[Node[K, V]]
}

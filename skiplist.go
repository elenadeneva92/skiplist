package skiplist

import (
	"math/rand"
	"time"
)

type Value interface{}

var defaultValue = struct{}{}

func NewSkipList(height int, prob float64) SkipList {
	return &genericSkipList{
		height: height,
		head: &node{
			next: make([]*node, height),
		},
	}
}

type Key interface {
	LessEq(Key) bool
	Eq(Key) bool
	Less(Key) bool
}

type node struct {
	key   Key
	value Value
	next  []*node
}

func (n *node) level() int {
	return len(n.next)
}

type SkipList interface {
	Add(Key, Value)
	Find(Key) (Value, bool)
	Delete(Key) (Value, bool)
	Len() int
	Iterator() Iterator
}

type genericSkipList struct {
	height int
	head   *node
	size   int
	prob   float64
}

type Iterator interface {
	Next() bool
	Key() Key
	Value() Value
}

type skipListIterator struct {
	current *node
}

func (it *skipListIterator) Next() bool {
	if it.current == nil || it.current.next[0] == nil {
		return false
	}

	it.current = it.current.next[0]
	return true
}

func (it skipListIterator) Key() Key {
	return it.current.key
}

func (it skipListIterator) Value() Value {
	return it.current.value
}

// Find searched the value under the given key
// the second returned value show if an element with the given key is found
func (sl *genericSkipList) Find(key Key) (Value, bool) {
	var current = sl.head

	for level := sl.height - 1; level >= 0; level-- {
		for next := current.next[level]; next != nil && next.key.LessEq(key); {
			current = next
			next = current.next[level]
		}
	}

	if current == sl.head || !current.key.Eq(key) {
		return defaultValue, false
	}

	return current.value, true
}

// Push stores the value under the given key in the list
// if the key already exists, the value will be updated
func (sl *genericSkipList) Add(key Key, value Value) {
	var (
		current = sl.head
		update  = make([]*node, sl.height)
	)

	for level := sl.height - 1; level >= 0; level-- {
		for next := current.next[level]; next != nil && next.key.LessEq(key); {
			current = next
			next = current.next[level]
		}
		update[level] = current
	}

	if current != sl.head && key.Eq(current.key) {
		current.next[0].value = value
		return
	}

	var newNode = node{
		key:   key,
		value: value,
		next:  make([]*node, sl.randLevel()),
	}

	for level := newNode.level() - 1; level >= 0; level-- {
		newNode.next[level] = update[level].next[level]
		update[level].next[level] = &newNode
	}

	sl.size++
}

// Len return the size of elements in the skip list
func (sl *genericSkipList) Len() int {
	return sl.size
}

// Delete removes the element with the given key
// returns the deleted value and flag if the deletion was successful
func (sl *genericSkipList) Delete(key Key) (Value, bool) {
	var (
		value   Value = defaultValue
		ok            = false
		current       = sl.head
	)

	for level := sl.height - 1; level >= 0; level-- {
		for next := current.next[level]; next != nil && next.key.Less(key); {
			current = next
			next = current.next[level]
		}

		if current.next[level] != nil && current.next[level].key.Eq(key) {
			value = current.next[level].value
			ok = true
			current.next[level] = current.next[level].next[level]
		}
	}

	if ok {
		sl.size--
	}
	return value, ok
}

func (sl *genericSkipList) Iterator() Iterator {
	return &skipListIterator{
		current: sl.head.next[0],
	}
}

func (sl *genericSkipList) randLevel() int {
	var level = 1
	rand.Seed(time.Now().UTC().UnixNano()) // be sure it's really "random"
	for ; level < sl.height-1 && rand.Float64() < sl.prob; level++ {
	}
	return level
}

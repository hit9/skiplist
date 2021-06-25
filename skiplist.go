// Copyright 2016 Chao Wang <hit9@icloud.com>.

/*

Package skiplist implements in-memory skiplist.

Reference: https://en.wikipedia.org/wiki/Skip_list

Example

	sl := New(7)
	sl.Put(Int(3))
	sl.Put(Int(9))
	sl.Put(Int(2))
	...

And the result will be:

	Level[0]: 1 -> 2 -> 3 -> 4 -> 5 -> 6 -> 7 -> 8 -> 9 -> nil
	Level[1]: 1 -> 2 -> 4 -> 5 -> 7 -> 8 -> 9 -> nil
	Level[2]: 2 -> 4 -> 5 -> 7 -> 9 -> nil
	Level[3]: 4 -> 5 -> 7 -> 9 -> nil
	Level[4]: 4 -> 7 -> nil

Another example:

	type Item struct {
		score int
		value string
	}

	func (item Item) Less(than skiplist.Item) bool {
		return item.score < than.(Item).score
	}

	func main() {
		sl := skiplist.New(11)
		sl.Put(Item{3, "data1"})
		sl.Put(Item{5, "data2"})
		item := sl.Get(Item{score: 3})
		item.(Item).value // "data1"
	}

Iterator example:

	iter := sl.NewIterator(nil)
	for iter.Next() {
		item := iter.Item()
		...
	}

Complexity

Operation Put/Get/Delete time complexity are all O(logN). And the space
complexity is O(NlogN).

Goroutine Safety

No. Lock granularity depends on the use case.

*/
package skiplist // import "github.com/hit9/skiplist"

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

// Item is a single object in the skiplist.
type Item interface {
	// Less tests whether the item is less than given argument.
	// Must provide strict less,  we treat !a.Less(b) && !b.Less
	// to mean a == b.
	Less(than Item) bool
}

// equal tests whether the item equal the given argument.
func equal(item, than Item) bool {
	return !item.Less(than) && !than.Less(item)
}

// Int implements the Item interface for integers.
type Int int

// Less returns true if int(a) < int(b)
func (i Int) Less(j Item) bool {
	return i < j.(Int)
}

// node is an internel node in the skiplist.
type node struct {
	item     Item
	forwards []*node
}

// SkipList is an implementation of skiplist.
type SkipList struct {
	length   int
	level    int
	maxLevel int
	head     *node
	rand     *rand.Rand
	buf      []*node
}

// Iterator is skiplist iterator.
type Iterator struct {
	sl *SkipList
	n  *node
}

// FactorP is the propability to get the rand level.
var FactorP = 0.5

func newNode(level int, item Item) *node {
	return &node{
		item:     item,
		forwards: make([]*node, level, level),
	}
}

// New creates a new SkipList.
func New(maxLevel int) *SkipList {
	return NewWithRandSeed(maxLevel, time.Now().UnixNano())
}

// NewWithRandSeed creates a new SkipList with a given seed.
func NewWithRandSeed(maxLevel int, seed int64) *SkipList {
	if maxLevel < 2 {
		panic("skiplist: bad maxLevel")
	}
	return &SkipList{
		maxLevel: maxLevel,
		head:     newNode(maxLevel, nil),
		rand:     rand.New(rand.NewSource(seed)),
		buf:      make([]*node, maxLevel, maxLevel),
	}
}

// Len returns skiplist length.
func (sl *SkipList) Len() int { return sl.length }

// Level returns skiplist level.
func (sl *SkipList) Level() int { return sl.level }

// MaxLevel returns skiplist maxLevel.
func (sl *SkipList) MaxLevel() int { return sl.maxLevel }

// randLevel returns a level between 1 and maxLevel.
func (sl *SkipList) randLevel() int {
	level := 1
	for sl.rand.Int()&0xffff < int(FactorP*float64(0xffff)) {
		level++
	}
	if level < sl.maxLevel {
		return level
	}
	return sl.maxLevel
}

func (sl *SkipList) resetBuf() {
	for i := 0; i < sl.maxLevel; i++ {
		sl.buf[i] = nil
	}
}

// Put adds an item to the skiplist. O(logN)
func (sl *SkipList) Put(item Item) {
	// Reuse update array and find the node.
	sl.resetBuf()
	update := sl.buf
	n := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for n.forwards[i] != nil && n.forwards[i].item.Less(item) {
			n = n.forwards[i]
		}
		update[i] = n
	}
	// New level.
	level := sl.randLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.head
		}
		sl.level = level
	}
	// Add node.
	n = newNode(level, item)
	for i := 0; i < level; i++ {
		n.forwards[i] = update[i].forwards[i]
		update[i].forwards[i] = n
	}
	sl.length++
}

// Get an item from the skiplist, nil on not found. O(logN)
func (sl *SkipList) Get(item Item) Item {
	n := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for n.forwards[i] != nil && n.forwards[i].item.Less(item) {
			n = n.forwards[i]
		}
	}
	n = n.forwards[0]
	if n != nil && equal(n.item, item) {
		return n.item
	}
	return nil
}

// Has tests whether skiplist contains an item. O(logN)
func (sl *SkipList) Has(item Item) bool { return sl.Get(item) != nil }

// Delete an item from skiplist and return it, nil on not found. O(logN)
func (sl *SkipList) Delete(item Item) Item {
	// Find node.
	sl.resetBuf()
	update := sl.buf
	head := sl.head
	n := head
	for i := sl.level - 1; i >= 0; i-- {
		for n.forwards[i] != nil && n.forwards[i].item.Less(item) {
			n = n.forwards[i]
		}
		update[i] = n
	}
	n = n.forwards[0]
	if n == nil || !equal(n.item, item) {
		return nil
	}
	// Delete
	for i := 0; i < sl.level; i++ {
		if update[i].forwards[i] == n {
			update[i].forwards[i] = n.forwards[i]
		}
	}
	// Decrease level if need.
	for sl.level > 1 && head.forwards[sl.level-1] == nil {
		sl.level--
	}
	sl.length--
	return n.item
}

// First returns the first item, nil on not found. O(1)
func (sl *SkipList) First() Item {
	if sl.length == 0 {
		return nil
	}
	return sl.head.forwards[0].item
}

// PopFirst pops the first item and returns it, nil on empty. O(1)
func (sl *SkipList) PopFirst() Item {
	if sl.length == 0 {
		return nil
	}
	n := sl.head.forwards[0]
	for i := sl.level - 1; i >= 0; i-- { // Release upward
		if sl.head.forwards[i] == n {
			sl.head.forwards[i] = n.forwards[i]
		}
	}
	for sl.level > 1 && sl.head.forwards[sl.level-1] == nil {
		sl.level--
	}
	sl.length--
	return n.item
}

// Clear the skiplist.
func (sl *SkipList) Clear() {
	for sl.PopFirst() != nil {
	}
}

// NewIterator returns a new iterator on this skiplist with an item start,
// if the start is nil, iterator starts on head.
// Filter items >= start.
func (sl *SkipList) NewIterator(start Item) *Iterator {
	n := sl.head
	if start != nil {
		for i := sl.level - 1; i >= 0; i-- {
			for n.forwards[i] != nil && n.forwards[i].item.Less(start) {
				n = n.forwards[i]
			}
		}
	}
	return &Iterator{sl: sl, n: n}
}

// Next seeks iterator next, returns false on end.
func (iter *Iterator) Next() bool {
	iter.n = iter.n.forwards[0]
	return iter.n != nil
}

// Item returns current item on the iterator.
func (iter *Iterator) Item() Item {
	return iter.n.item
}

// Print the skiplist, debug purpose.
func (sl *SkipList) Print(w io.Writer) {
	for i := 0; i < sl.level; i++ {
		n := sl.head.forwards[i]
		fmt.Fprintf(w, "Level[%d]: ", i)
		for n != nil {
			fmt.Fprintf(w, "%v -> ", n.item)
			n = n.forwards[i]
		}
		fmt.Fprintf(w, "nil\n")
	}
}

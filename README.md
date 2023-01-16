SkipList
========

Package skiplist implements in-memory skiplist.

https://pkg.go.dev/github.com/hit9/skiplist

关于跳跃表的一个中文讲解 https://writings.sh/post/data-structure-skiplist

Example
-------

```go
package main

import (
	"fmt"
	"github.com/hit9/skiplist"
	"os"
)

// Item implements skiplist.Item
type Item struct {
	score int
	value string
}

// Less returns true if item is less than another.
func (item Item) Less(than skiplist.Item) bool {
	return item.score < than.(Item).score
}

func main() {
	// New skiplist with max level 16
	sl := skiplist.New(16)
	// Put some items.
	sl.Put(Item{3, "data1"})
	sl.Put(Item{1, "data2"})
	sl.Put(Item{5, "data3"})
	sl.Put(Item{2, "data4"})
	sl.Put(Item{6, "data5"})
	// Get one item.
	item := sl.Get(Item{score: 2})
	fmt.Println(item.(Item).value) // "data4"
	// Print the skiplist.
	sl.Print(os.Stdout)
}
```

Output:

```
Level[0]: {1 data2} -> {2 data4} -> {3 data1} -> {5 data3} -> {6 data5} -> nil
Level[1]: {1 data2} -> {2 data4} -> {5 data3} -> {6 data5} -> nil
Level[2]: {1 data2} -> {5 data3} -> {6 data5} -> nil
Level[3]: {1 data2} -> {6 data5} -> nil
Level[4]: {6 data5} -> nil
```

Complexity
----------

Operation Put/Get/Delete time complexity are all O(logN). And the space
complexity is O(NlogN).

License
-------

BSD.

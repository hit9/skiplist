// Copyright 2016 Chao Wang <hit9@icloud.com>.

// +build ignore

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

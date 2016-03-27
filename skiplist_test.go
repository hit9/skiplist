// Copyright 2016 Chao Wang <hit9@icloud.com>.

package skiplist

import (
	"math/rand"
	"runtime"
	"testing"
)

// Must asserts the given value is True for testing.
func Must(t *testing.T, v bool) {
	if !v {
		_, fileName, line, _ := runtime.Caller(1)
		t.Errorf("\n unexcepted: %s:%d", fileName, line)
	}
}

func TestPut(t *testing.T) {
	sl := New(16)
	n := 1024 * 10
	for i := 0; i < n; i++ {
		item := Int(rand.Int())
		sl.Put(item)
		// Must get
		Must(t, equal(sl.Get(item), item))
		// Must len++
		Must(t, sl.Len() == i+1)
	}
}

func TestGet(t *testing.T) {
	sl := New(16)
	n := 1024 * 10
	for i := 0; i < n; i++ {
		item := Int(rand.Int() % n)
		sl.Put(item)
		// Must get
		Must(t, equal(sl.Get(item), item))
		// Must cant get
		Must(t, sl.Get(Int(n+rand.Int())) == nil)
	}
}

func TestDelete(t *testing.T) {
	sl := New(16)
	n := 1024 * 10
	for i := 0; i < n; i++ {
		item := Int(rand.Int() % n)
		sl.Put(item)
		Must(t, sl.Len() == 1)
		// Must delete
		Must(t, sl.Delete(item) == item)
		// Must cant delete
		Must(t, sl.Delete(Int(n+rand.Int())) == nil)
		Must(t, sl.Len() == 0)
	}
}

func TestIteratorNil(t *testing.T) {
	sl := New(7)
	n := 1024
	for i := n - 1; i >= 0; i-- {
		sl.Put(Int(i))
	}
	iter := sl.NewIterator(nil)
	i := 0
	for iter.Next() {
		// Must equal
		Must(t, Int(i) == iter.Item())
		i++
	}
}

func TestIteratorStart(t *testing.T) {
	sl := New(7)
	n := 1024
	for i := n - 1; i >= 0; i-- {
		sl.Put(Int(i))
	}
	start := rand.Intn(n)
	iter := sl.NewIterator(Int(start))
	i := 0
	for iter.Next() {
		// Must equal
		Must(t, Int(i+start) == iter.Item())
		i++
	}
	Must(t, i == n-start)
}

// The maxLevel masters the bench results.
func BenchmarkPut(b *testing.B) {
	sl := New(50)
	for i := 0; i < b.N; i++ {
		sl.Put(Int(i))
	}
}

func BenchmarkGet(b *testing.B) {
	sl := New(50)
	for i := 0; i < b.N; i++ {
		sl.Put(Int(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Get(Int(i))
	}
}

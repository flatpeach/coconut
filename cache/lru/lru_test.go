// Copyright 2014 The coconut Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lru

import (
	"testing"
)

type cacheItem struct {
	v []byte
}

func (i *cacheItem) Size() uint64 {
	return uint64(len(i.v))
}

func TestLRUBasic(t *testing.T) {
	c := New(&Option{1 << 20, 0})

	if c.Capacity() != 1<<20 {
		t.Fatal("The capacity of LRU cache not matched!")
	}

	key := "hello"

	value := &cacheItem{[]byte("HelloWorld")}

	c.Set(key, value)

	v, ok := c.Get(key)

	if !ok {
		t.Fatal("Didn't get the key in memory")
		return
	}

	if v.(*cacheItem) != value {
		t.Fatal("Data mismatched!")
		return
	}

	if c.Size() != uint64(value.Size()) {
		t.Fatal("size not matched")
	}

	c.SetCapacity(1 << 30)

	if c.Capacity() != 1<<30 {
		t.Fatal("Set capacity failed")
	}

	c.Delete(key)

	if c.Size() != 0 {
		t.Fatal("Failed to delete one item")
	}

	v, ok = c.Get(key)
	if ok {
		t.Fatal("Failed to delete one element, still can access")
	}

}

func TestLRUEvict(t *testing.T) {
	c := New(&Option{0, 2})

	key1 := string("k1")
	key2 := string("k2")
	key3 := string("k3")

	val1 := &cacheItem{[]byte("HelloWorld")}
	val2 := &cacheItem{[]byte("HelloWorld")}
	val3 := &cacheItem{[]byte("HelloWorld")}

	c.Set(key1, val1)
	if c.ElementsCount() != 1 {
		t.Fatal("Count of elements should be equal to 1")
	}

	c.Set(key2, val2)
	if c.ElementsCount() != 2 {
		t.Fatal("Count of elements should be equal to 2")
	}

	c.Set(key3, val3)
	if c.ElementsCount() != 2 {
		t.Fatal("Count of elements should be equal to 2")
	}

	c.Evict(1)
	if c.ElementsCount() != 1 {
		t.Fatal("Count of elements should be equal to 1")
	}

	c.Evict(1)
	if c.ElementsCount() != 0 {
		t.Fatal("Count of elements should be equal to 0")
	}

}

func TestLRUIntKey(t *testing.T) {
	c := New(&Option{0, 2})

	val := &cacheItem{[]byte("ab")}

	c.Set(10, val)

	if v, ok := c.Get(10); !ok || v != val {
		t.Fatal("Int key tests failed")
	}

}

func TestLRUFull(t *testing.T) {
	c := New(&Option{0, 2})

	v := &cacheItem{[]byte("a")}

	c.Set(0, v)
	c.Set(1, v)

	if !c.Full() {
		t.Fatal("The cache should be full now.")
	}
}

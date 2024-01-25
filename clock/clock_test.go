package clock

import (
	"testing"
)

func TestBasic(t *testing.T) {
	keys := []int{1, 2, 3, 4, 5}
	cache := New[int, int](5)
	for _, k := range keys {
		cache.Set(k, k*10)
	}

	for _, k := range keys {
		val, ok := cache.Get(k)
		if !ok {
			t.Fatal(k, "not in cache")
		}
		if val != k*10 {
			t.Fatal(val, "wrong value in cache")
		}
	}
	cache.Get(1)
	if cache.list[0].usage != true {
		t.Fatal("hit not recorded")
	}
	cache.Set(13, 130)
	pos, ok := cache.index[13]
	if cache.list[pos].key == 130 || !ok {
		t.Fatal("fail evict-insert")
	}
}

func TestTriggerParallelHits(t *testing.T) {
	sv := New[int, int](2)
 	iters := 100
	for i := 0; i < iters; i++ {
		go sv.Set(i, i)
		go sv.Get(i)
		go sv.Set(i, i+1)
		go sv.Get(i)
	}
}

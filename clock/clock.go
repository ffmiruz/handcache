package clock

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	mu    sync.Mutex
	list  []object[K, V]
	index map[K]int
	size  int
	hand  int
}

type object[K comparable, V any] struct {
	key   K
	value V
	usage bool
}

func New[K comparable, V any](size int) *Cache[K, V] {
	return &Cache[K, V]{
		size:  size,
		list:  make([]object[K, V], 0, size),
		index: make(map[K]int, size),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	pos, ok := c.index[key]
	if ok {
		c.list[pos].value = value
		c.list[pos].usage = true
		return
	}
	// Eviction path.
	if len(c.list) == c.size {
		for c.hand < len(c.list) {
			// Target identified. Evict then insert in place.
			if c.list[c.hand].usage == false {
				delete(c.index, c.list[c.hand].key)
				c.list[c.hand] = object[K, V]{key, value, false}
				c.index[key] = c.hand
				c.hand++
				c.hand = c.hand % c.size
				return
			}
			c.list[c.hand].usage = false
			c.hand++
			c.hand = c.hand % c.size
		}
	}
	c.index[key] = len(c.list)
	c.list = append(c.list, object[K, V]{key, value, false})
	return
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	pos, ok := c.index[key]
	if ok {
		c.list[pos].usage = true
		return c.list[pos].value, true
	}
	var v V
	return v, false
}
